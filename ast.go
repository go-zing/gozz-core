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
	"path/filepath"
)

// AssertFuncType to assert interface fields as function type and try return name
func AssertFuncType(field *ast.Field) (name string, ft *ast.FuncType, ok bool) {
	ft, ok = field.Type.(*ast.FuncType)
	if !ok || len(field.Names) == 0 {
		return
	}
	name = field.Names[0].Name
	return
}

func ExtractAnonymousName(spec ast.Expr) (name *ast.Ident) {
	switch t := spec.(type) {
	case *ast.StarExpr:
		name, _ = t.X.(*ast.Ident)
	case *ast.SelectorExpr:
		name, _ = t.X.(*ast.Ident)
	case *ast.Ident:
		name = t
	}
	return
}

// ExtractStructFieldsNames extracts struct exported fields names
func ExtractStructFieldsNames(typ *ast.StructType) (names []string) {
	if typ.Fields == nil {
		return
	}

	add := func(ident *ast.Ident) {
		if ident != nil && ident.IsExported() {
			names = append(names, ident.Name)
		}
	}

	for _, field := range typ.Fields.List {
		// anonymous field
		if len(field.Names) == 0 {
			add(ExtractAnonymousName(field.Type))
			continue
		}

		// with name
		for _, name := range field.Names {
			add(name)
		}
	}
	return
}

// LookupTypSpec lookup typename in package src path.
func LookupTypSpec(name, dir, pkgPath string) (expr ast.Expr, srcFile *File) {
	if len(pkgPath) == 0 {
		return
	}

	pkgDir, err := ExecCommand(`go list -f "{{ .Dir }} " `+pkgPath, dir)
	if err != nil {
		return
	}

	_, _ = WalkPackage(pkgDir, func(file *File) (err error) {
		object := file.Lookup(name)
		if object == nil || object.Decl == nil {
			return
		}

		if spec, ok := object.Decl.(*ast.TypeSpec); ok {
			switch typ := spec.Type.(type) {
			case *ast.SelectorExpr:
				expr, srcFile = LookupTypSpec(typ.Sel.Name, dir, file.Imports().Which(UnsafeBytes2String(file.Node(typ.X))))
			case *ast.Ident:
				expr, srcFile = LookupTypSpec(typ.Name, dir, pkgPath)
			default:
				expr, srcFile = typ, file
			}
		}

		return filepath.SkipDir
	})
	return
}
