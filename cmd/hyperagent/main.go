package main

import (
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
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
	cfg := &Config{}
	cfg.init()
	r, err := NewRunner(cfg)
	if err != nil {
		return
	}
	var g errgroup.Group
	g.Go(r.Run)
	err = g.Wait()
	return
}
