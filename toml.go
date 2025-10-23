package gonfig

import (
	"github.com/pelletier/go-toml/v2"
)

type TOMLFile struct{}

func NewTOMLFile(options GonfigFileOptions) *TOMLFile {
	return &TOMLFile{}
}

func (tf *TOMLFile) encode(config any) ([]byte, error) {
	data, err := toml.Marshal(config)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (tf *TOMLFile) decode(data []byte, config any) error {
	return toml.Unmarshal(data, config)
}
