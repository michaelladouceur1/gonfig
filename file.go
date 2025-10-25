package gonfig

import (
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FileType string
type ValidationMode string

const (
	JSON FileType = "json"
	YAML FileType = "yaml"
	TOML FileType = "toml"

	VMRevert ValidationMode = "revert"
	VMWarn   ValidationMode = "warn"
)

type File struct {
	rootDir    string
	name       string
	path       string
	mu         sync.Mutex
	encoder    FileEncoder
	watcher    *fsnotify.Watcher
	pauseWatch bool
}

type FileEncoder interface {
	encode(config any) ([]byte, error)
	decode(data []byte, config any) error
}

type GonfigFileOptions struct {
	Type           FileType
	RootDir        string
	Name           string
	Watch          bool
	ValidationMode ValidationMode
}

func NewFile(options GonfigFileOptions) *File {
	var encoder FileEncoder

	switch options.Type {
	case JSON:
		encoder = &JSONFile{}
	case YAML:
		encoder = &YAMLFile{}
	case TOML:
		encoder = &TOMLFile{}
	default:
		log.Fatalf("Unsupported file type: %s", options.Type)
	}

	return &File{
		rootDir:    options.RootDir,
		name:       options.Name,
		path:       options.RootDir + "/" + options.Name + "." + string(options.Type),
		encoder:    encoder,
		pauseWatch: false,
	}
}

func (f *File) fileExists() bool {
	_, err := os.Stat(f.path)
	return !os.IsNotExist(err)
}

func (f *File) toString(config any) (string, error) {
	data, err := f.encoder.encode(config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *File) load(config any) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	data, err := os.ReadFile(f.path)
	if err != nil {
		return err
	}

	return f.encoder.decode(data, config)
}

func (f *File) save(config any) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	data, err := f.encoder.encode(config)
	if err != nil {
		return err
	}

	return os.WriteFile(f.path, data, 0644)
}

func (f *File) saveSilent(config any) error {
	f.pauseWatch = true
	err := f.save(config)
	f.pauseWatch = false

	return err
}

func (f *File) watchFile(callbackChan chan fsnotify.Event) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	f.watcher = watcher

	err = f.watcher.Add(f.path)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-f.watcher.Events:
				if !ok {
					return
				}
				if f.pauseWatch {
					continue
				}
				callbackChan <- event
			case err, ok := <-f.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("YAMLFile watch error: %v", err)
			}
		}
	}()

	return nil
}
