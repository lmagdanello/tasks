package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lmagdanello/tasks/db"
	"github.com/mitchellh/go-homedir"
)

func whichPath(bucketName []byte) string {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".tasks")
	if err := os.Mkdir(path, 0777); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	dbPath := []string{string(path), "/", string(bucketName), ".db"}
	dbPath_joined := strings.Join(dbPath, "")
	return dbPath_joined
}

func initBucket(bucketName []byte) error {
	var err error
	dbPath := whichPath(bucketName)
	if err != nil {
		log.Fatal("Error: ", err)
		os.Exit(1)
	}

	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		db.Create(dbPath, bucketName)
	} else {
		fmt.Printf("Bucket '%s' already exists!\n", string(bucketName))
		os.Exit(1)
	}

	return err
}

func chooseBucket(bucketName []byte) []byte {
	dbPath := whichPath(bucketName)
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Bucket '%s' does not exist! Create a new one with init!\n", string(bucketName))
		os.Exit(1)
	}
	db.Init(dbPath)
	return bucketName
}
