package main

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
)

func GetString(key, defaultValue string) (value string) {
	if value = viper.GetString(key); value == "" {
		value = defaultValue
	}
	return
}

func SetString(key, value string) {
	viper.Set(key, value)
}

// firstIP returns first ip address
func firstIP() (ipaddr string) {
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
	ip = GetString(RunnerIp, "")
	if ip == "" {
		ip = firstIP()
		SetString(RunnerIp, ip)
	}
	return
}

func runnerListen() (runner string) {
	runner = GetString(RunnerListen, "")
	if runner == "" {
		runner = fmt.Sprintf("tcp://%s:49160", runnerIp())
		SetString(RunnerListen, runner)
	}
	return
}

func leaderListen() (leader string) {
	leader = GetString(LeaderListen, "")
	return
}

func runnerName() (name string) {
	name = GetString(RunnerName, "")
	if name == "" {
		name = firstIP()
		SetString(RunnerName, name)
	}
	return
}
