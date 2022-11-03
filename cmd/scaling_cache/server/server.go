package main

import (
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
)

func main() {

	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	js, _ := nc.JetStream()

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "cache",
	})

	if err != nil {
		log.Println(err)
	}

	nc.Subscribe("cache", func(m *nats.Msg) {
		switch m.Header.Get("op") {
		case "GET":
			log.Println("GET op:", m.Header.Get("key"))
			value, err := kv.Get(m.Header.Get("key"))
			if err != nil {
				log.Println("miss")
				return
			}
			m.Header.Set("value", string(value.Value()))
			log.Println(string(value.Value()))
			nc.Publish(m.Reply, nil)
		case "SET":
			log.Println("SET op:", m.Header.Get("key"), m.Header.Get("value"))
			kv.PutString(m.Header.Get("key"), m.Header.Get("value"))
			m.Header.Set("key", "key01-done")
			m.Header.Set("value", "key01-done")
			log.Println(m.Header.Get("key"), m.Header.Get("value"))
			nc.Publish(m.Reply, nil)
		case "DEL":
			log.Println("DEL op:", m.Header.Get("key"))
			kv.Purge(m.Header.Get("key"))
			m.Header.Set("key", "key01-del")
			m.Header.Set("value", "key01-del")
			log.Println(m.Header.Get("key"), m.Header.Get("value"))
			nc.Publish(m.Reply, nil)
		}
	})

	for {
		runtime.Gosched()
	}
}
