package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

var searchFilter string

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search tasks by text",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			searchFilter = strings.Join(args, " ")
		}
		oldFilter := listFilter
		listFilter = searchFilter
		err := runList(cmd, args)
		listFilter = oldFilter
		return err
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&listDone, "done", false, "Include archived tasks")
	searchCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "Filter by priority")
}
