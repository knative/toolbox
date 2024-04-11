package environment_test

import (
	"testing"

	"gotest.tools/v3/assert"
	"knative.dev/toolbox/magetasks/pkg/environment"
)

func TestNew(t *testing.T) {
	e := environment.New("TAG=v4.5.6", "PUSH_RELEASE=1")
	assert.Equal(t, e[environment.Key("TAG")], environment.Value("v4.5.6"))
	assert.Equal(t, e[environment.Key("PUSH_RELEASE")], environment.Value("1"))
}
