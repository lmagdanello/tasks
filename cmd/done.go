package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Marks a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var idList []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the argument:", arg)
			} else {
				idList = append(idList, id)
			}
		}
		fmt.Println(idList)
	},
}

func init() {
	RootCmd.AddCommand(doneCmd)
}
