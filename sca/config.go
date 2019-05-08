package sca

import (
	"fmt"
	"os"
)

const (
	envRunnerListen = "HYPERAGENT_RUNNER_LISTEN"
	envRunnerName   = "HYPERAGENT_RUNNER_NAME"
	envRunnerIp     = "HYPERAGENT_RUNNER_IP"
	envLeaderListen = "HYPERAGENT_LEADER_LISTEN"
)

type Config struct {
	RunnerName   string
	RunnerListen string
	LeaderListen string
	RunnerIp     string
}

func (cfg *Config) FromEnv() (err error) {
	cfg.RunnerIp = runnerIp()
	cfg.RunnerName = runnerName()
	cfg.RunnerListen = runnerListen()
	cfg.LeaderListen = leaderListen()
	return
}

func runnerIp() (ip string) {
	if ip = os.Getenv(envRunnerIp); ip == "" {
		ip = FirstIPV4().String()
		os.Setenv(envRunnerIp, ip)
	}
	return
}

func runnerListen() (runner string) {
	if runner = os.Getenv(envRunnerListen); runner == "" {
		runner = fmt.Sprintf("tcp://%s:49160", runnerIp())
		os.Setenv(envRunnerListen, runner)
	}
	return
}

func leaderListen() (leader string) {
	leader = os.Getenv(envLeaderListen)
	return
}

func runnerName() (name string) {
	if name = os.Getenv(envRunnerName); name == "" {
		if hostname, err := os.Hostname(); err != nil {
			panic(err)
		} else {
			name = hostname
		}
		os.Setenv(envRunnerName, name)
	}
	return
}
