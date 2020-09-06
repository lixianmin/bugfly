package epoller

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestPoller(t *testing.T) {
	num := 10
	msgPerConn := 10

	poller, err := NewPoller()
	if err != nil {
		t.Fatal(err)
	}

	// start server
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}

			poller.Add(conn)
		}
	}()

	// create num connections and send msgPerConn messages per connection
	for i := 0; i < num; i++ {
		go func() {
			conn, err := net.Dial("tcp", ln.Addr().String())
			if err != nil {
				t.Error(err)
				return
			}
			time.Sleep(200 * time.Millisecond)
			for i := 0; i < msgPerConn; i++ {
				conn.Write([]byte("hello world"))
			}
			conn.Close()
		}()
	}

	time.Sleep(100 * time.Millisecond)

	// read those num * msgPerConn messages, and each message (hello world) contains 11 bytes.
	ch := make(chan struct{})
	var total int
	var count int
	var expected = num * msgPerConn * len("hello world")
	go func() {
		for {
			conns, err := poller.Wait(128)
			if err != nil {
				t.Fatal(err)
			}
			count++
			var buf = make([]byte, 11)
			for _, conn := range conns {
				n, err := conn.Read(buf)
				if err != nil {
					if err == io.EOF {
						conn.Close()
						poller.Remove(conn)
					} else {
						t.Error(err)
					}
				}
				total += n
			}

			if total == expected {
				break
			}
		}

		t.Logf("read all %d bytes, count: %d", total, count)
		close(ch)
	}()

	select {
	case <-ch:
	case <-time.After(2 * time.Second):
	}

	if total != expected {
		t.Fatalf("epoller does not work. expect %d bytes but got %d bytes", expected, total)
	}
}
