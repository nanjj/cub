package drilling

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"encoding/json"

	"golang.org/x/sync/errgroup"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/bus"
	"nanomsg.org/go/mangos/v2/protocol/pair"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/protocol/pull"
	"nanomsg.org/go/mangos/v2/protocol/push"
	"nanomsg.org/go/mangos/v2/protocol/rep"
	"nanomsg.org/go/mangos/v2/protocol/req"
	"nanomsg.org/go/mangos/v2/protocol/respondent"
	"nanomsg.org/go/mangos/v2/protocol/star"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	"nanomsg.org/go/mangos/v2/protocol/surveyor"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
	_ "nanomsg.org/go/mangos/v2/transport/tcp"
)

// Request/reply pattern
func TestReqRep(t *testing.T) {
	const (
		listen = "tcp://127.0.0.1:40899"
	)
	type Message struct {
		Index    int           `json:"index"`
		Duration time.Duration `json:"duration"`
		Input    string        `json:"input"`
		Output   string        `json:"output"`
	}
	start := func() mangos.Socket {
		var (
			b   []byte
			err error
		)
		sock, err := rep.NewSocket()
		if err != nil {
			t.Fatal(err)
		}
		if err = sock.Listen(listen); err != nil {
			t.Fatal(err)
		}
		go func() {
			for {
				b, err = sock.Recv()
				if err != nil {
					return
				}
				msg := &Message{}
				if err = json.Unmarshal(b, msg); err != nil {
					t.Fatal(err)
				}
				msg.Output = msg.Input
				msg.Index++
				if b, err = json.Marshal(msg); err != nil {
					t.Fatal(err)
				}
				if msg.Duration != 0 {
					time.Sleep(msg.Duration)
				}
				if err = sock.Send(b); err != nil {
					return
				}
			}
		}()
		return sock
	}
	server := start()
	defer func() {
		if err := server.Close(); err != nil {
			t.Log(err)
		}
	}()
	sock, err := req.NewSocket()
	if err != nil {
		t.Fatal(err)
	}
	if err = sock.Dial(listen); err != nil {
		t.Fatal(err)
	}
	var (
		b []byte
	)
	sendAndRecv := func(sock mangos.Socket, i int) (reply *Message, err error) {
		msg := &Message{Index: i, Input: "input", Duration: time.Microsecond * time.Duration(i)}
		if b, err = json.Marshal(msg); err != nil {
			return
		}
		if err = sock.Send(b); err != nil {
			return
		}
		reply = &Message{}
		if b, err = sock.Recv(); err != nil {
			return
		}
		if err = json.Unmarshal(b, reply); err != nil {
			return
		}
		return
	}
	// send and receive 10 times
	for i := 0; i < 10; i++ {
		msg, err := sendAndRecv(sock, i)
		if err != nil {
			t.Fatal(msg, err)
		}
		if msg.Output == "" || msg.Output != msg.Input || msg.Index != i+1 {
			t.Fatal(i, msg)
		}
	}

	// send 10 times first
	for i := 0; i < 10; i++ {
		msg := &Message{Index: i, Input: "input"}
		if b, err = json.Marshal(msg); err != nil {
			t.Fatal(i, err)
		}
		if err = sock.Send(b); err != nil {
			t.Fatal(i, err)
		}
	}
	{ // only the 10th message being received
		b, err = sock.Recv()
		msg := &Message{}
		if err = json.Unmarshal(b, msg); err != nil {
			t.Fatal(err)
		}
		if msg.Output == "" || msg.Output != msg.Input || msg.Index != 10 {
			t.Log(msg)
		}
	}
	if err = sock.Close(); err != nil {
		t.Fatal(err)
	}

	run := func(n int) (err error) {
		sock, err := req.NewSocket()
		if err != nil {
			return
		}
		if err = sock.Dial(listen); err != nil {
			return
		}
		for i := 0; i < 10; i++ {
			var msg *Message
			msg, err = sendAndRecv(sock, i+n*10)
			if err != nil {
				return
			}
			if msg.Input == "" || msg.Input != msg.Output || msg.Index != i+1+n*10 {
				err = fmt.Errorf("%d:%v", i, msg)
				return
			}
			t.Logf("client %d received message %v", n, msg)
		}
		err = sock.Close()
		return
	}

	rg := errgroup.Group{}
	rg.Go(func() error { return run(1) })
	rg.Go(func() error { return run(2) })
	if err = rg.Wait(); err != nil {
		t.Fatal(err)
	}
}

