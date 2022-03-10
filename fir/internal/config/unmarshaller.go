package config

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type unmarshaller = func([]byte, *Config) error

var unmarshallers = map[string]unmarshaller{
	".json": jsonUnmarshaller,
	".yaml": yamlUnmarshaller,
	".yml":  yamlUnmarshaller,
}

func jsonUnmarshaller(contents []byte, c *Config) error {
	return json.Unmarshal(contents, &c)
}

func yamlUnmarshaller(contents []byte, c *Config) error {
	return yaml.Unmarshal(contents, &c)
}
