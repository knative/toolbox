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

package output

import (
	"fmt"
	"log"
	"os"

	"knative.dev/toolbox/magetasks/pkg/output/color"
)

// Setup the output of tasks.
func Setup() {
	color.SetupMode()
}

// PrintPending works similar to log.Print, and displays the prefix.
func PrintPending(args ...interface{}) {
	args = appendPrefix(args...)
	printOrFail(args...)
}

// PrintEnd works similar to log.Print, and end line without adding the prefix.
func PrintEnd(args ...interface{}) {
	args = append(args, "\n")
	printOrFail(args...)
}

// Println works similar to log.Println.
func Println(args ...interface{}) {
	args = appendPrefix(args...)
	args = append(args, "\n")
	printOrFail(args...)
}

// Printlnf works similar to log.Printf.
func Printlnf(format string, args ...interface{}) {
	Println(fmt.Sprintf(format, args...))
}

// printOrFail works similar to log.Print.
func printOrFail(args ...interface{}) {
	_, err := fmt.Fprint(os.Stdout, args...)
	if err != nil {
		log.Fatal(err)
	}
}

func appendPrefix(args ...interface{}) []interface{} {
	return append([]interface{}{prefix(), " "}, args...)
}
