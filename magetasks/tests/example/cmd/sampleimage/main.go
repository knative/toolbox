// Sampleimage command
package main

import (
	"log"

	"knative.dev/toolbox/magetasks/tests/example/pkg/metadata"
)

func main() {
	log.Printf("Sample Image! Version: %s\n", metadata.Version)
}
