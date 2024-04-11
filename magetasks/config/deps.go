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

package config

import (
	"context"
	"fmt"

	"github.com/magefile/mage/sh"
	"knative.dev/toolbox/magetasks/pkg/tasks"
)

type Dependencies interface {
	Configurator
	Install(ctx context.Context, t *tasks.Task, dest string) error
	Count() int

	merge(other dependencies)
	installs() []string
}

func NewDependencies(deps ...string) Dependencies {
	s := make(map[string]bool, len(deps))
	for _, dep := range deps {
		s[dep] = exists
	}
	return dependencies{
		set: s,
	}
}

const exists = true

type dependencies struct {
	set map[string]bool
}

func (d dependencies) Count() int {
	return len(d.set)
}

func (d dependencies) Install(_ context.Context, t *tasks.Task, dest string) error {
	for _, dep := range d.installs() {
		pp := t.Part(fmt.Sprintf("Go install %q", dep)).Starting()
		env := map[string]string{
			"GOBIN": dest,
		}
		if err := sh.RunWith(env, "go", "install", dep); err != nil {
			pp.Done(err)
			return err
		}
		pp.Notify("installed")
	}
	return nil
}

func (d dependencies) installs() []string {
	keys := make([]string, len(d.set))

	i := 0
	for k := range d.set {
		keys[i] = k
		i++
	}

	return keys
}

func (d dependencies) Configure(cfg Configurable) {
	cfg.Config().Dependencies.Golang.merge(d)
}

func (d dependencies) merge(other dependencies) {
	for _, dep := range other.installs() {
		d.set[dep] = exists
	}
}
