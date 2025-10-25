package gonfig

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Gonfig[T any] struct {
	Config     *T
	file       *File
	validators []func(T) error
}

func NewGonfig[T any](config *T, fileOptions GonfigFileOptions) (*Gonfig[T], error) {
	gonfig := &Gonfig[T]{
		Config: config,
	}

	gonfig.file = NewFile(fileOptions)

	if err := gonfig.initialize(); err != nil {
		return nil, err
	}

	if fileOptions.Watch {
		go func() {
			if err := gonfig.watchFile(fileOptions.ValidationMode); err != nil {
				log.Println("Error watching config file:", err)
			}
		}()
	}

	return gonfig, nil
}

func (g *Gonfig[T]) AddValidator(validator func(T) error) {
	g.validators = append(g.validators, validator)
}

func (g *Gonfig[T]) Validate() error {
	return g.validate(*g.Config)
}

func (g *Gonfig[T]) Update(data T) error {
	err := g.validate(data)
	if err != nil {
		return err
	}
	g.Config = &data
	return nil
}

func (g *Gonfig[T]) Save() error {
	return g.file.save(g.Config)
}

func (g *Gonfig[T]) Load() error {
	return g.file.load(g.Config)
}

func (g *Gonfig[T]) PrintConfig() error {
	data, err := g.file.toString(g.Config)
	if err != nil {
		return err
	}
	println(data)
	return nil
}

func (g *Gonfig[T]) initialize() error {
	if !g.file.fileExists() {
		return g.file.save(g.Config)
	}
	return g.file.load(g.Config)
}

func (g *Gonfig[T]) validate(config T) error {
	for _, validator := range g.validators {
		if err := validator(config); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gonfig[T]) watchFile(mode ValidationMode) error {
	callbackChan := make(chan fsnotify.Event)
	if err := g.file.watchFile(callbackChan); err != nil {
		return err
	}

	for range callbackChan {
		if mode == VMRevert {
			copyConfig := *g.Config
			if err := g.file.load(&copyConfig); err != nil {
				log.Println("Error loading config on config update:", err)
				continue
			}
			err := g.validate(copyConfig)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				log.Println("Validation error on config update:", err)
				if saveErr := g.file.saveSilent(g.Config); saveErr != nil {
					log.Println("Error reverting to last known good config:", saveErr)
				}
				continue
			}
		} else if mode == VMWarn {
			if err := g.file.load(g.Config); err != nil {
				log.Println("Error loading config on config update:", err)
				continue
			}
			err := g.validate(*g.Config)
			if err != nil {
				log.Println("Validation warning on config update:", err)
			}
		} else {
			log.Println("Unknown validation mode:", mode)
		}
	}

	return nil
}
