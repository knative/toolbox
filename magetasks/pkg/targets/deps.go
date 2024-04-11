package targets

import (
	"context"
	"os"

	"github.com/magefile/mage/mg"
)

// Deps is a wrapper for mg.CtxDeps that can be disabled by setting the
// SKIP_DEPS environment variable to any value.
func Deps(ctx context.Context, fns ...any) {
	if useDeps() {
		mg.CtxDeps(ctx, fns...)
	}
}

func useDeps() bool {
	skipDeps := os.Getenv("SKIP_DEPS") != ""
	return !skipDeps
}