// pair pattern test
func TestPair(t *testing.T) {
	const (
		step2addr = "inproc://step2"
		step3addr = "inproc://step3"
	)
	closes := make(chan func() error, 10)
	connect := func(addr string) (sock mangos.Socket, err error) {
		if sock, err = pair.NewSocket(); err != nil {
			return
		}
		err = sock.Dial(addr)
		return
	}
	notify := func(sock mangos.Socket) (err error) {
		err = sock.Send([]byte("ready"))
		return
	}

	listen := func(addr string) (sock mangos.Socket, err error) {
		if sock, err = pair.NewSocket(); err != nil {
			return
		}
		err = sock.Listen(addr)
		return
	}

	wait := func(sock mangos.Socket) (err error) {
		var (
			b []byte
		)
		if b, err = sock.Recv(); err != nil {
			return
		}
		if "ready" != string(b) {
			err = fmt.Errorf("Not ready")
		}
		return
	}

	shutdown := func(sock mangos.Socket) (err error) {
		closes <- sock.Close
		return nil
	}

	g := errgroup.Group{}
	step1 := func() (err error) {
		t.Log("enter step 1")
		t.Log("step 1 is ready")
		sock, err := connect(step2addr)
		defer shutdown(sock)
		err = notify(sock)
		return
	}

	step2 := func() (err error) {
		t.Log("enter step2")
		sock, err := listen(step2addr)
		if err != nil {
			return
		}
		defer shutdown(sock)
		g.Go(step1)
		if err = wait(sock); err != nil {
			return err
		}
		t.Log("step2 is ready")
		sock, err = connect(step3addr)
		if err != nil {
			return
		}
		err = notify(sock)
		defer shutdown(sock)
		return
	}

	step3 := func() (err error) {
		t.Log("enter step3")
		sock, err := listen(step3addr)
		if err != nil {
			return
		}
		g.Go(step2)
		if err = wait(sock); err != nil {
			return
		}
		t.Log("step3 is ready")
		return
	}
	g.Go(step3)
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	for len(closes) > 0 {
		(<-closes)()
	}
}

