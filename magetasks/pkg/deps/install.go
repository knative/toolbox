package deps

import (
	"context"
	"fmt"

	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/dotenv"
	"knative.dev/toolbox/magetasks/pkg/files"
	"knative.dev/toolbox/magetasks/pkg/output"
	"knative.dev/toolbox/magetasks/pkg/targets"
	"knative.dev/toolbox/magetasks/pkg/tasks"
)

// Install install build dependencies.
func Install(ctx context.Context) error {
	targets.Deps(ctx, dotenv.Load, output.Setup, files.EnsureBuildDir)
	deps := config.Actual().Dependencies
	dest := fmt.Sprintf("%s/tools", files.BuildDir())
	t := tasks.Start("🔧", "Installing tools",
		deps.Golang.Count()+deps.Binaries.Count() > 0)
	if err := deps.Golang.Install(ctx, t, dest); err != nil {
		return err
	}
	if err := deps.Binaries.Install(ctx, t, dest); err != nil {
		return err
	}
	t.End()
	return nil
}
