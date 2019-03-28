package drilling

import (
	"fmt"
	"os"
	"testing"

	"github.com/grandecola/bigqueue"
	"golang.org/x/sync/errgroup"
)

func TestBigQueueUsage(t *testing.T) {
	bqd := ".test/bing_queue_usage.bqd"
	MakeQueueDir(bqd)
	q, err := bigqueue.NewBigQueue(bqd)
	if err != nil {
		t.Fatal(err)
	}
	Cleanup(q)
	for _, s := range []string{"hello", "world", "!"} {
		if err = q.Enqueue([]byte(s)); err != nil {
			t.Fatal(err)
		}
	}
	if err = q.Dequeue(); err != nil { // remove hello
		t.Fatal(err)
	}

	b, err := q.Peek()
	if err != nil { // read world
		if s := string(b); "world" != s {
			t.Fatal(s)
		}
	}

	if err = q.Dequeue(); err != nil { // remove world
		t.Fatal(err)
	}

	b, err = q.Peek() // read !
	if err != nil {
		t.Fatal(err)
	}
	if s := string(b); "!" != s {
		t.Fatal(s)
	}

	if err = q.Dequeue(); err != nil {
		t.Fatal(err)
	}
}

func TestBigQueueConcurrency(t *testing.T) {
	bqd := ".test/bing_queue_wait.bqd"
	os.MkdirAll(bqd, 0755)
	q, err := bigqueue.NewBigQueue(bqd)
	if err != nil {
		t.Fatal(err)
	}
	for !q.IsEmpty() {
		err = q.Dequeue()
		t.Fatal(err)
	}
	msgs := []string{}
	var g errgroup.Group
	for i := 0; i < 10; i++ {
		i := i
		g.Go(func() error {
			return q.Enqueue([]byte(fmt.Sprintf("message%02d", i)))
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	for !q.IsEmpty() {
		b, err := q.Peek()
		if err != nil {
			t.Fatal(err)
		}
		err = q.Dequeue()
		if err != nil {
			t.Fatal(err)
		}
		msgs = append(msgs, string(b))
	}
	if len(msgs) != 10 {
		t.Fatal(len(msgs), msgs)
	}
}
