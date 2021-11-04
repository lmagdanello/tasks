package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/lmagdanello/tasks/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			log.Println("Something went wrong creating task:", err.Error())
			return
		}

		fmt.Printf("Added \"%s\" to your task list\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
