package main_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rt := func() int {
		os.Setenv("HYPERAGENT_RUNNER_IP", "127.0.0.1")
		return m.Run()
	}()
	os.Exit(rt)
}
