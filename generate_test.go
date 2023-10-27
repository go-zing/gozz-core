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
	"testing"
)

type test struct {
	Value     string
	MultiLine string
}

func (t test) Name() string                                     { return "test" }
func (t test) Args() (args []string, options map[string]string) { return nil, nil }
func (t test) Description() string                              { return "" }
func (t test) Run(entities DeclEntities) (err error)            { return nil }

const testTemplate = `
var {{ lower .Value | title }} = {{ .Value | camel }}
var {{ .Value | upper }} = {{ quote .Value | kebab }} {{ comment .Value }}
{{ comment .MultiLine }}
`

const testRetTemplate = `// Code generated by gozz:test github.com/go-zing/gozz.

// test

package x

var Testid = testId
var TESTID = "test-id" // TestID
// line1
// line2
`

func TestRenderTemplate(t *testing.T) {
	b, err := RenderTemplate(test{
		Value:     "TestID",
		MultiLine: "line1\nline2",
	}, testTemplate, "x", true, "// test")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte(testRetTemplate)) {
		t.Fatalf("%s", b)
	}
}