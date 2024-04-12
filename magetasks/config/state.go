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
	"log"
)

var state Configurable //nolint:gochecknoglobals

// Actual returns actual config.
func Actual() Config {
	verifyState()
	cfg := state.Config()
	return *cfg
}

// WithContext allows to use context.Context and to amend it.
func WithContext(mutator func(ctx context.Context) context.Context) {
	verifyState()
	ctx := mutator(state.Config().Context)
	state.Config().Context = ctx
}

// Configure with provided Configurable.
func Configure(c Configurable) {
	state = c
	for _, override := range collectConfigurators(c) {
		override.Configure(state)
	}
}

func verifyState() {
	if state == nil {
		log.Fatal("Not configured yet, execute config.Configure() first!")
	}
}

func collectConfigurators(c Configurable) []Configurator {
	cnfrs := make([]Configurator, 0, len(c.Config().Overrides))
	cnfrs = append(cnfrs, c.Config().Overrides...)
	for _, task := range c.Config().Cleaning {
		cnfrs = append(cnfrs, task.Overrides...)
	}
	for _, task := range c.Config().Checks {
		cnfrs = append(cnfrs, task.Overrides...)
	}
	return cnfrs
}
