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
	"fmt"
	"go/format"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
)

var (
	TemplateFuncs = map[string]interface{}{
		"quote":   strconv.Quote,
		"title":   strings.Title,
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"snake":   strcase.SnakeCase,
		"camel":   strcase.LowerCamelCase,
		"kebab":   strcase.KebabCase,
		"comment": CommentLines,
	}

	templateStore = new(VersionStore)
)

const generateFormat = "// Code generated by %s:%s%s.\n\n"

func CommentLines(comment string) string {
	return "// " + strings.Replace(comment, "\n", "\n// ", -1)
}

// RenderTemplate render golang file template and generate headers
func RenderTemplate(plugin Plugin, templateText string, pkg string, editable bool, ext ...string) (data []byte, err error) {
	bf := BuffPool.Get().(*bytes.Buffer)
	bf.Reset()

	defer BuffPool.Put(bf)

	tips := ". DO NOT EDIT"
	if editable {
		tips = ""
	}

	// code generate comment
	_, _ = fmt.Fprintf(bf, generateFormat, ExecName, plugin.Name(), tips)

	// extra comments before package
	for i, str := range ext {
		bf.WriteString(str)
		bf.WriteRune('\n')
		if len(ext)-1 == i {
			bf.WriteRune('\n')
		}
	}

	// package
	_, _ = fmt.Fprintf(bf, "package %s\n\n", pkg)

	// execute template
	if err = ExecuteTemplate(plugin, templateText, bf); err != nil {
		return
	}

	if data, err = format.Source(bf.Bytes()); err != nil {
		fmt.Printf("%s\n", bf.Bytes())
		return
	}
	return
}

// getTemplate parse text as *template.Template
// parsed templates would be cached in templateStore with template text as key
func getTemplate(text string) (tmpl *template.Template, err error) {
	v, err := templateStore.Load(text, "newest", func() (interface{}, error) {
		return template.New("").Funcs(TemplateFuncs).Parse(text)
	})
	if err != nil {
		return
	}
	tmpl = v.(*template.Template)
	return
}

// ExecuteTemplate parse provide text template and execute template data into writer
func ExecuteTemplate(data interface{}, text string, writer io.Writer) (err error) {
	tmpl, err := getTemplate(text)
	if err != nil {
		return
	}
	return tmpl.Execute(writer, data)
}

// RenderWrite render golang file template and write into filename
func RenderWrite(plugin Plugin, templateText, filename, pkg string, editable bool, ext ...string) (err error) {
	data, err := RenderTemplate(plugin, templateText, pkg, editable, ext...)
	if err != nil {
		return
	}
	_, err = WriteFile(filename, data, 0o664)
	return
}

func RenderWithDefaultTemplate(plugin Plugin, templateText, filename, pkg string, editable bool, ext ...string) (err error) {
	tmpl, err := GetOrWriteDefault(filename+".tmpl", UnsafeString2Bytes(templateText))
	if err != nil {
		return
	}
	return RenderWrite(plugin, string(tmpl), filename, pkg, editable, ext...)
}

// GetOrWriteDefault try read filename or write default data
func GetOrWriteDefault(filename string, defaultData []byte) ([]byte, error) {
	if data, _, err := ReadFile(filename); err == nil {
		return data, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if _, err := WriteFile(filename, defaultData, 0o664); err != nil {
		return nil, err
	}
	return defaultData, nil
}
