/*
Copyright 2019 The Knative Authors

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

// Package metautil supports various needs for running tests
package metautil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"knative.dev/toolbox/pkg/prow"
)

const (
	filename = "metadata.json"
)

// Client holds metadata as a string:string map, as well as path for storing
// metadata
type Client struct {
	metadata map[string]string
	Path     string
}

// NewClient creates a client, takes custom directory for storing `metadata.json`.
// It reads existing `metadata.json` file if it exists, otherwise creates it.
// Errors out if there is any file i/o problem other than file not exist error.
func NewClient(dir string) (*Client, error) {
	c := &Client{
		metadata: make(map[string]string),
	}
	if dir == "" {
		log.Println("Getting artifacts dir from prow")
		dir = prow.GetLocalArtifactsDir()
	}
	c.Path = path.Join(dir, filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			return nil, fmt.Errorf("failed creating directory: %w", err)
		}
	}
	return c, nil
}

// sync is shared by Get and Set, invoked at the very beginning of each, makes
// sure the file exists, and loads the content of file into c.metadata
func (c *Client) sync() error {
	_, err := os.Stat(c.Path)
	if os.IsNotExist(err) {
		body, _ := json.Marshal(&c.metadata)
		err = ioutil.WriteFile(c.Path, body, 0777)
	} else {
		var body []byte
		body, err = ioutil.ReadFile(c.Path)
		if err == nil {
			err = json.Unmarshal(body, &c.metadata)
		}
	}

	return err
}

// Set sets key:val pair, and overrides if it exists
func (c *Client) Set(key, val string) error {
	err := c.sync()
	if err != nil {
		return err
	}
	if oldVal, ok := c.metadata[key]; ok {
		log.Printf("Overriding meta %q:%q with new value %q", key, oldVal, val)
	}
	c.metadata[key] = val
	body, _ := json.Marshal(c.metadata)
	return ioutil.WriteFile(c.Path, body, 0777)
}

// Get gets val for key
func (c *Client) Get(key string) (string, error) {
	if _, err := os.Stat(c.Path); err != nil && os.IsNotExist(err) {
		return "", fmt.Errorf("file %q doesn't exist", c.Path)
	}
	var res string
	err := c.sync()
	if err == nil {
		if val, ok := c.metadata[key]; ok {
			res = val
		} else {
			err = fmt.Errorf("key %q doesn't exist", key)
		}
	}
	return res, err
}
