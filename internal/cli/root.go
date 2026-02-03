package cli

import (
	"github.com/spf13/cobra"
)

var bucket string

var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "tasks is a CLI TODO list manager",
	Long: `tasks is a CLI TODO list manager based on John Calhoun's Gophercises.

Tasks are organized in buckets (lists). Use 'tasks init --bucket <name>' to create
a new bucket, then add and manage tasks within it.

Run without arguments for the interactive TUI.`,
}

// Execute runs the CLI. Returns error if command fails.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "", "Bucket (list) to work with")
}
