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

package entrypoint

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/magefile/mage/mage"
	"github.com/wavesoftware/go-retcode"
)

var errUnknownCommand = errors.New("unknown command")

// Directories is a configuration of project directories.
type Directories struct {
	BuildDir   string
	ProjectDir string
	CacheDir   string
}

// InvocationOption is a function that can be used to modify mage.Invocation.
type InvocationOption func(mage.Invocation) mage.Invocation

// Context is a context for the mage execution.
type Context struct {
	Directories
	Options []InvocationOption
}

// Execute is the main entry point for the mage command.
// It will change the current working directory to the build directory given
// as Context BuildDir parameter, if given.
func Execute(ctx Context) int {
	if ctx.BuildDir != "" {
		chdirBuild(ctx.BuildDir)
	}
	return parseAndRun(os.Stdout, os.Stderr, os.Stdin, os.Args[1:], ctx)
}

func parseAndRun(stdout, stderr io.Writer, stdin io.Reader, args []string, ctx Context) int {
	inv, cmd, err := mage.Parse(stderr, stdout, args)
	inv.Stdin = stdin

	if ctx.ProjectDir != "" {
		inv.WorkDir = ctx.ProjectDir
	}
	if ctx.BuildDir != "" {
		inv.Dir = ctx.BuildDir
	}
	if ctx.CacheDir != "" {
		inv.CacheDir = ctx.CacheDir
	}
	for _, opt := range ctx.Options {
		inv = opt(inv)
	}
	return run(inv, cmd, err, stderr)
}

// TODO: remove this function, when https://github.com/magefile/mage/pull/442 lands.
func run(inv mage.Invocation, cmd mage.Command, err error, stderr io.Writer) int {
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	errlog := log.New(inv.Stderr, "", 0)
	if err != nil {
		errlog.Println("Error:", err)
		return retcode.Calc(err)
	}
	inv.Stderr = stderr

	switch cmd {
	case mage.Version, mage.Init, mage.Clean:
		return mage.ParseAndRun(inv.Stdout, inv.Stderr, inv.Stdin, inv.Args)
	case mage.CompileStatic:
		return mage.Invoke(inv)
	case mage.None:
		return mage.Invoke(inv)
	default:
		panic(fmt.Errorf("%w type: %v", errUnknownCommand, cmd))
	}
}

func chdirBuild(builddir string) {
	if err := os.Chdir(builddir); err != nil {
		panic(err)
	}
}
