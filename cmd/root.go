package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "tasks is a CLI TODO list manager",
}
