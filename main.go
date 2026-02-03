// This file enables 'go run .' from the project root.
// The main binary entry point is in cmd/tasks/.

package main

import (
	"log"
	"os"

	"github.com/lmagdanello/tasks/internal/cli"
	"github.com/lmagdanello/tasks/internal/tui"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 || (len(args) == 1 && args[0] == "--tui") {
		app := tui.New()
		if err := app.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}

	if err := cli.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
