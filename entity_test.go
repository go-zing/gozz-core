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
	"testing"
)

func TestParseAnnotation(t *testing.T) {
	if args, opt, ok := parseAnnotation(`test:arg0:arg1:arg2:k1=v1:k1=v2:k2=\:v2`,
		"test", 2, map[string]string{"k2": "v4", "k3": "v3"},
	); !ok || len(args) != 2 || len(opt) != 4 || opt["arg2"] != "" || opt["k1"] != "v1,v2" || opt["k2"] != ":v2" || opt["k3"] != "v3" {
		t.Fatal(args, opt, ok)
	}
}
