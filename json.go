package gonfig

import (
	"encoding/json"
)

type JSONFile struct{}

func NewJSONFile(options GonfigFileOptions) *JSONFile {
	return &JSONFile{}
}

func (jf *JSONFile) encode(config any) ([]byte, error) {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (jf *JSONFile) decode(data []byte, config any) error {
	return json.Unmarshal(data, config)
}
