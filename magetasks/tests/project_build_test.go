//go:build e2e

package tests_test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestProjectBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("short tests only")
	}
	c := mkCmd("./example", "./mage", "clean", "build")
	assertCommandStarted(t, c)
	assertCommandSucceded(t, c)
	c = mkCmd("./example/build/_output/bin",
		fmt.Sprintf("./other-%s-%s", runtime.GOOS, runtime.GOARCH))
	assertCommandStarted(t, c)
	assertCommandSucceded(t, c)
}

func mkCmd(dir, name string, args ...string) *exec.Cmd {
	c := exec.Command(name, args...)
	c.Env = append(
		env(filterOutByName{names: []string{"GOOS", "GOARCH", "GOARM"}}),
		"GOTRACEBACK=all",
	)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func assertCommandStarted(tb testing.TB, c *exec.Cmd) {
	tb.Helper()
	assert.NilError(tb, c.Start())
	tb.Logf("Started `%q` with pid %d",
		append([]string{c.Path}, c.Args...),
		c.Process.Pid)
	tb.Log("Process env:")
	for _, e := range c.Env {
		tb.Logf(" * %s", e)
	}
}

func assertCommandSucceded(tb testing.TB, c *exec.Cmd) {
	tb.Helper()
	assert.NilError(tb, c.Wait())
}

func env(filter envFilter) []string {
	ret := make([]string, 0, len(os.Environ()))
	for _, e := range os.Environ() {
		envPair := strings.SplitN(e, "=", 2)
		key := envPair[0]
		if !filter.include(key) {
			continue
		}

		ret = append(ret, e)
	}
	return ret
}

type filterOutByName struct {
	names []string
}

func (f filterOutByName) include(name string) bool {
	for _, n := range f.names {
		if name == n {
			return false
		}
	}
	return true
}

type envFilter interface {
	include(name string) bool
}
