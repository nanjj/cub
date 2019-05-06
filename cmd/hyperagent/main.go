package main

import (
	"github.com/nanjj/cub/sca"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:  "hyperagent",
		RunE: RunHyperAgentE,
	}
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		RootCmd.Println(err)
	}
}

func RunHyperAgentE(cmd *cobra.Command, args []string) (err error) {
	cfg := &sca.Config{}
	cfg.FromEnv()
	r, err := sca.NewRunner(cfg)
	if err != nil {
		return
	}
	err = r.Wait()
	return
}
