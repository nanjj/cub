package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/uber/jaeger-client-go/config"
)

const (
	envRunnerListen                = "HYPERAGENT_RUNNER_LISTEN"
	envRunnerName                  = "HYPERAGENT_RUNNER_NAME"
	envRunnerIp                    = "HYPERAGENT_RUNNER_IP"
	envLeaderListen                = "HYPERAGENT_LEADER_LISTEN"
	envJaegerServiceName           = "JAEGER_SERVICE_NAME"
	envJaegerTags                  = "JAEGER_TAGS"
	envJaegerSamplerType           = "JAEGER_SAMPLER_TYPE"
	envJaegerSamplerParam          = "JAEGER_SAMPLER_PARAM"
	envJaegerReporterMaxQueueSize  = "JAEGER_REPORTER_MAX_QUEUE_SIZE"
	envJaegerReporterFlushInterval = "JAEGER_REPORTER_FLUSH_INTERVAL"
)

type Config struct {
	RunnerName   string
	RunnerListen string
	LeaderListen string
	RunnerIp     string
	Jaeger       *config.Configuration
}

func (cfg *Config) init() (err error) {
	cfg.RunnerIp = runnerIp()
	cfg.RunnerName = runnerName()
	cfg.RunnerListen = runnerListen()
	cfg.LeaderListen = leaderListen()
	err = cfg.initJaeger()
	return
}

func (cfg *Config) initJaeger() (err error) {
	os.Setenv(envJaegerServiceName, "Hyperagent")
	os.Setenv(envJaegerSamplerType, "const")
	os.Setenv(envJaegerSamplerParam, "1")
	if os.Getenv(envJaegerReporterMaxQueueSize) == "" {
		os.Setenv(envJaegerReporterMaxQueueSize, "64")
	}
	if os.Getenv(envJaegerReporterFlushInterval) == "" {
		os.Setenv(envJaegerReporterFlushInterval, "10s")
	}
	if tags, runnerTag := os.Getenv(envJaegerTags),
		fmt.Sprintf("runner=%s", cfg.RunnerName); tags == "" {
		os.Setenv(envJaegerTags, runnerTag)
	} else if strings.Contains(tags, runnerTag) {
		os.Setenv(envJaegerTags, fmt.Sprintf("%s,%s", tags, runnerTag))
	}
	cfg.Jaeger, err = config.FromEnv()
	return
}

// firstIp returns first ip address
func firstIp() (ipaddr string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ipaddr = ip.String()
			return
		}
	}
	return
}

func runnerIp() (ip string) {
	if ip = os.Getenv(envRunnerIp); ip == "" {
		ip = firstIp()
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
		name = runnerIp()
		os.Setenv(envRunnerName, name)
	}
	return
}
