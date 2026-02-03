package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lmagdanello/tasks/internal/store"
	"github.com/spf13/cobra"
)

var (
	editPriority string
	editDue      string
)

var editCmd = &cobra.Command{
	Use:   "edit <id> [new text]",
	Short: "Edit a task's text, priority, or due date",
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

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		tasks, err := s.AllTasks()
		if err != nil {
			return fmt.Errorf("listing tasks: %w", err)
		}

		if id <= 0 || id > len(tasks) {
			return fmt.Errorf("invalid task ID: %d", id)
		}

		task := &tasks[id-1]
		newText := strings.TrimSpace(strings.Join(args[1:], " "))

		opts := make([]store.UpdateOption, 0, 3)
		if editPriority != "" {
			opts = append(opts, store.UpdatePriority(parsePriority(editPriority)))
		}
		if editDue == "none" || editDue == "-" {
			opts = append(opts, store.UpdateDue(nil))
		} else if editDue != "" {
			t, err := parseDue(editDue)
			if err != nil {
				return fmt.Errorf("invalid due date %q: %w", editDue, err)
			}
			opts = append(opts, store.UpdateDue(&t))
		}

		if newText == "" && editPriority == "" && editDue == "" {
			printBucketHeader()
			fmt.Printf("Task %d: %s\n", id, task.Text)
			fmt.Println("Use: tasks edit <id> [new text] -p <priority> -d <due>")
			return nil
		}

		textToUse := newText
		if textToUse == "" {
			textToUse = task.Text
		}
		if err := s.UpdateTask(task.Key, textToUse, opts...); err != nil {
			return fmt.Errorf("updating task: %w", err)
		}

		printBucketHeader()
		fmt.Printf("Task %d updated.\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVarP(&editPriority, "priority", "p", "", "Priority: alta, media, baixa")
	editCmd.Flags().StringVarP(&editDue, "due", "d", "", "Due date: YYYY-MM-DD or none")
}
