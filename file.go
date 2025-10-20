package gonfig

import (
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FileType string

const (
	JSON FileType = "json"
	YAML FileType = "yaml"
	// TOMLFileType FileType = "toml"
)

type File struct {
	rootDir string
	name    string
	path    string
	mu      sync.Mutex
	watcher *fsnotify.Watcher
}

type JSONFile struct {
	File
}

type YAMLFile struct {
	File
}

type GonfigFile interface {
	fileExists() bool
	load(config any) error
	save(config any) error
	watchFileChanges(chan fsnotify.Event) error
	toString() (string, error)
}

type GonfigFileOptions struct {
	Type    FileType
	RootDir string
	Name    string
	Watch   bool
}

func (f *File) fileExists() bool {
	_, err := os.Stat(f.path)
	return !os.IsNotExist(err)
}

func (f *File) toString() (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	data, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (f *File) watchFileChanges(callbackChan chan fsnotify.Event) error {
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
