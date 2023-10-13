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
	"plugin"
)

const (
	ExecSuffix       = "zz"
	ExecName         = "go" + ExecSuffix
	AnnotationIdent  = "+"
	AnnotationPrefix = AnnotationIdent + ExecSuffix + ":"
)

type (
	// Plugin represents interface to register as plugin and handles entities
	// builtin Plugin would automate registered on process init.
	// also here supports load external plugin extension from ".so" plugin.
	// external plugin should provide symbol named "Z" and implements Plugin interface
	Plugin interface {

		// Name represents plugin's unique name to register and annotations prefix
		Name() string

		// Args represents arguments and options of plugin to uses.
		// plugin should implement Args to provide args and options description.
		// these infos would be showed in command "list"
		//
		// additionally args count effects annotations parsing and options offsets from name prefix
		// use ":" to split args name with description like "arg_name:arg_help"
		//
		// options provides optional control for plugins
		// plugin can return extra key-value options to describe name-help.
		Args() (args []string, options map[string]string)

		// Description represents summary of plugin .
		// these infos would be showed in command "list"
		Description() string

		// Run is entry of plugin make use of parsed entities from command "run".
		// plugin can use entities parsed from provided name prefix to do awesome things
		Run(entities DeclEntities) (err error)
	}

	// PluginEntity represents Plugin instance and extra options from execute command
	PluginEntity struct {
		Plugin
		Options map[string]string
	}

	PluginEntities []PluginEntity
)

// plugin provides simple registry store for all registered plugins with name
var pluginRegistry = map[string]Plugin{}

func PluginRegistry() map[string]Plugin { return pluginRegistry }

func RegisterPlugin(plugin Plugin) {
	pluginRegistry[plugin.Name()] = plugin
}

func (entities PluginEntities) Run(filename string) (err error) {
	for _, entity := range entities {
		if err = entity.run(filename); err != nil {
			return
		}
	}
	return
}

func (entity PluginEntity) run(filename string) (err error) {
	decls, err := ParseFileOrDirectory(filename, AnnotationPrefix)
	if err != nil {
		return
	}
	return entity.Plugin.Run(decls.Parse(entity, entity.Options))
}

// LoadExtension load filename and lookup symbol named "Z"
// symbol object should implement Plugin or OrmSchemaDriver
func LoadExtension(filename string) (err error) {
	p, err := plugin.Open(filename)
	if err != nil {
		return
	}
	// lookup symbol
	symbol, err := p.Lookup("Z")
	if err != nil {
		return
	}
	// register symbol type
	switch v := symbol.(type) {
	case Plugin:
		RegisterPlugin(v)
	case OrmSchemaDriver:
		RegisterOrmSchemaDriver(v)
	}
	return
}
