package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lmagdanello/tasks/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new bucket (list)",
	Long:  "Create a new bucket to organize your tasks. Each bucket is a separate list.",
	RunE: func(*cobra.Command, []string) error {
		if !requireBucket() {
			return nil
		}

		if err := store.CreateBucket(bucket); err != nil {
			if err == store.ErrBucketExists {
				color.New(color.FgYellow).Printf("Bucket \"%s\" already exists!\n", bucket)
				return nil
			}
			return fmt.Errorf("creating bucket: %w", err)
		}

		color.New(color.FgGreen).Printf("Bucket \"%s\" created successfully!\n", bucket)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
