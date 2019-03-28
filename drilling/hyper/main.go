package main

import "github.com/spf13/cobra"

var (
	RootCmd = &cobra.Command{
		Use:  "clerk",
		RunE: RunClerk,
	}
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		RootCmd.Println(err)
	}
}
