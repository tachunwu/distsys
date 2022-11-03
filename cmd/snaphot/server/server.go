package main

import (
	"log"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
	corev1 "github.com/tachunwu/distsys/pkg/proto/core/v1"
	"google.golang.org/protobuf/proto"
)

var streamName = "CACHE"
var shardId = pflag.String("shard-id", "", "Shard id for this instance")
var proxyPeers = pflag.StringSlice("proxy-peers", []string{"localhost:4222", "localhost:4223", "localhost:4224"}, "IP address of each peer in the proxy layer")
var shardSubject = strings.Join([]string{"CACHE.", *shardId}, "")

func main() {
	// Parse flag
	pflag.Parse()

	// Connect to NATS
	nc, _ := nats.Connect(strings.Join(*proxyPeers, ","))
	defer nc.Drain()

	// Create JetStream Context
	js, _ := nc.JetStream()

	// Create Stream(Server-side)
	js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"CACHE.>"},
	})

	// Create a Consumer (In-memory push-based client-side)
	js.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:        *shardId,
		MemoryStorage:  true,
		DeliverSubject: shardSubject,
		AckPolicy:      nats.AckExplicitPolicy,
		AckWait:        time.Second,
	})

	// Log handers
	nc.Subscribe(shardSubject, LogHandler())

}

func NewCache() *ristretto.Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return cache
}

func LogHandler() func(msg *nats.Msg) {
	return func(msg *nats.Msg) {}
}

func StartNewTransaction(txn *corev1.Transaction) {
	log.Println("txn: ", txn.TxnId)
}

func SnapshotHandler(js nats.JetStreamContext) {

	// Start subscribe shard stream for this topic
	sub, err := js.SubscribeSync(
		shardSubject,
		nats.Durable(*shardId),
		nats.AckExplicit(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Snapshot handers
	for {
		msg, err := sub.NextMsg(time.Second)
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Tansaction
		txn := &corev1.Transaction{}
		err = proto.Unmarshal(msg.Data, txn)
		if err != nil {
			log.Println(err)
			continue
		}

		// Process Transaction
		StartNewTransaction(txn)
		msg.Respond([]byte("transaction process success"))
	}
}
