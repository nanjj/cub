package drilling

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"os/exec"
	"runtime"

	"github.com/uber/jaeger-client-go/thrift"
	"github.com/uber/jaeger-client-go/thrift-gen/jaeger"
	"github.com/uber/jaeger-client-go/utils"
)

func Open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	c := exec.Command(cmd, args...)
	if err := c.Start(); err != nil {
		return err
	}
	if err := c.Wait(); err != nil {
		return err
	}
	return nil
}

func AgentProxy(listen string, ch chan []byte, closers chan io.Closer) {
	pc, err := net.ListenPacket("udp", listen)
	if err != nil {
		return
	}
	closers <- pc
	go func() {
		for {
			b := make([]byte, utils.UDPPacketMaxLength)
			n, _, err := pc.ReadFrom(b)
			if err != nil {
				return
			}
			if n == 0 {
				continue
			}
			b = b[0:n]
			ch <- b
		}
	}()
}

func NewUDPWriter(addr string) (w io.WriteCloser, err error) {
	destAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}

	connUDP, err := net.DialUDP(destAddr.Network(), nil, destAddr)
	if err != nil {
		return
	}

	if err = connUDP.SetWriteBuffer(utils.UDPPacketMaxLength); err != nil {
		return
	}
	w = connUDP
	return
}

type EmitBatch struct {
	Name   string                     `json:"name"`
	Seqid  int32                      `json:"seqid"`
	Typeid thrift.TMessageType        `json:"typeid"`
	Args   *jaeger.AgentEmitBatchArgs `json:"args"`
}

func (a *EmitBatch) Encode() (b []byte, err error) {
	trans := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(b),
	}
	oprot := thrift.NewTCompactProtocolFactory().GetProtocol(trans)
	if err = oprot.WriteMessageBegin(a.Name, a.Typeid, a.Seqid); err != nil {
		return
	}
	if a.Args != nil {
		if err = a.Args.Write(oprot); err != nil {
			return
		}
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	if err = oprot.Flush(); err != nil {
		return
	}
	b = trans.Buffer.Bytes()
	return
}

func (a *EmitBatch) Decode(b []byte) (err error) {
	trans := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(b),
	}
	iprot := thrift.NewTCompactProtocolFactory().GetProtocol(trans)
	a.Name, a.Typeid, a.Seqid, err = iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if a.Args == nil {
		a.Args = &jaeger.AgentEmitBatchArgs{}
	}
	if err = a.Args.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	return
}

func (a EmitBatch) String() (s string) {
	if b, err := json.MarshalIndent(&a, "", "  "); err == nil {
		s = string(b)
	}
	return
}
