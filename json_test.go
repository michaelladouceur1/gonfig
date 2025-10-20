// go
package gonfig

import (
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	Name string
	Age  int
}

func TestNewJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	opts := GonfigFileOptions{
		RootDir: tmpDir,
		Name:    "testconfig",
	}
	jf, err := NewJSONFile(opts)
	if err != nil {
		t.Fatalf("NewJSONFile error: %v", err)
	}
	expectedPath := filepath.Join(tmpDir, "testconfig.json")
	if jf.path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, jf.path)
	}
}

func TestJSONFileSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	opts := GonfigFileOptions{
		RootDir: tmpDir,
		Name:    "testconfig",
	}
	jf, _ := NewJSONFile(opts)

	cfg := testConfig{Name: "Alice", Age: 42}
	if err := jf.save(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(jf.path); err != nil {
		t.Fatalf("file not created: %v", err)
	}

	// Load back
	var loaded testConfig
	if err := jf.load(&loaded); err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Name != "Alice" || loaded.Age != 42 {
		t.Errorf("loaded config mismatch: %+v", loaded)
	}
}

func TestJSONFileLoadFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	opts := GonfigFileOptions{
		RootDir: tmpDir,
		Name:    "missingfile",
	}
	jf, _ := NewJSONFile(opts)
	var cfg testConfig
	err := jf.load(&cfg)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestJSONFileLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	opts := GonfigFileOptions{
		RootDir: tmpDir,
		Name:    "invalidjson",
	}
	jf, _ := NewJSONFile(opts)
	// Write invalid JSON
	os.WriteFile(jf.path, []byte("{invalid json"), 0644)
	var cfg testConfig
	err := jf.load(&cfg)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
