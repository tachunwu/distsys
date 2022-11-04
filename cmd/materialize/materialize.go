package main

import (
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/tachunwu/distsys/pkg/jetstream"
	"github.com/tachunwu/distsys/pkg/materialize/consumers"
	"github.com/tachunwu/distsys/pkg/materialize/streams"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	js, _ := nc.JetStream()

	// Create stream
	jetstream.CreateStreams(streams.Streams)

	materializeConsumer := consumers.NewSnapshotConsumer(js, "MaterializeConsumer", "Materialize")
	materializeConsumer.Start()

	for {
		runtime.Gosched()
	}
}