// One-way data distribution:
// http://zguide.zeromq.org/page:all#Getting-the-Message-Out
// #+begin_src artist
//
//                           +-------------+
//                           |    pub      |
//         +----------------->  (listen)   <-----+
//         |                 +-------------+     |
//         |                     ^               |
//         |                     |               |
//    +----+----+       +--------+--+      +-----+----+
//    |  sub1   |       |   sub2    |      |  sub3    |
//    | (dial)  |       |  (dial)   |      | (dial)   |
//    +---------+       +-----------+      +----------+
//
//   Note:
//   1. pub listen and send message
//   2. sub1  sub2, sub3 dial to pub and receive message
// #+end_src
//
func TestWeatherUpdates(t *testing.T) {
	const (
		addr = "tcp://127.0.0.1:59999"
	)

	rand.Seed(time.Now().UnixNano())

	type WheatherUpdateMessage struct {
		Zipcode     int `json:"zipcode"`
		Temperature int `json:"temperature"`
		Humidity    int `json:"humidity"`
	}

	newMessage := func() (m *WheatherUpdateMessage) {
		m = &WheatherUpdateMessage{
			Zipcode:     (rand.Intn(10)+1)*100000 + rand.Intn(10000),
			Temperature: rand.Intn(130) - 80,
			Humidity:    rand.Intn(80) + 10,
		}
		return
	}
	closes := make(chan func() error, 1024)
	defer func() {
		max := len(closes)
		for i := 0; i < max; i++ {
			f := <-closes
			f()
		}
	}()
	bind := func() (sock mangos.Socket, err error) {
		if sock, err = pub.NewSocket(); err != nil {
			return
		}
		if err = sock.Listen(addr); err == nil {
			closes <- sock.Close
		}
		return
	}

	publish := func(sock mangos.Socket) (err error) {
		msg := newMessage()
		b, err := json.Marshal(msg)
		if err != nil {
			return
		}
		err = sock.Send(b)
		return
	}

	connect := func() (sock mangos.Socket, err error) {
		if sock, err = sub.NewSocket(); err != nil {
			return
		}
		if err = sock.Dial(addr); err != nil {
			return
		}
		sock.SetOption(mangos.OptionSubscribe, []byte(""))
		closes <- sock.Close
		return
	}

	dial := func(sock mangos.Socket) (msg *WheatherUpdateMessage, err error) {
		var b []byte
		if b, err = sock.Recv(); err != nil {
			return
		}
		msg = &WheatherUpdateMessage{}
		err = json.Unmarshal(b, msg)
		closes <- sock.Close
		return
	}

	msgs := make(chan *WheatherUpdateMessage, 1000)
	ready := make(chan bool, 1000)
	subscribe := func() (err error) {
		sock, err := connect()
		if err != nil {
			return
		}
		ready <- true
		msg, err := dial(sock)
		if err != nil {
			t.Fatal(err)
		}
		msgs <- msg
		return
	}
	// start publisher
	publisher, err := bind()
	if err != nil {
		t.Fatal(err)
	}
	g := errgroup.Group{}
	// start subscriber 1
	g.Go(subscribe)
	// start subscriber 2
	g.Go(subscribe)
	// start subscriber 3
	g.Go(subscribe)
	// wait subscribe ready
	for i := 0; i < 3; i++ {
		<-ready
	}
	if err := publish(publisher); err != nil {
		t.Fatal(err)
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 3 {
		t.Fatal(len(msgs))
	}
	// check result
	msg1 := <-msgs
	msg2 := <-msgs
	msg3 := <-msgs
	if !reflect.DeepEqual(msg1, msg2) || !reflect.DeepEqual(msg1, msg3) {
		t.Fatal(msg1, msg2, msg3)
	}
}

// Divide and conquer:
// http://zguide.zeromq.org/page:all#Divide-and-Conquer
//
// if you running it under macos, you may hit the maxfiles limit,
// change it as below:
//
// 1. create file =/Library/LaunchDaemons/limit.maxfiles.plist= with
//    below content:
//    #+begin_src xml
//      <?xml version="1.0" encoding="UTF-8"?>
//      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
//       "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
//      <plist version="1.0">
//       <dict>
//       <key>Label</key>
//       <string>limit.maxfiles</string>
//       <key>ProgramArguments</key>
//       <array>
//       <string>launchctl</string>
//       <string>limit</string>
//       <string>maxfiles</string>
//       <string>64000</string>
//       <string>524288</string>
//       </array>
//       <key>RunAtLoad</key>
//       <true/>
//       <key>ServiceIPC</key>
//       <false/>
//       </dict>
//      </plist>
//    #+end_src
// 2. =sudo launchctl load -w /Library/LaunchDaemons/limit.maxfiles.plist=
//
func TestDivideAndConquer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	defer Do()
	const (
		senderAddr = "tcp://127.0.0.1:55557"
		sinkerAddr = "tcp://127.0.0.1:55558"
		maxWorkers = 100
	)
	type Task struct {
		Id       int           `json:"id"`
		Name     string        `json:"name"`
		Duration time.Duration `json:"duration"`
		Worker   string        `json:"worker"`
	}

	tLogTask := func(task *Task) {
		tLogf("%02d %s %s %v\n", task.Id, task.Name, task.Worker, task.Duration)
	}

	newSender := func() (sock mangos.Socket, err error) {
		sock, err = push.NewSocket()
		if err != nil {
			tLog(err)
			return
		}
		Defer(sock.Close)
		err = sock.Listen(senderAddr)
		if err != nil {
			tLog(err)
		}
		return
	}

	send := func(sock mangos.Socket, task *Task) (err error) {
		var (
			b []byte
		)
		if b, err = json.Marshal(task); err != nil {
			return
		}
		err = sock.Send(b)
		return
	}

	dial := func(sock mangos.Socket) (task *Task, err error) {
		var (
			b []byte
		)
		if b, err = sock.Recv(); err != nil {
			return
		}
		task = &Task{}
		err = json.Unmarshal(b, task)
		return
	}

	newSleepTask := func(id int) *Task {
		return &Task{
			Id:       id,
			Name:     "sleep",
			Duration: time.Millisecond * time.Duration(rand.Intn(10)+1),
		}
	}

	newDoneTask := func(id int) *Task {
		return &Task{
			Id:   id,
			Name: "done",
		}
	}

	work := func(idx int) (err error) {
		var (
			tasks   mangos.Socket
			results mangos.Socket
			name    = fmt.Sprintf("worker-%02d", idx)
		)
		tasks, err = pull.NewSocket()
		if err != nil {
			tLog(name, err)
			t.Fatal(name, err)
			return
		}
		Defer(tasks.Close)
		err = tasks.Dial(senderAddr)
		if err != nil {
			tLog(name, err)
			t.Fatal(name, err)
			return
		}
		results, err = push.NewSocket()
		if err != nil {
			tLog(name, err)
			t.Fatal(name, err)
			return
		}
		Defer(results.Close)
		err = results.Dial(sinkerAddr)
		if err != nil {
			tLog(name, err)
			t.Fatal(name, err)
			return
		}
		Add()
		for {
			var task *Task
			if task, err = dial(tasks); err != nil {
				tLog(name, err)
				t.Fatal(name, err)
				return
			}
			task.Worker = name
			done := false
			switch task.Name {
			case "sleep":
				time.Sleep(task.Duration)
			case "done":
				done = true
			default: // discard others
				continue
			}
			if err = send(results, task); err != nil {
				tLog(name, err)
				t.Fatal(name, err)
				return
			}

			if done {
				break
			}
		}
		return
	}
	tpool := make(chan *Task, 1024)
	sink := func() (err error) {
		var (
			sock mangos.Socket
			task *Task
		)
		sock, err = pull.NewSocket()
		if err != nil {
			return
		}
		Defer(sock.Close)
		if err = sock.Listen(sinkerAddr); err != nil {
			return
		}
		Add()
		offduties := 0
		startTime := time.Now()
		for {
			if task, err = dial(sock); err != nil {
				tLog(err)
				return
			}
			if task.Name == "done" {
				offduties++
			}
			tpool <- task
			if offduties == maxWorkers {
				break
			}
		}
		endTime := time.Now()
		tLogf("Total elapsed time:%v\n", endTime.Sub(startTime))
		return
	}

	g := errgroup.Group{}
	// start sinker
	g.Go(sink)
	// wait sinker ready
	Wait(1)

	// start sender
	sender, err := newSender()
	if err != nil {
		t.Fatal(err)
	}

	// start workers
	for i := 0; i < maxWorkers; i++ {
		i := i
		f := func() error {
			return work(i)
		}
		g.Go(f)
	}
	// send 200 sleep tasks
	for i := 0; i < 200; i++ {
		task := newSleepTask(i)
		if err = send(sender, task); err != nil {
			t.Fatal(err)
		}
	}

	// Wait 100 task being handled
	for i := 0; i < 200; i++ {
		tLogTask(<-tpool)
	}

	// send done tasks
	for i := 0; i < maxWorkers; i++ {
		task := newDoneTask(i)
		if err = send(sender, task); err != nil {
			t.Fatal(err)
		}
	}

	// WAIT all
	if err = g.Wait(); err != nil {
		t.Fatal(err)
	}

	if dones := len(tpool); dones != maxWorkers {
		t.Fatal(dones)
	}

	for i := 0; i < maxWorkers; i++ {
		tLogTask(<-tpool)
	}
}

