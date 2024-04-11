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

package errors

import (
	"errors"
	"testing"
)

// Check if expected and actual error values are correct.
func Check(tb testing.TB, got, want error) {
	tb.Helper()
	if want != nil {
		if !errors.Is(got, want) {
			tb.Fatalf("want %v, got %v", want, got)
		}
	} else {
		if got != nil {
			tb.Fatal("got unexpected:", got)
		}
	}
}
