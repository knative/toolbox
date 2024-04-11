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

import (
	"os"
	"strings"

	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/config/buildvars"
	"knative.dev/toolbox/magetasks/pkg/artifact"
)

// InfluenceableReference defines an image reference that can be influenced by
// environment variables. To disable referencing images by they built hash,
// caller should set variable `DONT_REFERENCE_IMAGE_BY_DIGEST=true` - images
// will use version and basename to resolve full OCI image name.
//
// By setting variable defined in EnvVariable field, caller can override the
// full OCI name of the image.
type InfluenceableReference struct {
	Path        string
	EnvVariable string
	artifact.Image
}

func (r InfluenceableReference) Operation() buildvars.Operation {
	return func(builder buildvars.Builder) buildvars.Builder {
		return builder.ConditionallyAdd(
			referenceImageByDigest,
			r.Path,
			environmentOverrideImageReference(r.EnvVariable,
				artifact.ImageReferenceOf(r.Image)),
		)
	}
}

func environmentOverrideImageReference(
	envVariable string, fallbackResolver config.Resolver,
) config.Resolver {
	return func() string {
		if val, ok := os.LookupEnv(envVariable); ok {
			return val
		}
		return fallbackResolver()
	}
}

func dontReferenceImageByDigest() bool {
	if val, ok := os.LookupEnv("DONT_REFERENCE_IMAGE_BY_DIGEST"); ok {
		return strings.ToLower(val) == "true"
	}
	return false
}

func referenceImageByDigest() bool {
	return !dontReferenceImageByDigest()
}
