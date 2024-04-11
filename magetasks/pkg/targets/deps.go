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

package targets

import (
	"context"
	"os"

	"github.com/magefile/mage/mg"
)

// Deps is a wrapper for mg.CtxDeps that can be disabled by setting the
// SKIP_DEPS environment variable to any value.
func Deps(ctx context.Context, fns ...any) {
	if useDeps() {
		mg.CtxDeps(ctx, fns...)
	}
}

func useDeps() bool {
	skipDeps := os.Getenv("SKIP_DEPS") != ""
	return !skipDeps
}
