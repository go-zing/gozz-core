/*
 * Copyright (c) 2023 Maple Wu <justmaplewu@gmail.com>
 *   National Electronics and Computer Technology Center, Thailand
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zcore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const cacheFileName = ".gozzcache"

var (
	importNameCache        = new(sync.Map)
	importPathCache        = new(sync.Map)
	importPackageNameCache = new(sync.Map)
	importPackageDirCache  = new(sync.Map)
	modFileCache           = new(sync.Map)

	importPathMutex sync.Mutex

	cacheKey = map[*sync.Map]string{
		importNameCache:        "importName",
		importPathCache:        "importPath",
		importPackageNameCache: "importPackageName",
		importPackageDirCache:  "importPackageDir",
		modFileCache:           "modFile",
	}
)

func InitCacheStore() {
	cache, _ := ioutil.ReadFile(cacheFileName)
	mm := make(map[string]map[string]string)
	if json.Unmarshal(cache, &mm) == nil {
		for store, key := range cacheKey {
			if v, ok := mm[key]; ok {
				for k, vv := range v {
					store.Store(k, vv)
				}
			}
		}
	}
}

func FlushCacheStore() {
	mm := make(map[string]map[string]string)
	for store, key := range cacheKey {
		mm[key] = make(map[string]string)
		store.Range(func(k, v interface{}) bool {
			kk, _ := k.(string)
			vk, _ := v.(string)
			mm[key][kk] = vk
			return true
		})
	}
	cache, _ := json.Marshal(mm)
	_ = ioutil.WriteFile(cacheFileName, cache, 0644)
}

// loadWithStore try loads key from sync.Map or execute provided fn to store valid results
func loadWithStore(key string, m *sync.Map, fn func() string) (r string) {
	if v, ok := m.Load(key); ok {
		return v.(string)
	} else if r = fn(); len(r) > 0 {
		m.Store(key, r)
	}
	return
}

func GetPackageImportName(pkg, dir string) (output string) {
	return loadWithStore(fmt.Sprintf("%s#%s", pkg, dir), importPackageNameCache, func() string {
		ret, _ := ExecCommand(`go list -f "{{ .Name }}" `+strconv.Quote(pkg), dir)
		return ret
	})
}

func GetPackageImportDir(pkg, dir string) (output string) {
	return loadWithStore(fmt.Sprintf("%s#%s", pkg, dir), importPackageDirCache, func() string {
		ret, _ := ExecCommand(`go list -f "{{ .Dir }}" `+strconv.Quote(pkg), dir)
		return ret
	})
}

// ExecCommand execute command in provide directory and get stdout,stderr as string,error
func ExecCommand(command, dir string) (output string, err error) {
	stderr := &bytes.Buffer{}
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	cmd.Stderr = stderr
	r, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s:\n%s", err.Error(), stderr.String())
	}
	return UnsafeBytes2String(bytes.TrimSpace(r)), nil
}

// GetModFile get directory direct mod file by execute "go env GOMOD"
func GetModFile(dir string) string {
	return loadWithStore(dir, modFileCache, func() string {
		modFile, _ := ExecCommand("go env GOMOD", dir)
		return modFile
	})
}

// IsStandardImportPath check import path is whether golang standard library
func IsStandardImportPath(path string) bool {
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	elem := path[:i]
	return !strings.Contains(elem, ".")
}

// GetImportName get filename or directory module import name
// if file is not exist then return a relative calculated result from module environments
func GetImportName(filename string) string {
	return loadWithStore(filename, importNameCache, func() (name string) {
		name, dir := executeWithDir(filename, `go list -f "{{.Name}}"`)
		if len(dir) == 0 || len(name) > 0 {
			return
		}
		// use import path base
		if p := GetImportPath(dir); len(p) > 0 {
			return importNameReplacer.Replace(path.Base(p))
		}
		// use directory base
		return importNameReplacer.Replace(filepath.Base(dir))
	})
}

// GetImportName get filename or directory module import path
// if file is not exist then return a relative calculated result from module environments
func GetImportPath(filename string) string {
	importPathMutex.Lock()
	defer importPathMutex.Unlock()
	return loadWithStore(filename, importPathCache, func() (p string) {
		p, dir := executeWithDir(filename, `go list -f "{{.ImportPath}}"`)
		if len(dir) == 0 || len(p) > 0 {
			return
		}

		// get exist directory
		tmp := dir
		for {
			if _, e := os.Stat(tmp); e == nil {
				break
			}
			tmp = filepath.Dir(tmp)
		}

		// get nearest module path
		modDir := filepath.Dir(GetModFile(tmp))
		modName, err := ExecCommand("go list -m", modDir)
		if err != nil {
			return
		}

		// computed module package import path
		rel, err := filepath.Rel(modDir, dir)
		if err != nil {
			return
		}
		return path.Join(modName, strings.Replace(rel, string(filepath.Separator), "/", -1))
	})
}

// executeInDir try executes command in provided directory or parent if filename is not directory
// return execute output and directory
func executeWithDir(filename string, command string) (ret, dir string) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return
	}

	// check file exist and is directory
	if st, e := os.Stat(filename); (e == nil && st.IsDir()) || !strings.HasSuffix(filename, ".go") {
		dir = filename
	} else {
		dir = filepath.Dir(filename)
	}

	ret, _ = ExecCommand(command+" "+strconv.Quote(dir), dir)
	return
}

// FixPackage modify or add selector package to provide name according to src and dst import module info
func FixPackage(name, srcImportPath, dstImportPath string, srcImports, dstImports Imports) string {
	name, ok := TrimPrefix(name, "*")
	ptr := ""
	if ok {
		ptr = "*"
	}

	sp := strings.Split(name, ".")
	if len(sp) == 1 {
		if token.IsExported(name) && srcImportPath != dstImportPath {
			return ptr + dstImports.Add(srcImportPath) + "." + name
		}
		return ptr + name
	}

	if pkgImportPath := srcImports.Which(sp[0]); pkgImportPath == dstImportPath {
		return ptr + sp[1]
	} else if len(pkgImportPath) == 0 {
		return ptr + name
	} else {
		return ptr + dstImports.Add(pkgImportPath) + "." + sp[1]
	}
}
