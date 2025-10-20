// go
package gonfig

import (
	"errors"
	"testing"

	"github.com/fsnotify/fsnotify"
)

// Sample config struct for testing
type TestConfig struct {
	Name string
	Age  int
}

// Mock GonfigFile for testing
type MockFile struct {
	saved   *TestConfig
	loaded  *TestConfig
	exists  bool
	toStr   string
	watchCh chan struct{}
}

func (m *MockFile) save(cfg any) error {
	m.saved = cfg.(*TestConfig)
	return nil
}
func (m *MockFile) load(cfg any) error {
	if m.loaded != nil {
		*cfg.(*TestConfig) = *m.loaded
	}
	return nil
}
func (m *MockFile) toString() (string, error) {
	return m.toStr, nil
}
func (m *MockFile) fileExists() bool {
	return m.exists
}
func (m *MockFile) watchFileChanges(ch chan fsnotify.Event) error {
	m.watchCh = make(chan struct{})
	return nil
}

// Test NewGonfig with mock file
func TestNewGonfig(t *testing.T) {
	cfg := TestConfig{Name: "Alice", Age: 30}
	g := &Gonfig[TestConfig]{Config: cfg}
	mockFile := &MockFile{exists: false}
	g.file = mockFile

	// Test initialize (should save since file doesn't exist)
	if err := g.initialize(); err != nil {
		t.Fatalf("initialize failed: %v", err)
	}
	if mockFile.saved.Name != "Alice" {
		t.Errorf("expected saved name Alice, got %s", mockFile.saved.Name)
	}

	// Test loading config
	mockFile.exists = true
	mockFile.loaded = &TestConfig{Name: "Bob", Age: 40}
	if err := g.initialize(); err != nil {
		t.Fatalf("initialize (load) failed: %v", err)
	}
	if g.Config.Name != "Bob" {
		t.Errorf("expected loaded name Bob, got %s", g.Config.Name)
	}
}

// Test validators
func TestValidators(t *testing.T) {
	cfg := TestConfig{Name: "Alice", Age: 30}
	g := &Gonfig[TestConfig]{Config: cfg}
	called := false
	g.AddValidator(func(c TestConfig) error {
		called = true
		if c.Age < 0 {
			return errors.New("age must be positive")
		}
		return nil
	})

	// Should pass
	if err := g.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
	if !called {
		t.Error("validator not called")
	}

	// Should fail
	if err := g.Update(TestConfig{Name: "Alice", Age: -1}); err == nil {
		t.Error("expected validation error for negative age")
	}
}

// Test Save and Load
func TestSaveLoad(t *testing.T) {
	cfg := TestConfig{Name: "Alice", Age: 30}
	g := &Gonfig[TestConfig]{Config: cfg}
	mockFile := &MockFile{}
	g.file = mockFile

	if err := g.Save(); err != nil {
		t.Errorf("Save failed: %v", err)
	}
	if mockFile.saved.Name != "Alice" {
		t.Errorf("expected saved name Alice, got %s", mockFile.saved.Name)
	}

	mockFile.loaded = &TestConfig{Name: "Carol", Age: 25}
	if err := g.Load(); err != nil {
		t.Errorf("Load failed: %v", err)
	}
	if g.Config.Name != "Carol" {
		t.Errorf("expected loaded name Carol, got %s", g.Config.Name)
	}
}
