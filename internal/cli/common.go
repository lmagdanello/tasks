package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lmagdanello/tasks/internal/store"
)

func requireBucket() bool {
	if bucket == "" {
		color.New(color.FgRed).Fprintln(os.Stderr, "Error: bucket is required. Use --bucket <name> or -b <name>")
		fmt.Fprintln(os.Stderr, "Create a new bucket with: tasks init --bucket <name>")
		return false
	}
	return true
}

func printBucketHeader() {
	if bucket != "" {
		color.New(color.FgGreen).Printf("Bucket: %s\n", bucket)
	}
}

// openStore opens the store for the current bucket. Prints error and returns nil if bucket not found.
func openStore() (*store.Store, error) {
	s, err := store.Open(bucket)
	if err != nil {
		if err == store.ErrBucketNotFound {
			color.New(color.FgRed).Printf("Bucket \"%s\" does not exist. Create it with: tasks init --bucket %s\n", bucket, bucket)
			return nil, nil
		}
		return nil, fmt.Errorf("opening bucket: %w", err)
	}
	return s, nil
}
