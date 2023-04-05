package toolbox_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
	"knative.dev/toolbox/pkg/gowork"
)

func TestBannedDeps(t *testing.T) {
	ban := []string{
		"k8s.io/test-infra",
		"istio.io/test-infra",
	}
	s := gowork.RealSystem{}
	mod, err := gowork.Current(s, s)
	require.NoError(t, err)
	var bytes []byte
	modPath := fmt.Sprintf("/%s/go.mod", mod.Path)
	bytes, err = os.ReadFile(modPath)
	require.NoError(t, err)
	file, err := modfile.Parse(modPath, bytes, nil)
	require.NoError(t, err)
	for _, r := range file.Require {
		// Look for requirements that have the prefix of domain.
		for _, b := range ban {
			assert.False(t, strings.HasPrefix(r.Mod.Path, b))
		}
	}
}
