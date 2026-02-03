package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doneDelete bool

var doneCmd = &cobra.Command{
	Use:   "done [task ids...]",
	Short: "Archive task(s) as complete (or delete with --delete)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		if !requireBucket() {
			return nil
		}

		s, err := openStore()
		if err != nil || s == nil {
			return err
		}
		defer s.Close()

		tasks, err := s.AllTasks()
		if err != nil {
			return fmt.Errorf("listing tasks: %w", err)
		}

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				printBucketHeader()
				fmt.Printf("Invalid task ID: %s\n", arg)
				continue
			}

			if id <= 0 || id > len(tasks) {
				printBucketHeader()
				fmt.Printf("Invalid task ID: %d\n", id)
				continue
			}

			task := &tasks[id-1]
			if doneDelete {
				if err := s.DeleteTask(task.Key); err != nil {
					return fmt.Errorf("deleting task %d: %w", id, err)
				}
				printBucketHeader()
				fmt.Printf("Task \"%d\" deleted.\n", id)
			} else {
				if err := s.ArchiveTask(task.Key); err != nil {
					return fmt.Errorf("archiving task %d: %w", id, err)
				}
				printBucketHeader()
				fmt.Printf("Task \"%d\" archived.\n", id)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
	doneCmd.Flags().BoolVar(&doneDelete, "delete", false, "Permanently delete instead of archive")
}
