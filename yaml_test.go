package gonfig

import (
	"os"
	"path/filepath"
	"testing"
)

type testYAMLConfig struct {
	Name string
	Age  int
}

func TestNewYAMLFile(t *testing.T) {
	tmpDir := t.TempDir()
	opts := GonfigFileOptions{
		RootDir: tmpDir,
		Name:    "testconfig",
	}
	yf, err := NewYAMLFile(opts)
	if err != nil {
		t.Fatalf("NewYAMLFile error: %v", err)
	}
	expectedPath := filepath.Join(tmpDir, "testconfig.yaml")
	if yf.path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, yf.path)
	}
}

func TestYAMLFileSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	opts := GonfigFileOptions{
		RootDir: tmpDir,
		Name:    "testconfig",
	}
	yf, _ := NewYAMLFile(opts)

	cfg := testYAMLConfig{Name: "Bob", Age: 33}
	if err := yf.save(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	if _, err := os.Stat(yf.path); err != nil {
		t.Fatalf("file not created: %v", err)
	}

	var loaded testYAMLConfig
	if err := yf.load(&loaded); err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Name != "Bob" || loaded.Age != 33 {
		t.Errorf("loaded config mismatch: %+v", loaded)
	}
}

func TestYAMLFileLoadFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	opts := GonfigFileOptions{
		RootDir: tmpDir,
		Name:    "missingfile",
	}
	yf, _ := NewYAMLFile(opts)
	var cfg testYAMLConfig
	err := yf.load(&cfg)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestYAMLFileLoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	opts := GonfigFileOptions{
		RootDir: tmpDir,
		Name:    "invalidyaml",
	}
	yf, _ := NewYAMLFile(opts)
	os.WriteFile(yf.path, []byte("invalid: [unclosed"), 0644)
	var cfg testYAMLConfig
	err := yf.load(&cfg)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}
