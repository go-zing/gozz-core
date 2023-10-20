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
	"strconv"
	"strings"
)

type (
	// DeclEntity represents annotated ast.Decl with parsed args and options
	DeclEntity struct {
		*AnnotatedDecl

		Plugin  string
		Args    []string
		Options Options
	}

	DeclEntities []DeclEntity

	// FieldEntity represents annotated ast.Field with parsed args and options
	FieldEntity struct {
		*AnnotatedField

		Args    []string
		Options Options
	}

	FieldEntities []FieldEntity

	// Options carries parsed key-value options from annotations
	Options map[string]string
)

const (
	AnnotationSeparator       = ":"
	EscapeAnnotationSeparator = `\u003A`
	KeyValueSeparator         = "="
)

// Get option value by key from Options map. if empty return default value from def
func (opt Options) Get(key string, def string) string {
	if v, ok := opt[key]; ok && len(v) > 0 {
		return v
	}
	return def
}

// Exist checks key in Options map. if exist and not empty then return strconv.ParseBool result
func (opt Options) Exist(key string) bool {
	if v, ok := opt[key]; ok {
		ok2, _ := strconv.ParseBool(v)
		return len(v) == 0 || ok2
	}
	return false
}

// parseAnnotation parse annotation string
// annotation strings would split by ":" and check first matches provided name
// if not matched then return ok=false
// on matched rests items would be divided into args and options according to args count
// args is strings slice and options is key-value pairs split by "="
// extOptions would fill options if extOptions key not in parsed options
//
// annotation format  $name:$args1:$args2:...$argsN:$key1=$value1:$key2=$value2:...
//
// for example
//
// params:
// 	 annotation  foo:args1:args2:key1=value1:key2=value2
//   name        foo
//   argsCount   2
//   extOptions  [key3:value3 key4:value4]
//
// returns:
// 	 args        [args1 args2]
//   options     [key1:value1 key2:value2 key3:value3 key4:value4]
//   ok          true
func parseAnnotation(annotation, name string, argsCount int, extOptions map[string]string) (args []string, options map[string]string, ok bool) {
	sp := strings.Split(EscapeAnnotation(annotation), AnnotationSeparator)
	if sp[0] != name || len(sp)-1 < argsCount {
		return
	}
	options = make(map[string]string)
	SplitKVSlice2Map(sp[1+argsCount:], KeyValueSeparator, options)

	for k, v := range options {
		options[k] = UnescapeAnnotation(v)
	}

	for k, v := range extOptions {
		if _, exist := options[k]; exist {
			continue
		}
		options[k] = v
	}
	return sp[1 : 1+argsCount], options, true
}

func EscapeAnnotation(str string) string {
	return strings.Replace(str, `\:`, EscapeAnnotationSeparator, -1)
}

func UnescapeAnnotation(str string) string {
	return strings.Replace(str, EscapeAnnotationSeparator, AnnotationSeparator, -1)
}

// GroupByDir groups entities into string map by declaration file dir
func (entities DeclEntities) GroupByDir() (m map[string]DeclEntities) {
	return entities.GroupBy(func(entity DeclEntity) string {
		return filepath.Dir(entity.File.Path)
	})
}

// GroupBy groups entities into string map by function return a string from entity
func (entities DeclEntities) GroupBy(fn func(entity DeclEntity) string) (m map[string]DeclEntities) {
	m = make(map[string]DeclEntities)
	for _, entity := range entities {
		if key := fn(entity); len(key) > 0 {
			m[key] = append(m[key], entity)
		}
	}
	return
}

// ParseFields parses decl fields annotation and returns FieldEntities
func (entity *DeclEntity) ParseFields(argsCount int, options map[string]string) (fields FieldEntities) {
	for _, field := range entity.Fields {
		fields = append(fields, field.Parse(entity.Plugin, argsCount, options)...)
	}
	return
}
