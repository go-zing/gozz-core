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
	"path/filepath"
	"reflect"
	"testing"
)

var pkg = reflect.TypeOf(File{}).PkgPath()

func TestGetImportPath(t *testing.T) {
	rel := "./test/xxx"
	ret := GetImportPath(rel)
	if ret != filepath.Join(pkg, rel) {
		t.Fatal(ret)
	}
	if GetImportPath(".") != pkg {
		t.Fatal(ret)
	}
}

func TestGetImportName(t *testing.T) {
	rel := "./test/xxx"
	ret := GetImportName(rel)
	if ret != "xxx" {
		t.Fatal(ret)
	}
}
