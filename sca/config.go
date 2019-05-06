package sca

import (
	"fmt"
	"net"
	"os"
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
}

func (cfg *Config) FromEnv() (err error) {
	cfg.RunnerIp = runnerIp()
	cfg.RunnerName = runnerName()
	cfg.RunnerListen = runnerListen()
	cfg.LeaderListen = leaderListen()
	return
}

func init() {
	os.Setenv(envJaegerSamplerType, "const")
	os.Setenv(envJaegerSamplerParam, "1")
	if os.Getenv(envJaegerReporterMaxQueueSize) == "" {
		os.Setenv(envJaegerReporterMaxQueueSize, "64")
	}
	if os.Getenv(envJaegerReporterFlushInterval) == "" {
		os.Setenv(envJaegerReporterFlushInterval, "10s")
	}
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
		if hostname, err := os.Hostname(); err != nil {
			panic(err)
		} else {
			name = hostname
		}
		os.Setenv(envRunnerName, name)
	}
	return
}
