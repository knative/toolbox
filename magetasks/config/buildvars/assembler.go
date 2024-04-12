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

package buildvars

import "knative.dev/toolbox/magetasks/config"

// Operation performs some operation on Builder and return modified one.
type Operation func(Builder) Builder

type Operator interface {
	Operation() Operation
}

// Assemble will assemble a set of operations into the build variables.
func Assemble(operators []Operator) config.BuildVariables {
	b := Builder{}
	for _, operator := range operators {
		b = operator.Operation()(b)
	}
	return b.Build()
}
