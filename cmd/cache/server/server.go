package main

import (
	"log"
	"runtime"
	"strings"

	"github.com/dgraph-io/ristretto"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
	corev1 "github.com/tachunwu/distsys/pkg/proto/core/v1"
	"google.golang.org/protobuf/proto"
)

var shardId = pflag.String("shard-id", "", "Shard id for this instance")
var proxyPeers = pflag.StringSlice("proxy-peers", []string{"localhost:4222", "localhost:4223", "localhost:4224"}, "IP address of each peer in the proxy layer")
var shardSubject = strings.Join([]string{"CACHE.", *shardId}, "")

func main() {
	// Parse flag
	pflag.Parse()

	// Connect to NATS
	nc, _ := nats.Connect(strings.Join(*proxyPeers, ","))
	defer nc.Drain()

	// Cache init
	cache := NewCache()

	// Cache handers
	nc.Subscribe("CACHE.shard-0", CacheHandler(cache))
	for {
		runtime.Gosched()
	}
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

func CacheHandler(cache *ristretto.Cache) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {

		// Unmarshal
		txn := &corev1.Transaction{}
		err := proto.Unmarshal(msg.Data, txn)
		if err != nil {
			log.Println("unmarshal fail:", err)
			return
		}
		// Info
		// log.Println("TxnId:", txn.TxnId, "Type:", txn.TxnType, "Key:", txn.Key, "Value:", txn.Value)

		// Cache
		switch txn.GetTxnType() {
		case corev1.Transaction_GET:
			value, ok := cache.Get(txn.GetKey())
			if ok {
				msg.Respond(value.([]byte))
				return
			}
			msg.Respond([]byte("GET fail"))
		case corev1.Transaction_SET:
			ok := cache.Set(txn.GetKey(), txn.GetValue(), 1)
			if ok {
				msg.Respond([]byte("SET success"))
				return
			}
			msg.Respond([]byte("SET fail"))
		case corev1.Transaction_DELETE:
			cache.Del(txn.GetKey())
			msg.Respond([]byte("DELETE success"))
		default:
			log.Println("need txn type")
		}

	}
}
