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

package ldflags

import (
	"fmt"
	"strings"

	"knative.dev/toolbox/magetasks/config"
)

// Builder builds the LD flags by adding values resolvers.
type Builder interface {
	// Add a name and a resolver to the builder.
	Add(name string, resolver config.Resolver) Builder
	// Build into slice args for ldflags.
	Build() []string
	// BuildOnto provided args.
	BuildOnto(args []string) []string
}

// NewBuilder creates a new builder.
func NewBuilder() Builder {
	return &defaultBuilder{
		resolvers: make(map[string]config.Resolver),
	}
}

type defaultBuilder struct {
	resolvers map[string]config.Resolver
}

func (d *defaultBuilder) Add(name string, resolver config.Resolver) Builder {
	d.resolvers[name] = resolver
	return d
}

func (d *defaultBuilder) Build() []string {
	collected := make([]string, 0, len(d.resolvers))
	if len(d.resolvers) == 0 {
		return collected
	}
	for name, resolver := range d.resolvers {
		collected = append(collected, fmt.Sprintf("-X %s=%s", name, resolver()))
	}
	return collected
}

func (d *defaultBuilder) BuildOnto(args []string) []string {
	return append(args, "-ldflags", strings.Join(d.Build(), " "))
}
