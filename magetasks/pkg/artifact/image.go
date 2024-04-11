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

package artifact

import (
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/artifact/platform"
	"knative.dev/toolbox/magetasks/pkg/output"
	"knative.dev/toolbox/magetasks/pkg/output/color"
)

const imageReferenceKey = "oci.image.reference"

// Image is an OCI image that will be built from a binary.
type Image struct {
	config.Metadata
	Labels        map[string]config.Resolver
	Architectures []platform.Architecture
}

func (i Image) GetType() string {
	return "💿"
}

// ImageReferenceOf will try to fetch an image reference from image build result.
func ImageReferenceOf(img Image) config.Resolver {
	return func() string {
		result, ok := config.Actual().Context.Value(BuildKey(img)).(config.Result)
		if !ok || result.Failed() {
			return noImageReference(img)
		}
		ref, ok := result.Info[imageReferenceKey]
		if !ok {
			return noImageReference(img)
		}
		str, ok := ref.(string)
		if !ok {
			return noImageReference(img)
		}
		return str
	}
}

func noImageReference(artifact config.Artifact) string {
	output.Println(color.Yellow("WARNING"),
		" can't resolve image reference for: ", artifact.GetName())
	return ""
}
