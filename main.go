package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/lmagdanello/tasks/cmd"
	"github.com/lmagdanello/tasks/db"
	"github.com/mitchellh/go-homedir"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")
	must(db.Init(dbPath))
	must(cmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
