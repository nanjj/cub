package drilling

import (
	"os"

	"github.com/grandecola/bigqueue"
)

func Cleanup(q *bigqueue.BigQueue) (err error) {
	for !q.IsEmpty() {
		err = q.Dequeue()
		if err != nil {
			return
		}
	}
	return
}

func MakeQueueDir(name string) (err error) {
	if err = os.MkdirAll(name, 0755); err != nil {
		return
	}
	return
}
