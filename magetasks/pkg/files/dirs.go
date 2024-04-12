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

package files

import (
	"os"
	"path"
	"strings"
	"sync"

	"github.com/magefile/mage/mg"
	rand "github.com/thanhpk/randstr"
	"github.com/wavesoftware/go-ensure"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/dotenv"
	"knative.dev/toolbox/magetasks/pkg/output"
)

const randomDirLength = 12

var (
	randomDir    = "mage-build-" + rand.String(randomDirLength)
	buildDirOnce sync.Once
)

// EnsureBuildDir creates a build directory.
func EnsureBuildDir() {
	buildDirOnce.Do(func() {
		mg.Deps(dotenv.Load, output.Setup)
		d := path.Join(BuildDir(), "bin")
		ensure.NoError(os.MkdirAll(d, os.ModePerm))
		ensure.NoError(os.MkdirAll(ReportsDir(), os.ModePerm))
		if strings.Contains(ReportsDir(), randomDir) {
			output.Println("📁 Reports directory: ", ReportsDir())
		}
	})
}

// BuildDir returns project build dir.
func BuildDir() string {
	buildDir := os.Getenv("MAGE_BUILD_DIR")
	if buildDir != "" {
		return buildDir
	}
	buildDir = os.Getenv("BUILD_DIR")
	if buildDir != "" {
		return buildDir
	}
	return relativeTo(ProjectDir(), config.Actual().BuildDirPath...)
}

// ReportsDir returns project reports directory.
func ReportsDir() string {
	artifacts := os.Getenv("ARTIFACTS")
	if artifacts != "" {
		return path.Join(artifacts, randomDir)
	}
	return relativeTo(BuildDir(), "reports")
}

// ProjectDir returns project repo directory.
func ProjectDir() string {
	if config.Actual().ProjectDir != "" {
		return config.Actual().ProjectDir
	}
	repoDir, err := os.Getwd()
	ensure.NoError(err)
	return repoDir
}

func relativeTo(to string, paths ...string) string {
	fullpath := make([]string, len(paths)+1)
	fullpath[0] = to
	for ix, elem := range paths {
		fullpath[ix+1] = elem
	}
	return path.Join(fullpath...)
}