// Bus: http://250bpm.com/blog:17
// #+begin_src artist
//               +----------+
//               |   p1     |
//      +------->| (listen) +--------+
//      |        |          |        |
//      |        +----------+        |
//      | dial                  dial |
//      |                            |
//      |                            v
// +----+-----+                +----------+
// |   p3     |     dial       |   p2     |
// | (listen) |<---------------+ (listen) |
// |          |                |          |
// +----------+                +----------+
// Note:
// p1 sends message, p2, p3 should receive the message,
// p1 should not receive the message
// #+end_src
func TestBusBasic(t *testing.T) {
	type Message struct {
		Action   string `json:"action"`
		Sender   string `json:"sender"`
		Receiver string `json:"receiver"`
	}
	var (
		g          errgroup.Group
		msgs       = make(chan *Message, 1024)
		p1, p2, p3 mangos.Socket
		err        error
		adds       = map[string]string{
			"p1": "tcp://127.0.0.1:55555",
			"p2": "tcp://127.0.0.1:55556",
			"p3": "tcp://127.0.0.1:55557",
		}
	)

	recv := func(sock mangos.Socket, name string) (err error) {
		var b []byte
		for {
			if b, err = sock.Recv(); err != nil {
				tLog(err)
				return
			}
			msg := &Message{}
			if err = json.Unmarshal(b, msg); err != nil {
				tLog(err)
				return
			}
			if msg.Action == "quit" {
				return
			}
			msg.Receiver = name
			msgs <- msg
		}
	}

	listen := func(name string) (sock mangos.Socket, err error) {
		sock, err = bus.NewSocket()
		if err != nil {
			tLog(err)
			return
		}
		addr := adds[name]
		err = sock.Listen(addr)
		if err != nil {
			tLog(err)
		}
		g.Go(func() (err error) {
			return recv(sock, name)
		})
		return
	}

	dial := func(sock mangos.Socket, name string) (err error) {
		addr := adds[name]
		if err = sock.Dial(addr); err != nil {
			tLog(err)
			return
		}
		return
	}

	send := func(sock mangos.Socket, sender, action string) (err error) {
		var (
			b []byte
		)
		if b, err = json.Marshal(&Message{
			Action: action,
			Sender: sender}); err != nil {
			tLog(err)
			return
		}
		err = sock.Send(b)
		return
	}

	if p1, err = listen("p1"); err != nil {
		tLog(err)
		t.Fatal(err)
	}

	defer p1.Close()
	if p2, err = listen("p2"); err != nil {
		tLog(err)
		t.Fatal(err)
	}

	defer p2.Close()
	if p3, err = listen("p3"); err != nil {
		tLog(err)
		t.Fatal(err)
	}
	defer p3.Close()

	// p1 dial p2
	if err = dial(p1, "p2"); err != nil {
		tLog(err)
		t.Fatal(err)
	}
	// p2 dial p3
	if err = dial(p2, "p3"); err != nil {
		tLog(err)
		t.Fatal(err)
	}

	// p3 dial p1
	if err = dial(p3, "p1"); err != nil {
		tLog(err)
		t.Fatal(err)
	}

	tLog("now send message")
	if err = send(p1, "p1", "work"); err != nil {
		tLog(err)
		t.Fatal(err)
	}
	tLogf("now check messages received: %d\n", len(msgs))
	msg1 := <-msgs
	tLogf("Get message 1: %v\n", msg1)
	msg2 := <-msgs
	tLogf("Get message 2: %v\n", msg2)
	if l := len(msgs); l != 0 {
		for len(msgs) > 0 {
			tLog(<-msgs)
		}
		t.Fatal(l)
	}
	if msg1.Sender != "p1" || msg2.Sender != "p1" ||
		msg1.Receiver == msg2.Receiver ||
		msg1.Receiver == "p1" || msg2.Receiver == "p1" ||
		msg1.Action != "work" || msg2.Action != "work" {
		tLog(msg1, msg2)
		t.Fatal(msg1, msg2)
	}
	// quit p2, p3
	if err = send(p1, "p1", "quit"); err != nil {
		tLog(err)
		t.Fatal(err)
	}
	// quit p1
	if err = send(p2, "p2", "quit"); err != nil {
		tLog(err)
		t.Fatal(err)
	}
	if err = g.Wait(); err != nil {
		tLog(err)
		t.Fatal(err)
	}
}

