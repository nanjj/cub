package sca_test

import (
	"encoding/binary"
	"testing"

	"github.com/nanjj/cub/sca"
)

func TestFirstIP(t *testing.T) {
	if ip := sca.FirstIPV4(); len(ip) != 4 {
		t.Fatalf("%0x", ip)
	} else {
		port := uint16(49160)
		b := make([]byte, 6)
		copy(b, ip)
		binary.BigEndian.PutUint16(b[4:], port)
		t.Logf("%x", b)
	}
}
