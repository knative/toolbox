package cli

import (
	"knative.dev/toolbox/pkg/gowork"
)

// OS represents a virtual operating system.
type OS interface {
	gowork.FileSystem
	gowork.Environment
	Abs(filepath string) string
}

var _ OS = gowork.RealSystem{}
