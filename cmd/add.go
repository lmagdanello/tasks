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
		bucket, _ := cmd.Flags().GetString("bucket")
		if bucket != "" {
			bucketName := chooseBucket([]byte(bucket))
			task := strings.Join(args, " ")
			_, err := db.CreateTask(task, bucketName)
			if err != nil {
				log.Println("Something went wrong creating task:", err.Error())
				return
			}

			fmt.Printf("Added \"%s\" to your task list\n", task)
		} else {
			log.Println("Choose a bucket dumbass!")

		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().String("bucket", "", "Choose a bucket to work")

}
