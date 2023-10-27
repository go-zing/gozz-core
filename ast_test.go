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
	"go/ast"
	"go/parser"
	"reflect"
	"strconv"
	"testing"
)

func TestExtractStructFieldsNames(t *testing.T) {
	v, err := parser.ParseExpr("struct{F1 string;F2 int;F3 bool;int;pkg.F4;*pkg.F5}")
	if err != nil {
		t.Fatal(err)
	}
	names := ExtractStructFieldsNames(v.(*ast.StructType))
	if len(names) != 5 {
		t.Fatal(names)
	}
	for i, name := range names {
		if name != "F"+strconv.Itoa(i+1) {
			t.Fatal(i, name)
		}
	}
}

func TestLookupTypSpec(t *testing.T) {
	exp, f := LookupTypSpec(reflect.TypeOf(File{}).Name(), "./", pkg)
	if f == nil {
		t.Fatal("not found")
	}
	if _, ok := exp.(*ast.StructType); !ok {
		t.Fatal("not found")
	}
}
