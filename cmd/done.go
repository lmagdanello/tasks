package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Marks a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("done called")
	},
}

func init() {
	RootCmd.AddCommand(doneCmd)
}
