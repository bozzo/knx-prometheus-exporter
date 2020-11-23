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
)

type Configuration struct {
	dptByGroup    map[string]string
	dptNameByType map[string]string
}

type YmlConfiguration struct {
	Version int                           `yaml:"version"`
	Mapping []YmlConfigurationTypeMapping `yaml:"mapping"`
}

type YmlConfigurationTypeMapping struct {
	DptId   string   `yaml:"dptId"`
	Groups  []string `yaml:"groups"`
	DptName string   `yaml:"dptName"`
}

func (config *Configuration) loadConfiguration() error {
	file := os.Getenv("CONFIG_FILE")
	if file == "" {
		file = "config.yml"
	}
	ymlFile, err := os.Open(file)
	if err != nil {
		return err
	}
	var ymlConfig YmlConfiguration

	decoder := yaml.NewDecoder(ymlFile)
	err = decoder.Decode(&ymlConfig)
	if err != nil {
		return err
	}

	if config.dptByGroup == nil {
		config.dptByGroup = map[string]string{}
		config.dptNameByType = map[string]string{}
	}

	for _, obj := range ymlConfig.Mapping {
		for _, group := range obj.Groups {
			config.dptByGroup[group] = obj.DptId
		}
		config.dptNameByType[obj.DptId] = obj.DptName
	}

	return nil
}

func (config *Configuration) getDptId(group string) string {
	return config.dptByGroup[group]
}

func (config *Configuration) getDptName(dpt string) string {
	return config.dptNameByType[dpt]
}

func (config *Configuration) getDptAndName() map[string]string {
	return config.dptNameByType
}

func (config *Configuration) getGroups() []string {
	groups := make([]string, 0, len(config.dptByGroup))
	for k := range config.dptByGroup {
		groups = append(groups, k)
	}
	return groups
}
