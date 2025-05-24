package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/tnaucoin/stringer/types"
)

type CacheFile struct {
	Hash    string                  `json:"hash"`
	Actions []types.CompositeAction `json:"actions"`
}

func SaveActionsWithHash(actions []types.CompositeAction, rootdir, filepath string) error {
	hash, err := hashDirectory(rootdir)
	if err != nil {
		return fmt.Errorf("failed to hash directory: %w", err)
	}

	cache := CacheFile{
		Hash:    hash,
		Actions: actions,
	}

	data, err := json.MarshalIndent(cache, "", "	")
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}
	return nil
}

func IsCacheValid(rootDir, cachePath string) (bool, error) {
	cache, err := LoadCache(cachePath)
	if err != nil {
		return false, err
	}
	currentHash, err := hashDirectory(rootDir)
	if err != nil {
		return false, err
	}
	return currentHash == cache.Hash, nil
}

func SaveActions(actions []types.CompositeAction, filepath string) error {
	data, err := json.MarshalIndent(actions, "", "	")
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}
	return nil
}

func LoadActions(filepath string) (*CacheFile, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var cache CacheFile
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return &cache, nil
}

func LoadCache(filepath string) (*CacheFile, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}
	var cache CacheFile
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache file: %w", err)
	}
	return &cache, nil
}

var hashDirectory = func(rootDir string) (string, error) {
	var entries []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			entries = append(entries, fmt.Sprintf("%s:%d", path, info.ModTime().UnixNano()))
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	sort.Strings(entries)
	h := sha256.New()
	for _, entry := range entries {
		h.Write([]byte(entry))
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
