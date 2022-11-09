package main

import (
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/tachunwu/distsys/pkg/jetstream"
	"github.com/tachunwu/distsys/pkg/seq/streams"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	// Create stream
	jetstream.CreateStreams(streams.Streams)

	for {
		runtime.Gosched()
	}
}
