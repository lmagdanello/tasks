package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lmagdanello/tasks/db"
	"github.com/mitchellh/go-homedir"
)

func whichPath(bucketName []byte) (string, string) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".tasks")
	dbPath := filepath.Join(path, string(bucketName)+".db")
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Bucket \"%s\" doesn't exist! \nCreate a new one with 'init'\n", string(bucketName))
		os.Exit(1)
	}

	return path, dbPath
}

func initBucket(bucketName []byte) error {
	path, dbPath := whichPath(bucketName)
	if err := os.Mkdir(path, 0777); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err := db.Create(dbPath, bucketName)
	if err != nil {
		log.Println("Something went wrong with your bucket:", err.Error())
	}

	return err
}

func chooseBucket(bucketName []byte) []byte {
	_, dbPath := whichPath(bucketName)
	_, err := db.Init(dbPath)
	if err != nil {
		log.Println("Something went wrong initializing bucket: ", err)
	}

	return bucketName
}
