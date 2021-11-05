package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/lmagdanello/tasks/db"
	"github.com/mitchellh/go-homedir"
)

// chooseBucket already init the bucket (*.db inside HOMEDIR/.tasks)
func chooseBucket(bucketName []byte) []byte {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".tasks")
	dbPath := filepath.Join(path, string(bucketName)+".db")
	if err := os.Mkdir(path, 0777); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err := db.Init(dbPath, bucketName)
	if err != nil {
		log.Println("Something went wrong with your bucket:", err.Error())
	}

	return bucketName
}

// function removeBucket
// remove bucket from (AKA remove .db files)
//
//
