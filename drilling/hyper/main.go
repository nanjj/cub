package main

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	RootCmd = &cobra.Command{
		Use:  "hyper",
		RunE: RunHyper,
	}
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HYPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		RootCmd.Println(err)
	}
}
