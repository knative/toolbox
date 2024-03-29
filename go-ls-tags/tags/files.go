package tags

/*
Copyright 2022 The Knative Authors

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

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"knative.dev/toolbox/pkg/logging"
)

func (l Lister) files(ctx context.Context) ([]string, error) {
	var files []string
	err := filepath.Walk(l.Directory, func(fullpath string, info os.FileInfo, err error) error {
		if err != nil {
			return errwrap(err)
		}
		if info.IsDir() || l.isExcluded(ctx, fullpath) {
			return nil
		}

		if filepath.Ext(fullpath) == "."+l.Extension {
			files = append(files, fullpath)
		}
		return nil
	})
	if err != nil {
		return nil, errwrap(err)
	}
	return files, nil
}

func (l Lister) isExcluded(ctx context.Context, fullpath string) bool {
	log := logging.FromContext(ctx)
	relative, err := filepath.Rel(l.Directory, fullpath)
	if err != nil {
		log.Fatal(err)
	}
	relative = path.Clean(relative)
	for _, exclude := range l.Exclude {
		parts := strings.Split(relative, "/")
		for _, part := range parts {
			if part == exclude {
				log.Debugf("Excluding file %s", relative)
				return true
			}
		}
	}
	return false
}
