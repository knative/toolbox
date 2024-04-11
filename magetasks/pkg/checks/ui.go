package checks

import (
	"fmt"

	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/output/color"
)

func skipBecauseOfMissingConfig(notifier config.Notifier, configFiles interface{}) {
	skipBecauseOf(notifier,
		fmt.Sprintf("config file(s) %v don't exists. Skipping.", configFiles))
}

func skipBecauseOf(notifier config.Notifier, reason string) {
	notifier.Notify(color.Yellow(reason))
}
