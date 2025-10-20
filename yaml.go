package gonfig

import (
	"os"

	"gopkg.in/yaml.v3"
)

func NewYAMLFile(options GonfigFileOptions) (*YAMLFile, error) {
	yf := &YAMLFile{
		File: File{
			rootDir: options.RootDir,
			name:    options.Name,
			path:    options.RootDir + "/" + options.Name + ".yaml",
		},
	}
	return yf, nil
}

func (yf *YAMLFile) save(config any) error {
	yf.mu.Lock()
	defer yf.mu.Unlock()

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(yf.path, data, 0644)
}

func (yf *YAMLFile) load(config any) error {
	yf.mu.Lock()
	defer yf.mu.Unlock()

	data, err := os.ReadFile(yf.path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}
