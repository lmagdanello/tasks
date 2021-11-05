package main

import (
	"log"
	"os"

	"github.com/lmagdanello/tasks/cmd"
)

func main() {
	must(cmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
