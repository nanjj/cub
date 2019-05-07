package sca_test

import "testing"
import "github.com/nanjj/cub/sca"

func TestFromEnv(t *testing.T) {
	c := sca.Config{}
	if err := c.FromEnv(); err != nil {
		t.Fatal(err)
	}
	if c.RunnerIp == "" || c.RunnerName == "" || c.RunnerListen == "" {
		t.Fatal(c)
	}
}
