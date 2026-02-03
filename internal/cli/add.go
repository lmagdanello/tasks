package cli

import (
	"fmt"
	"strings"

	"github.com/lmagdanello/tasks/internal/store"
	"github.com/spf13/cobra"
)

var (
	addPriority string
	addDue      string
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task to your list",
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

		task := strings.Join(args, " ")

		opts := make([]store.CreateOption, 0, 2)
		opts = append(opts, store.WithPriority(parsePriority(addPriority)))
		if addDue != "" {
			t, err := parseDue(addDue)
			if err != nil {
				return fmt.Errorf("invalid due date %q: %w", addDue, err)
			}
			opts = append(opts, store.WithDue(t))
		}

		if _, err = s.CreateTask(task, opts...); err != nil {
			return fmt.Errorf("creating task: %w", err)
		}

		printBucketHeader()
		fmt.Printf("Added \"%s\" to your task list!\n", task)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addPriority, "priority", "p", "media", "Priority: alta, media, baixa")
	addCmd.Flags().StringVarP(&addDue, "due", "d", "", "Due date: YYYY-MM-DD")
}
