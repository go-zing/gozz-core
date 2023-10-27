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
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	testParseData = `package x

// +zz:test
// comment
type T struct{}

// +zz:test
/*
lines comment
*/
type ( 
	T2 interface{
		// field
		// +zz:test
		Foo()
	}
)

// +zz:test
// comment
var (
	V0 = 0
	// +zz:test
	V1 = 1
)

// +zz:test
var V2 = 2

// +zz:test
func F0(){}
`
)

func TestParse(t *testing.T) {
	if err := os.WriteFile("test.go", []byte(testParseData), 0o644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.go")
	decls, err := ParseFileOrDirectory(".", AnnotationPrefix)
	if err != nil {
		t.Fatal(err)
	}

	for _, decl := range decls {
		rel := decl.RelFilename("{{ .Package }}_{{ .Name }}_{{ .Filename }}", "")
		if !strings.HasSuffix(rel, fmt.Sprintf("%s_%s_%s", "x", decl.Name(), "test.go")) {
			t.Fatal(rel)
		}
	}

	entities := decls.Parse(test{}, nil)
	for _, entity := range entities {
		if len(entity.Args) != 0 {
			t.Fatal()
		}
	}
}
