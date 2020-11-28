/*
 *
 *    Copyright 2020 Boris Barnier <bozzo@users.noreply.github.com>
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type configuration struct {
	dptByGroup    map[string]string
	dptNameByType map[string]string
}

type ymlConfiguration struct {
	Version int                           `yaml:"version"`
	Mapping []ymlConfigurationTypeMapping `yaml:"mapping"`
}

type ymlConfigurationTypeMapping struct {
	DptID   string   `yaml:"dptID"`
	Groups  []string `yaml:"groups"`
	DptName string   `yaml:"dptName"`
}

func (config *configuration) loadConfigurationFromFile() error {
	file := os.Getenv("CONFIG_FILE")
	if file == "" {
		file = "config.yml"
	}
	ymlFile, err := os.Open(filepath.Clean(file))
	if err != nil {
		return err
	}
	var ymlConfig ymlConfiguration

	decoder := yaml.NewDecoder(ymlFile)
	err = decoder.Decode(&ymlConfig)
	if err != nil {
		return err
	}

	config.parseConfig(ymlConfig)

	return nil
}

func (config *configuration) parseConfig(ymlConfig ymlConfiguration) {
	for _, obj := range ymlConfig.Mapping {
		for _, group := range obj.Groups {
			config.dptByGroup[group] = obj.DptID
		}
		config.dptNameByType[obj.DptID] = obj.DptName
	}
}

func (config *configuration) getDptID(group string) string {
	return config.dptByGroup[group]
}

func (config *configuration) getDptName(dpt string) string {
	return config.dptNameByType[dpt]
}

func (config *configuration) getDptAndName() map[string]string {
	return config.dptNameByType
}
