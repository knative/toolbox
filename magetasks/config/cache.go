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

	"knative.dev/toolbox/magetasks/pkg/cache"
)

// Cache return a cache.Cache implementation based on context.Context.
func Cache() cache.Cache {
	return contextCache{}
}

type contextCache struct{}

func (c contextCache) Compute(
	key interface{},
	provider func() (interface{}, error),
) (interface{}, error) {
	value := fromContext(key)
	if value != nil {
		return value, nil
	}
	value, err := provider()
	if err != nil {
		return nil, err
	}
	saveInContext(key, value)
	return value, nil
}

func (c contextCache) Drop(key interface{}) interface{} {
	value := fromContext(key)
	saveInContext(key, nil)
	return value
}

func saveInContext(cacheKey interface{}, value interface{}) {
	WithContext(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, cacheKey, value)
	})
}

func fromContext(cacheKey interface{}) interface{} {
	return Actual().Context.Value(cacheKey)
}
