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

package image

import "os"

// BaseName returning a basename of the images.
func BaseName() string {
	return env("localhost", "KO_DOCKER_REPO", "IMAGE_BASENAME")
}

// BaseNameSeparator returns a separator between name and basename of OCI image.
func BaseNameSeparator() string {
	return env("/", "IMAGE_BASENAME_SEPARATOR")
}

func env(defaultVal string, keys ...string) string {
	for _, key := range keys {
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
	}
	return defaultVal
}
