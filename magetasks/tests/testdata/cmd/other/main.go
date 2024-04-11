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

// Other command
package main

import (
	"log"

	"knative.dev/toolbox/magetasks/tests/testdata/pkg/metadata"
)

func main() {
	log.Printf("Version: %s\n", metadata.Version)
	log.Println("Images:")
	images := []metadata.Image{
		metadata.PlaceholderImage,
		metadata.SampleImage,
	}
	for _, image := range images {
		log.Printf(" %s: %s\n", image, metadata.ResolveImage(image))
	}
}
