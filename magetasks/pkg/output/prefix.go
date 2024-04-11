package output

import (
	"github.com/fatih/color"
)

func prefix() string {
	return color.New(color.FgCyan).Sprint("[MAGE]")
}
