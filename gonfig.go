package gonfig

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Gonfig[T any] struct {
	Config     T
	file       GonfigFile
	validators []func(T) error
}

func NewGonfig[T any](config T, fileOptions GonfigFileOptions) (*Gonfig[T], error) {
	gonfig := &Gonfig[T]{
		Config: config,
	}

	var file GonfigFile
	switch fileOptions.Type {
	case JSON:
		jsonConfig, err := NewJSONFile(fileOptions)
		if err != nil {
			return nil, err
		}
		file = jsonConfig
	case YAML:
		yamlConfig, err := NewYAMLFile(fileOptions)
		if err != nil {
			return nil, err
		}
		file = yamlConfig
	case TOML:
		tomlConfig, err := NewTOMLFile(fileOptions)
		if err != nil {
			return nil, err
		}
		file = tomlConfig
	default:
		return nil, nil
	}

	gonfig.file = file

	if err := gonfig.initialize(); err != nil {
		return nil, err
	}

	if fileOptions.Watch {
		callbackChan := make(chan fsnotify.Event)
		if err := gonfig.file.watchFileChanges(callbackChan); err != nil {
			return nil, err
		}
		go func() {
			for range callbackChan {
				// The config will be reloaded regardless of validation outcome
				if err := gonfig.file.load(&gonfig.Config); err != nil {
					log.Println("Error loading config on config update:", err)
					continue
				}
				err := gonfig.validate(gonfig.Config)
				if err != nil {
					log.Println("Validation error on config update:", err)
					continue
				}
				gonfig.PrintConfig()
			}
		}()
	}

	return gonfig, nil
}

func (g *Gonfig[T]) AddValidator(validator func(T) error) {
	g.validators = append(g.validators, validator)
}

func (g *Gonfig[T]) Validate() error {
	return g.validate(g.Config)
}

func (g *Gonfig[T]) Update(data T) error {
	err := g.validate(data)
	if err != nil {
		return err
	}
	g.Config = data
	return nil
}

func (g *Gonfig[T]) Save() error {
	return g.file.save(&g.Config)
}

func (g *Gonfig[T]) Load() error {
	return g.file.load(&g.Config)
}

func (g *Gonfig[T]) PrintConfig() error {
	data, err := g.file.toString()
	if err != nil {
		return err
	}
	println(string(data))
	return nil
}

func (g *Gonfig[T]) initialize() error {
	if !g.file.fileExists() {
		return g.file.save(&g.Config)
	}
	return g.file.load(&g.Config)
}

func (g *Gonfig[T]) validate(config T) error {
	for _, validator := range g.validators {
		if err := validator(config); err != nil {
			return err
		}
	}
	return nil
}