// Survey (Everybody Votes)
//
// The surveyor pattern is used to send a timed survey out,
// responses are individually returned until the survey has
// expired. This pattern is useful for service discovery
// and voting algorithms.
//
// #+begin_src artist
//
//         +----------+
//         |    r1    |
//         |respondent+------------+
//         |          |            |
//         +----------+            |
//                                 |
//                             +---+-------+           +-----------+
//                             |   sur     |           |   r3      |
//                             | surveyor  +-----------+respondent |
//                             |           |           |           |
//                             +---+-------+           +-----------+
//         +----------+            |
//         |   r2     |            |
//         |respondent+------------+
//         |          |
//         +----------+
// #+end_src
//
func TestSurvey(t *testing.T) {
	const (
		addr = "tcp://127.0.0.1:59999"
	)

	type Message struct {
		Action     string        `json:"action"`
		Respondent string        `json:"respondent"`
		Reply      string        `json:"reply"`
		Duration   time.Duration `json:"duration"`
	}

	var (
		replies = make(chan *Message, 1024)
		locker  = &sync.Mutex{}
		cond    = sync.NewCond(locker)
	)
	closes := make(chan func() error, 1024)
	defer func() {
		max := len(closes)
		for i := 0; i < max; i++ {
			f := <-closes
			f()
		}
	}()
	listen := func() (sock mangos.Socket, err error) {
		sock, err = surveyor.NewSocket()
		if err != nil {
			tLog("listen", err)
			return
		}
		if err = sock.Listen(addr); err != nil {
			tLog("listen", err)
		}
		closes <- sock.Close
		sock.SetOption(mangos.OptionSurveyTime, time.Millisecond*10)
		cond.Broadcast()
		return
	}

	dial := func() (sock mangos.Socket, err error) {
		sock, err = respondent.NewSocket()
		if err != nil {
			tLog("dial", err)
			return
		}
		locker.Lock()
		cond.Wait()
		locker.Unlock()
		if err = sock.Dial(addr); err != nil {
			tLog("dial", err)
		}
		closes <- sock.Close
		return
	}

	ping := func(sock mangos.Socket) (err error) {
		m := &Message{
			Action: "ping",
		}

		b, err := json.Marshal(m)
		if err != nil {
			tLog("ping", err)
			return
		}
		// wait a little while
		time.Sleep(time.Millisecond)
		startTime := time.Now()
		if err = sock.Send(b); err != nil {
			tLog("ping", "send", err)
			return
		}

		for {
			b, err = sock.Recv()
			if err != nil {
				tLog("ping", "recv", err)
				return
			}
			m := &Message{}
			if err = json.Unmarshal(b, m); err != nil {
				tLog("ping", err)
				return
			}
			m.Duration = time.Now().Sub(startTime)
			replies <- m
		}
		return
	}

	tong := func(sock mangos.Socket, name string) (err error) {
		b, err := sock.Recv()
		if err != nil {
			tLog("tong", err)
			return
		}
		m := &Message{}
		if err = json.Unmarshal(b, m); err != nil {
			tLog("tong", err)
			return
		}
		m.Respondent = name
		if b, err = json.Marshal(m); err != nil {
			tLog("tong", err)
			return
		}
		if err = sock.Send(b); err != nil {
			tLog("tong", err)
			return
		}
		return
	}
	var g errgroup.Group
	// surveyer
	g.Go(func() (err error) {
		sock, err := listen()
		if err != nil {
			tLog(err)
			return
		}
		return ping(sock)
	})

	// 3 responsents
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("r%d", i)
		g.Go(func() (err error) {
			sock, err := dial()
			if err != nil {
				tLog(err)
				return
			}
			return tong(sock, name)
		})
	}
	err := g.Wait()
	t.Log(err)
	if len(replies) != 3 {
		t.Fatal(len(replies))
	}
	for i := 0; i < 3; i++ {
		m := <-replies
		tLog(m)
		if m.Duration > time.Millisecond*10 {
			t.Fatal(m)
		}
	}
}

