// Other command
package main

import (
	"log"

	"knative.dev/toolbox/magetasks/tests/example/pkg/metadata"
)

func main() {
	log.Printf("Version: %s\n", metadata.Version)
	log.Println("Images:")
	images := []metadata.Image{
		metadata.DummyImage,
		metadata.SampleImage,
	}
	for _, image := range images {
		log.Printf(" %s: %s\n", image, metadata.ResolveImage(image))
	}
}
