package factory

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ReadConfig(f string) (*Config, error) {
	if content, err := ioutil.ReadFile(f); err != nil {
		return nil, err
	} else {
		config := &Config{}

		if yamlErr := yaml.Unmarshal(content, config); yamlErr != nil {
			return nil, yamlErr
		}
		return config, nil
	}
}
