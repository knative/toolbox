package artifact

import (
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/artifact/platform"
	"knative.dev/toolbox/magetasks/pkg/output"
	"knative.dev/toolbox/magetasks/pkg/output/color"
)

const imageReferenceKey = "oci.image.reference"

// Image is an OCI image that will be built from a binary.
type Image struct {
	config.Metadata
	Labels        map[string]config.Resolver
	Architectures []platform.Architecture
}

func (i Image) GetType() string {
	return "💿"
}

// ImageReferenceOf will try to fetch an image reference from image build result.
func ImageReferenceOf(img Image) config.Resolver {
	return func() string {
		result, ok := config.Actual().Context.Value(BuildKey(img)).(config.Result)
		if !ok || result.Failed() {
			return noImageReference(img)
		}
		ref, ok := result.Info[imageReferenceKey]
		if !ok {
			return noImageReference(img)
		}
		str, ok := ref.(string)
		if !ok {
			return noImageReference(img)
		}
		return str
	}
}

func noImageReference(artifact config.Artifact) string {
	output.Println(color.Yellow("WARNING"),
		" can't resolve image reference for: ", artifact.GetName())
	return ""
}