// Star
// #+begin_src artist
//
//   +----------+               +----------+              +----------+
//   |    p1    |               |    p0    |              |    p3    |
//   |  (dial)  |<------------->|  (recv)  |<------------>|  (dial)  |
//   |          |               |          |              |          |
//   +----------+               +----------+              +----------+
//                                    ^
//                                    |
//                                    v
//                              +----------+
//                              |    p2    |
//                              |  (dial)  |
//                              |          |
//                              +----------+
//
//   Notes:
//   1. p1 send message, p3, p2 should receive the message,
//   2. p0 should receive the message
// #+end_src
//
func TestStar(t *testing.T) {
	const (
		addr = "tcp://127.0.0.1:59998"
	)

	// type Message struct {
	// 	Sender   string `json:"sender"`
	// 	Receiver string `json:"receiver"`
	// 	Content  string `json:"content"`
	// }

	listen := func() (sock mangos.Socket, err error) {
		if sock, err = star.NewSocket(); err == nil {
			err = sock.Listen(addr)
		}
		return
	}

	dial := func() (sock mangos.Socket, err error) {
		if sock, err = star.NewSocket(); err == nil {
			err = sock.Dial(addr)
		}
		return
	}

	p0, err := listen()
	if err != nil {
		t.Fatal(err)
	}
	defer p0.Close()
	p1, err := dial()
	defer p1.Close()
	p2, err := dial()
	defer p2.Close()
	p3, err := dial()
	defer p3.Close()
	msgs := make(chan string, 1024)
	recv := func(sock mangos.Socket, name string) (err error) {
		var b []byte
		for {
			b, err = sock.Recv()
			if err != nil {
				return
			}
			msg := string(b)
			msg = fmt.Sprintf("%s received at %s", msg, name)
			msgs <- msg
			if msg == "done" {
				return
			}
		}
		return
	}

	send := func(sock mangos.Socket, msg string) (err error) {
		if err = sock.Send([]byte(msg)); err != nil {
			t.Fatal(err)
		}
		return
	}

	var g errgroup.Group
	// p1 recv
	g.Go(func() error {
		return recv(p1, "p1")
	})

	// p2 recv
	g.Go(func() error {
		return recv(p2, "p2")
	})

	// p3 recv
	g.Go(func() error {
		return recv(p3, "p3")
	})

	// p0 recv
	g.Go(func() error {
		return recv(p0, "p0")
	})

	// p1 send hello
	if err := send(p1, "hello"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		m := <-msgs
		t.Log(m)
		if strings.Contains(m, "p1") {
			t.Fatal(m)
		}
	}
	// p0 send hello
	if err := send(p0, "hello"); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		m := <-msgs
		t.Log(m)
		if strings.Contains(m, "p0") {
			t.Fatal(m)
		}
	}
	// p0, p1 send done
	for _, sock := range []mangos.Socket{p0, p1} {
		if err := send(sock, "done"); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 4; i++ {
		m := <-msgs
		t.Log(m)
	}
	if l := len(msgs); l > 0 {
		t.Fatal(l)
	}
}

func TestMangosPushPull(t *testing.T) {
	const addr = "tcp://127.0.0.1:59999"
	sockPull, err := pull.NewSocket()
	if err != nil {
		t.Fatal(err)
	}

	if err = sockPull.Listen(addr); err != nil {
		t.Fatal(err)
	}
	defer sockPull.Close()
	buf := make(chan []byte, 1024)
	g := errgroup.Group{}
	g.Go(func() (err error) {
		b, err := sockPull.Recv()
		if err != nil {
			return
		}
		buf <- b
		return
	})
	sockPush, err := push.NewSocket()
	if err != nil {
		t.Fatal(err)
	}

	if err = sockPush.Dial(addr); err != nil {
		t.Fatal(err)
	}
	defer sockPush.Close()
	err = sockPush.Send([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	g.Wait()

	if l := len(buf); l != 1 {
		t.Fatal(l)
	}

	want := string(<-buf)
	if want != "hello" {
		t.Fatal(want)
	}
}
