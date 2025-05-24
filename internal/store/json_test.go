package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/tnaucoin/stringer/parser"
)

func TestSaveActionsWithHash(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	// Create test actions
	actions := []parser.CompositeAction{
		{
			Name:        "Test Action",
			Description: "A test action",
			Inputs:      map[string]any{"input1": map[string]any{"description": "Test input"}},
			Outputs:     map[string]any{"output1": map[string]any{"description": "Test output"}},
			Path:        "test/path",
		},
	}

	// Test saving actions with hash
	err := SaveActionsWithHash(actions, tmpDir, cachePath)
	if err != nil {
		t.Fatalf("SaveActionsWithHash failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatalf("Cache file was not created")
	}

	// Load and verify content
	cache, err := LoadCache(cachePath)
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	if len(cache.Actions) != 1 || cache.Actions[0].Name != "Test Action" {
		t.Errorf("Cache content doesn't match expected actions")
	}

	if cache.Hash == "" {
		t.Errorf("Hash was not generated")
	}
}

func TestIsCacheValid(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	// Create a mock cache file with a known hash
	constantHash := "mockhash123456789"
	cache := CacheFile{
		Hash: constantHash,
		Actions: []parser.CompositeAction{
			{
				Name:        "Test Action",
				Description: "A test action",
			},
		},
	}

	// Write the mock cache file directly
	data, err := json.MarshalIndent(cache, "", "\t")
	if err != nil {
		t.Fatalf("Failed to marshal mock cache: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write mock cache file: %v", err)
	}

	// Mock the hashDirectory function to return a constant hash
	originalHashDir := hashDirectory
	// Restore the original function when the test completes
	defer func() { hashDirectory = originalHashDir }()

	// Replace with mock function
	hashDirectory = func(rootDir string) (string, error) {
		return constantHash, nil
	}

	// Test cache validity - should be valid with our mock
	valid, err := IsCacheValid(tmpDir, cachePath)
	if err != nil {
		t.Fatalf("IsCacheValid failed: %v", err)
	}
	if !valid {
		t.Errorf("Cache should be valid but was reported as invalid")
	}

	// Now change the mock to return a different hash
	hashDirectory = func(rootDir string) (string, error) {
		return "differenthash", nil
	}

	// Test cache validity again - should be invalid now
	valid, err = IsCacheValid(tmpDir, cachePath)
	if err != nil {
		t.Fatalf("IsCacheValid failed after modification: %v", err)
	}
	if valid {
		t.Errorf("Cache should be invalid after directory modification but was reported as valid")
	}
}

func TestSaveAndLoadActions(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "actions.json")

	// Create test actions
	actions := []parser.CompositeAction{
		{
			Name:        "Action 1",
			Description: "First test action",
		},
		{
			Name:        "Action 2",
			Description: "Second test action",
		},
	}

	// Test SaveActions
	err := SaveActions(actions, filePath)
	if err != nil {
		t.Fatalf("SaveActions failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Actions file was not created")
	}

	// Read the file directly instead of using LoadActions
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read actions file: %v", err)
	}

	// Unmarshal directly to actions slice
	var loadedActions []parser.CompositeAction
	if err := json.Unmarshal(data, &loadedActions); err != nil {
		t.Fatalf("Failed to unmarshal actions: %v", err)
	}

	if len(loadedActions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(loadedActions))
	}

	if loadedActions[0].Name != "Action 1" || loadedActions[1].Name != "Action 2" {
		t.Errorf("Loaded actions don't match saved actions")
	}
}

func TestLoadCache(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	// Create a test cache file
	cache := CacheFile{
		Hash: "testhash123",
		Actions: []parser.CompositeAction{
			{
				Name:        "Cached Action",
				Description: "A cached action",
			},
		},
	}

	data, err := json.MarshalIndent(cache, "", "\t")
	if err != nil {
		t.Fatalf("Failed to marshal test cache: %v", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write test cache file: %v", err)
	}

	// Test LoadCache
	loadedCache, err := LoadCache(cachePath)
	if err != nil {
		t.Fatalf("LoadCache failed: %v", err)
	}

	if loadedCache.Hash != "testhash123" {
		t.Errorf("Expected hash 'testhash123', got '%s'", loadedCache.Hash)
	}

	if len(loadedCache.Actions) != 1 || loadedCache.Actions[0].Name != "Cached Action" {
		t.Errorf("Loaded cache doesn't match expected content")
	}
}

func TestLoadCacheErrors(t *testing.T) {
	// Test loading non-existent file
	_, err := LoadCache("nonexistent.json")
	if err == nil {
		t.Errorf("Expected error when loading non-existent file, got nil")
	}

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "invalid.json")

	// Create an invalid JSON file
	if err := os.WriteFile(cachePath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid cache file: %v", err)
	}

	// Test loading invalid JSON
	_, err = LoadCache(cachePath)
	if err == nil {
		t.Errorf("Expected error when loading invalid JSON, got nil")
	}
}

func TestHashDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Get initial hash
	hash1, err := hashDirectory(tmpDir)
	if err != nil {
		t.Fatalf("hashDirectory failed: %v", err)
	}

	if hash1 == "" {
		t.Errorf("Expected non-empty hash")
	}

	// Create a file and check that hash changes
	filePath := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hash2, err := hashDirectory(tmpDir)
	if err != nil {
		t.Fatalf("hashDirectory failed after adding file: %v", err)
	}

	if hash1 == hash2 {
		t.Errorf("Hash should change after adding a file")
	}

	// Modify the file and check that hash changes again
	if err := os.WriteFile(filePath, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	hash3, err := hashDirectory(tmpDir)
	if err != nil {
		t.Fatalf("hashDirectory failed after modifying file: %v", err)
	}

	if hash2 == hash3 {
		t.Errorf("Hash should change after modifying a file")
	}
}
