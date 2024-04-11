package magetasks

import (
	"os"

	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/files"
	"knative.dev/toolbox/magetasks/pkg/tasks"
)

// Clean will clean project files.
func Clean() {
	t := tasks.Start("🚿", "Cleaning", len(config.Actual().Cleaning) > 0)
	err := os.RemoveAll(files.BuildDir())
	errs := make([]error, 0, 1)
	errs = append(errs, err)
	for _, task := range config.Actual().Cleaning {
		p := t.Part(task.Name)
		pp := p.Starting()
		pp.Done(task.Operation(pp))
	}
	t.End(errs...)
}
