// Dummy command
package main

import (
	"log"

	"knative.dev/toolbox/magetasks/tests/example/pkg/metadata"
)

func main() {
	log.Printf("Dummy code! Version: %s\n", metadata.Version)
}
