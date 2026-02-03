package cli

import (
	"fmt"

	"github.com/lmagdanello/tasks/internal/store"
	"github.com/spf13/cobra"
)

var (
	listDone     bool
	listFilter   string
	listPriority string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&listDone, "done", false, "List archived tasks")
	listCmd.Flags().StringVarP(&listFilter, "filter", "f", "", "Filter by text")
	listCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "Filter by priority: alta, media, baixa")
}

func runList(_ *cobra.Command, _ []string) error {
	if !requireBucket() {
		return nil
	}

	s, err := openStore()
	if err != nil || s == nil {
		return err
	}
	defer s.Close()

	opts := make([]store.ListOption, 0, 2)
	if listFilter != "" {
		opts = append(opts, store.FilterQuery(listFilter))
	}
	if listPriority != "" {
		opts = append(opts, store.FilterPriority(parsePriority(listPriority)))
	}

	var tasks []store.Task
	if listDone {
		tasks, err = s.DoneTasks(opts...)
	} else {
		tasks, err = s.AllTasks(opts...)
	}
	if err != nil {
		return fmt.Errorf("listing tasks: %w", err)
	}

	printBucketHeader()
	if listDone {
		fmt.Println("Archived:")
	} else {
		fmt.Println("List of tasks:")
	}
	if len(tasks) == 0 {
		fmt.Println("  (none)")
		return nil
	}

	for i := range tasks {
		line := formatTaskLine(i+1, &tasks[i])
		fmt.Println(line)
	}
	return nil
}

func formatTaskLine(idx int, t *store.Task) string {
	line := fmt.Sprintf("  %d. %s", idx, t.Text)
	if t.Priority != store.PriorityMedium {
		line += fmt.Sprintf(" [%s]", t.Priority)
	}
	if t.Due != nil {
		line += fmt.Sprintf(" (due: %s)", t.Due.Format("2006-01-02"))
	}
	return line
}
