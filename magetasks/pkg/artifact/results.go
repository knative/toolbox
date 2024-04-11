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
