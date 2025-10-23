package gonfig

import (
	"gopkg.in/yaml.v3"
)

type YAMLFile struct{}

func NewYAMLFile(options GonfigFileOptions) *YAMLFile {
	return &YAMLFile{}
}

func (yf *YAMLFile) encode(config any) ([]byte, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (yf *YAMLFile) decode(data []byte, config any) error {
	return yaml.Unmarshal(data, config)
}
