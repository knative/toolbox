package artifact

import "knative.dev/toolbox/magetasks/pkg/artifact/platform"

// Platform to built binary for.
type Platform struct {
	platform.OS
	platform.Architecture
}
