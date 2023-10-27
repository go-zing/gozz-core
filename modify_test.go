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
	"io/ioutil"
	"os"
	"testing"
)

var (
	testModifyData = `package s

import (
	"time"
)

var _ = new(time.Time)
`
	testModifyRetData = `package x

import (
	"context"
	time2 "host.com/time"
	"time"
)

var _ = new(time.Time)
`
)

func TestModify(t *testing.T) {
	_ = ioutil.WriteFile("test", []byte(testModifyData), 0o664)
	defer os.Remove("test")
	f, err := ParseFile("test")
	if err != nil {
		t.Fatal(err)
	}
	set := ModifySet{}
	m := set.Add("test")
	m.Imports = f.Imports()
	m.Imports.Add("context")
	m.Imports.Add("host.com/time")
	m.Nodes[f.Ast.Name] = []byte("x")
	if err = m.Apply(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("test")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, []byte(testModifyRetData)) {
		t.Fatal()
	}
}
