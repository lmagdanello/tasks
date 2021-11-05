package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/lmagdanello/tasks/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, _ := cmd.Flags().GetString("bucket")
		if bucket != "" {
			bucketName := chooseBucket([]byte(bucket))
			tasks, err := db.AllTasks(bucketName)
			if err != nil {
				log.Println("Something went wrong: ", err)
				os.Exit(1)
			}
			if len(tasks) == 0 {
				fmt.Println("You have no tasks!")
				return
			}
			fmt.Println("List tasks:")
			for i, task := range tasks {
				fmt.Printf("%d. %s\n", i+1, task.Value)
			}
		} else {
			log.Println("Choose a bucket dumbass!")
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().String("bucket", "", "Choose a bucket to work")
}
