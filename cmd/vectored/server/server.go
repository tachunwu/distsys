package main

import (
	"runtime"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	server(nc, "shard-0")
	server(nc, "shard-1")
	server(nc, "shard-2")

	for {
		runtime.Gosched()
	}
}

func server(nc *nats.Conn, shard string) {
	nc.Subscribe("vectored-req", func(m *nats.Msg) {
		nc.Publish("vectored-rep", []byte(shard))
	})
}
