package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/toolbox/modscope/cli"
)

func TestCurrentCmd(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	app := cli.App{}
	c := app.Command()
	c.SetOut(&buf)
	c.SetArgs([]string{"current"})
	err := c.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "knative.dev/toolbox\n", buf.String())
}
