package metadata_test

import (
	"testing"

	"gotest.tools/v3/assert"
	"knative.dev/toolbox/magetasks/tests/example/pkg/metadata"
)

func TestImagePath(t *testing.T) {
	p := metadata.ImagePath(metadata.SampleImage)
	assert.Check(t, len(p) > 2)
}
