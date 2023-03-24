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

// config.go contains configurations for flaky tests reporting

package config

import (
	"embed"
	"log"
	"os"
	"path"

	"sigs.k8s.io/yaml"
)

//go:embed config.yaml
var configFs embed.FS

// configFile saves all information we need, this path is caller based
const configFile = "config/config.yaml"

var JobConfigs []JobConfig

// Config contains all job configs for flaky tests reporting
type Config struct {
	JobConfigs []JobConfig `yaml:"jobConfigs"`
}

// JobConfig is initial configuration for a given repo, defines which job to scan
type JobConfig struct {
	Name          string         `yaml:"name"` // name of job to analyze
	Org           string         `yaml:"org"`  // org to test job on
	Repo          string         `yaml:"repo"` // repository to test job on
	Type          string         `yaml:"type"`
	IssueRepo     string         `yaml:"issueRepo,omitempty"`
	SlackChannels []SlackChannel `yaml:"slackChannels,omitempty"`
}

// SlackChannel contains Slack channels info
type SlackChannel struct {
	Name     string `yaml:"name"`
	Identity string `yaml:"identity"`
}

func init() {
	contents, err := os.ReadFile(configFile)
	if err != nil {
		// If there's no config file, we will take contents of embedded config
		contents, err = configFs.ReadFile(path.Base(configFile))
	}
	if err != nil {
		log.Printf("Failed to load the config file: %v", err)
		return
	}

	config := &Config{}
	if err = yaml.Unmarshal(contents, config); err != nil {
		log.Printf("Failed to unmarshal %v", contents)
	} else {
		JobConfigs = config.JobConfigs
	}
}
