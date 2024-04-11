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
	"fmt"

	"knative.dev/toolbox/magetasks/config"
)

// BuildKey returns the config.ResultKey for a build command.
func BuildKey(artifact config.Artifact) config.ResultKey {
	return config.ResultKey(fmt.Sprintf("%s-%s-%s",
		artifact.GetType(), artifact.GetName(), ResultBuilt))
}

// PublishKey returns the config.ResultKey for a publish command.
func PublishKey(artifact config.Artifact) config.ResultKey {
	return config.ResultKey(fmt.Sprintf("%s-%s-%s",
		artifact.GetType(), artifact.GetName(), ResultPublication))
}
