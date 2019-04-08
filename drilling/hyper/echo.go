package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"github.com/ugorji/go/codec"
	"nanjj.github.io/nanomsg/mangos"
	"nanjj.github.io/nanomsg/mangos/protocol/sub"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	_ "nanomsg.org/go/mangos/v2/transport/tcp"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	captainCmd := &cobra.Command{
		Use:  "captain",
		RunE: CaptainRunE,
	}

	soldierCmd := &cobra.Command{
		Use:  "soldier",
		RunE: SoldierRunE,
	}
	RootCmd.AddCommand(captainCmd, soldierCmd)
}

func GetCaptainListen() string {
	return GetString("captain.listen", "tcp://127.0.0.1:9999")
}

type TimeMessage struct {
	Time time.Time `codec:"time"`
}

func CaptainRunE(cmd *cobra.Command, args []string) (err error) {
	sock, err := pub.NewSocket()
	if err != nil {
		return
	}
	if err = sock.Listen(GetCaptainListen()); err != nil {
		return
	}
	out := make([]byte, 0, 128)
	cborHand := &codec.CborHandle{}
	for {
		// sleep random millseconds
		// time.Sleep(time.Millisecond)
		out = out[:]
		enc := codec.NewEncoderBytes(&out, cborHand)
		now := time.Now().Round(0)
		if err = enc.Encode(&now); err != nil {
			return
		}
		if err = sock.Send(out); err != nil {
			return
		}
	}
	return
}

func SoldierRunE(cmd *cobra.Command, args []string) (err error) {
	sock, err := sub.NewSocket()
	if err != nil {
		return
	}
	if err = sock.Dial(GetCaptainListen()); err != nil {
		return
	}
	sock.SetOption(mangos.OptionSubscribe, []byte(""))
	var out []byte
	cborHand := &codec.CborHandle{}
	count := 0
	startTime := time.Now()
	max := 100000
	for {
		out, err = sock.Recv()
		if err != nil {
			return
		}
		count++
		dec := codec.NewDecoderBytes(out, cborHand)
		rNow := time.Time{}
		if err = dec.Decode(&rNow); err != nil {
			return
		}
		if count >= max {
			now := time.Now().Round(0)
			d := now.Sub(startTime)
			speed := max * int(time.Millisecond) / int(d)
			fmt.Printf("%d k/s (%d in %v)\n", speed, max, d)
			count = 0
			startTime = now
		}
	}
	return
}
