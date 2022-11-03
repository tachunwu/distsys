package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
	corev1 "github.com/tachunwu/distsys/pkg/proto/core/v1"
	"google.golang.org/protobuf/proto"
)

var proxyPeers = pflag.StringSlice("proxy-peers", []string{"localhost:4222", "localhost:4223", "localhost:4224"}, "IP address of each peer in the proxy layer")
var txnNum = 10000

var txnPool = sync.Pool{
	New: func() interface{} {
		return new(corev1.Transaction)
	},
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(*bytes.Buffer)
	},
}

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(strings.Join(*proxyPeers, ","))
	defer nc.Drain()

	// Publish some msg
	start := time.Now()
	// Set
	TestSet(nc)
	// Get
	TestGet(nc)
	// Delete
	TestDelete(nc)
	log.Println("Process", txnNum*3, "txn in", time.Since(start))
}

func TestGet(nc *nats.Conn) {
	for i := 0; i < txnNum; i++ {
		// Marshal
		txn := txnPool.Get().(*corev1.Transaction)
		txn.TxnId = uint64(i)
		txn.TxnType = corev1.Transaction_GET
		txn.Key = []byte("key" + strconv.Itoa(i))

		// log.Println("TxnId:", i, "Type:", txn.TxnType, "Key:", txn.Key)
		b, err := proto.Marshal(txn)
		if err != nil {
			log.Panicln(err)
			return
		}
		txnPool.Put(txn)

		// Request-Reply
		rply, err := nc.Request("CACHE.shard-0", b, time.Second)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(rply.Data))
	}
}

func TestSet(nc *nats.Conn) {
	for i := 0; i < txnNum; i++ {
		// Marshal
		txn := txnPool.Get().(*corev1.Transaction)
		txn.TxnId = uint64(i)
		txn.TxnType = corev1.Transaction_SET
		txn.Key = []byte("key" + strconv.Itoa(i))
		txn.Value = []byte("value" + strconv.Itoa(i))
		// log.Println("TxnId:", i, "Type:", txn.TxnType, "Key:", txn.Key, "Value:", txn.Value)
		b, err := proto.Marshal(txn)
		if err != nil {
			log.Panicln(err)
			return
		}
		txnPool.Put(txn)

		// Request-Reply
		rply, err := nc.Request("CACHE.shard-0", b, time.Second)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(rply.Data))
	}
}

func TestDelete(nc *nats.Conn) {
	for i := 0; i < txnNum; i++ {
		// Marshal
		txn := txnPool.Get().(*corev1.Transaction)
		txn.TxnId = uint64(i)
		txn.TxnType = corev1.Transaction_DELETE
		txn.Key = []byte("key" + strconv.Itoa(i))

		// log.Println("TxnId:", i, "Type:", txn.TxnType, "Key:", txn.Key)
		b, err := proto.Marshal(txn)
		if err != nil {
			log.Panicln(err)
			return
		}
		txnPool.Put(txn)

		// Request-Reply
		rply, err := nc.Request("CACHE.shard-0", b, time.Second)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(rply.Data))
	}
}
