package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lmagdanello/tasks/db"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Marks a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var idList []int
		bucket, _ := cmd.Flags().GetString("bucket")
		if bucket != "" {
			bucketName := chooseBucket([]byte(bucket))
			for _, arg := range args {
				id, err := strconv.Atoi(arg)
				if err != nil {
					fmt.Println("Failed to parse the argument:", arg)
				} else {
					idList = append(idList, id)
				}
			}
			tasks, err := db.AllTasks(bucketName)
			if err != nil {
				log.Println("Something went wrong: ", err)
				return
			}
			for _, id := range idList {
				if id <= 0 || id > len(tasks) {
					fmt.Println("Invalid task ID: ", id)
					continue
				}
				task := tasks[id-1]
				err := db.DeleteTask(task.Key, bucketName)
				if err != nil {
					fmt.Printf("Failed to delete task: \"%d\". Error: %s\n", id, err)
				} else {
					fmt.Printf("Task \"%d\" completed\n", id)
				}
			}
		} else {
			log.Println("Choose a bucket dumbass!")
		}
	},
}

func init() {
	RootCmd.AddCommand(doneCmd)
	doneCmd.Flags().String("bucket", "", "Choose a bucket to work")
}
