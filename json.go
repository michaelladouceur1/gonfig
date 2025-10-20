package gonfig

import (
	"encoding/json"
	"os"
)

func NewJSONFile(options GonfigFileOptions) (*JSONFile, error) {
	jf := &JSONFile{
		File: File{
			rootDir: options.RootDir,
			name:    options.Name,
			path:    options.RootDir + "/" + options.Name + ".json",
		},
	}
	return jf, nil
}

func (jf *JSONFile) save(config any) error {
	jf.mu.Lock()
	defer jf.mu.Unlock()

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(jf.path, data, 0644)
}

func (jf *JSONFile) load(config any) error {
	jf.mu.Lock()
	defer jf.mu.Unlock()

	data, err := os.ReadFile(jf.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &config)
}
