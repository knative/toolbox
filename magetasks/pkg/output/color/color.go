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

package color

import (
	"os"

	"github.com/fatih/color"
)

var (
	// Red text color.
	Red = color.New(color.FgHiRed).Add(color.Bold).SprintFunc()
	// Green text color.
	Green = color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	// Yellow text color.
	Yellow = color.New(color.FgYellow).SprintFunc()
	// Blue text color.
	Blue = color.New(color.FgBlue).SprintFunc()
)

// SetupMode will set the output color mode.
func SetupMode() {
	if val, envset := os.LookupEnv("FORCE_COLOR"); envset && val == "true" {
		color.NoColor = false
	}
}
