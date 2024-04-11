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

package environment

import (
	"errors"
	"fmt"
)

// ErrNotSupported when operation is not supported.
var ErrNotSupported = errors.New("not supported")

// NewVersionResolver creates a VersionResolver using options.
func NewVersionResolver(options ...VersionResolverOption) VersionResolver {
	r := VersionResolver{}

	for _, option := range options {
		option(&r)
	}

	return r
}

// WithValuesSupplier allows to set the values supplier.
func WithValuesSupplier(vs ValuesSupplier) VersionResolverOption {
	return func(r *VersionResolver) {
		r.ValuesSupplier = vs
	}
}

// VersionResolverOption is used to customize creation of VersionResolver.
type VersionResolverOption func(*VersionResolver)

// VersionResolver is used to resolve version information solely on
// environment variables.
type VersionResolver struct {
	VersionKey   Key
	IsApplicable []Check
	ValuesSupplier
}

// Check is used to verify environment values.
type Check Pair

func (e VersionResolver) Version() string {
	values := e.environment()
	if !e.isApplicable(values) {
		return ""
	}
	return string(values[e.VersionKey])
}

func (e VersionResolver) IsLatest(_ string) (bool, error) {
	return false, fmt.Errorf(
		"%w: IsLatest(versionRange string) by environment.VersionResolver",
		ErrNotSupported)
}

func (e VersionResolver) environment() Values {
	supplier := Current
	if e.ValuesSupplier != nil {
		supplier = e.ValuesSupplier
	}
	return supplier()
}

func (e VersionResolver) isApplicable(values Values) bool {
	for _, check := range e.IsApplicable {
		if !check.test(values) {
			return false
		}
	}
	return true
}

func (c Check) test(values Values) bool {
	if c.Value == "" {
		_, ok := values[c.Key]
		return ok
	}
	val, ok := values[c.Key]
	return ok && val == c.Value
}
