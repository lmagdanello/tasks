package cmd

import (
	"log"
	"path/filepath"

	"github.com/lmagdanello/tasks/db"
	"github.com/mitchellh/go-homedir"
)

func chooseBucket(bucketName []byte) []byte {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, string(bucketName)+".db")
	err := db.Init(dbPath, bucketName)
	if err != nil {
		log.Println("Something went wrong with your bucket:", err.Error())
	}

	return bucketName
}

// function removeBucket
// remove bucket from dbPath (remove .db files)
//
//
