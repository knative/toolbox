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

// Builder for a build variables.
type Builder struct {
	bv config.BuildVariables
}

// Add a key/value pair.
func (b Builder) Add(key string, resolver config.Resolver) Builder {
	b.ensureBuildVariables()
	b.bv[key] = resolver
	return b
}

// ConditionallyAdd a key/value pair if cnd is true.
func (b Builder) ConditionallyAdd(cnd func() bool, key string, resolver config.Resolver) Builder {
	if cnd() {
		return b.Add(key, resolver)
	}
	return b
}

// Build a build variables instance.
func (b Builder) Build() config.BuildVariables {
	b.ensureBuildVariables()
	return b.bv
}

func (b *Builder) ensureBuildVariables() {
	if b.bv == nil {
		b.bv = make(config.BuildVariables)
	}
}
