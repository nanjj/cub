package main

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	RootCmd = &cobra.Command{
		Use:  "hyperagent",
		RunE: RunHyperAgentE,
	}
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HYPERAGENT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		RootCmd.Println(err)
	}
}

func RunHyperAgentE(cmd *cobra.Command, args []string) (err error) {
	name := runnerName()
	listen := runnerListen()
	leader := leaderListen()
	r, err := NewRunner(name, listen, leader)
	if err != nil {
		return
	}
	var g errgroup.Group
	g.Go(r.Run)
	err = g.Wait()
	return
}
