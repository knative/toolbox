/*
 Copyright 2024 The Knative Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

// Package metadata contains meta information about the project.
package metadata

import (
	pkgimage "knative.dev/toolbox/magetasks/pkg/image"
)

var (
	// DummyImageRef holds information about companion dummy image.
	DummyImageRef = "" //nolint:gochecknoglobals
	// SampleImageRef holds information about companion sample image.
	SampleImageRef = "" //nolint:gochecknoglobals
	// ImageBasename holds a basename of a image, so the development reference
	// could be built from it.
	ImageBasename = "" //nolint:gochecknoglobals
	// ImageBasenameSeparator holds a separator between image basename and name.
	ImageBasenameSeparator = "/" //nolint:gochecknoglobals
)

// Image represents an image.
type Image string

const (
	// DummyImage represents a dummy image.
	DummyImage Image = "dummy"
	// SampleImage represents a sample image.
	SampleImage Image = "sampleimage"
)

// ResolveImage will try to resolve the image reference from set values. If
// DummyImageRef is given it will be used, otherwise the ImageBasename and Version will
// be used.
func ResolveImage(image Image) string {
	ref := func() string {
		switch image {
		case DummyImage:
			return DummyImageRef
		case SampleImage:
			return SampleImageRef
		}
		return ""
	}()
	if ref == "" {
		return pkgimage.FloatToRelease(
			ImageBasename, string(image), ImageBasenameSeparator, Version,
			pkgimage.FloatDirectionDown)
	}
	return ref
}

// ImagePath return a path to the image variable.
func ImagePath(image Image) string {
	variable := func() string {
		switch image {
		case DummyImage:
			return "DummyImageRef"
		case SampleImage:
			return "SampleImageRef"
		}
		panic("unsupported image: " + image)
	}()
	return importPath(variable)
}

// ImageBasenamePath return a path to the image basename variable.
func ImageBasenamePath() string {
	return importPath("ImageBasename")
}

// ImageBasenameSeparatorPath return a path to the image basename separator
// variable.
func ImageBasenameSeparatorPath() string {
	return importPath("ImageBasenameSeparator")
}
