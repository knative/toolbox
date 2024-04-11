package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cardil/ghet/pkg/ghet/download"
	"github.com/cardil/ghet/pkg/ghet/install"
	"knative.dev/toolbox/magetasks/pkg/tasks"
)

type Binaries interface {
	Configurator
	Install(ctx context.Context, t *tasks.Task, destination string) error
	Count() int

	merge(b *binaries)
}

func NewBinaries(bins ...string) Binaries {
	return &binaries{
		bins: bins,
	}
}

type binaries struct {
	bins []string
}

func (b *binaries) Count() int {
	return len(b.bins)
}

func (b *binaries) Install(ctx context.Context, t *tasks.Task, destination string) error {
	for _, binSpec := range b.bins {
		p := t.Part(fmt.Sprintf("Pre-built binary %q", binSpec))
		pp := p.Starting()
		args := download.Args{
			Args:        install.Parse(binSpec),
			Destination: destination,
		}
		bin := args.Asset.FileName.ToString()
		path := fmt.Sprintf("%s/%s", destination, bin)
		if fileExist(path) {
			log.Println("Skipping installation of", binSpec,
				"because it already exists:", path)
			p.Skip("already installed")
			continue
		}
		if err := download.Action(ctx, args); err != nil {
			pp.Done(err)
			return err
		}
		pp.Notify("installed")
	}
	return nil
}

func (b *binaries) Configure(cfg Configurable) {
	cfg.Config().Dependencies.Binaries.merge(b)
}

func (b *binaries) merge(b2 *binaries) {
	b.bins = append(b.bins, b2.bins...)
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
