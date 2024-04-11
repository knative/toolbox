package artifact

import "knative.dev/toolbox/magetasks/config"

// ConfigureDefaults will configure default builders and publishers to be used.
func ConfigureDefaults() {
	config.DefaultBuilders = append(config.DefaultBuilders,
		BinaryBuilder{},
		KoBuilder{},
	)
	config.DefaultPublishers = append(config.DefaultPublishers,
		ListPublisher{},
		KoPublisher{},
	)
}
