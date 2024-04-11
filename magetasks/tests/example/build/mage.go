//go:build ignored

package main

import (
	"os"
	"path"
	"runtime"

	"knative.dev/toolbox/magetasks/entrypoint"
)

func main() {
	bdir := builddir()
	os.Exit(entrypoint.Execute(entrypoint.Context{
		Directories: entrypoint.Directories{
			BuildDir:   bdir,
			ProjectDir: path.Dir(bdir),
			CacheDir:   path.Join(bdir, "_output"),
		},
	}))
}

func builddir() string {
	_, file, _, _ := runtime.Caller(0) //nolint:dogsled
	return path.Dir(file)
}
