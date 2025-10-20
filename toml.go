package gonfig

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func NewTOMLFile(options GonfigFileOptions) (*TOMLFile, error) {
	tf := &TOMLFile{
		File: File{
			rootDir: options.RootDir,
			name:    options.Name,
			path:    options.RootDir + "/" + options.Name + ".toml",
		},
	}
	return tf, nil
}

func (tf *TOMLFile) save(config any) error {
	tf.mu.Lock()
	defer tf.mu.Unlock()

	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(tf.path, data, 0644)
}

func (tf *TOMLFile) load(config any) error {
	tf.mu.Lock()
	defer tf.mu.Unlock()

	data, err := os.ReadFile(tf.path)
	if err != nil {
		return err
	}

	return toml.Unmarshal(data, config)
}
