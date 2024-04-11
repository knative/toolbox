/*
 Copyright 2024 The Knative Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package environment_test

import (
	"testing"

	"gotest.tools/v3/assert"
	"knative.dev/toolbox/magetasks/pkg/environment"
)

func TestNew(t *testing.T) {
	e := environment.New("TAG=v4.5.6", "PUSH_RELEASE=1")
	assert.Equal(t, e[environment.Key("TAG")], environment.Value("v4.5.6"))
	assert.Equal(t, e[environment.Key("PUSH_RELEASE")], environment.Value("1"))
}
