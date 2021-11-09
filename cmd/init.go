package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init a new bucket",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, _ := cmd.Flags().GetString("bucket")
		if bucket != "" {
			err := initBucket([]byte(bucket))
			if err != nil {
				log.Println("Something went wrong creating your bucket:", err)
			}
			fmt.Printf("Bucket \"%s\" created with success!\n", bucket)
		} else {
			fmt.Println("Choose a bucket!!!!!")
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().String("bucket", "", "Choose your bucket name!")

}
