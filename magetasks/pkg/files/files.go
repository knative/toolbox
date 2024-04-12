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
	"os/exec"
)

// DontExists will check if target file dont exist.
func DontExists(file string) bool {
	_, err := os.Stat(file)
	return err != nil && os.IsNotExist(err)
}

// ExecutableAvailable will return true if given executable in available in
// system env.PATH's.
func ExecutableAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
