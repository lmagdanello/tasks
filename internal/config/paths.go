package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	tasksDirName = ".tasks"
	dbExt        = ".db"
)

// TasksDir returns the base directory where buckets are stored (~/.tasks).
func TasksDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, tasksDirName), nil
}

// DBPath returns the full path to the .db file for a bucket.
func DBPath(bucket string) (string, error) {
	dir, err := TasksDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, bucket+dbExt), nil
}

// ListBuckets returns names of all existing buckets.
func ListBuckets() ([]string, error) {
	dir, err := TasksDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	buckets := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, dbExt) && len(name) > len(dbExt) {
			buckets = append(buckets, name[:len(name)-len(dbExt)])
		}
	}
	return buckets, nil
}

// BucketExists checks if the bucket exists on disk.
func BucketExists(bucket string) (bool, error) {
	path, err := DBPath(bucket)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
