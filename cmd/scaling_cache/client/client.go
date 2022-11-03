package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	Set(nc)
	Get(nc)
	Del(nc)

}

func Get(nc *nats.Conn) {
	// Get
	msg := &nats.Msg{
		Subject: "cache",
		Header:  nats.Header{},
	}
	msg.Header.Add("op", "GET")
	msg.Header.Add("key", "key01")

	r, err := nc.RequestMsg(msg, 1*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Get cache with key:", r.Header.Get("key"), "value:", r.Header.Get("value"))
}

func Set(nc *nats.Conn) {
	// Set
	msg := &nats.Msg{
		Subject: "cache",
		Header:  nats.Header{},
	}
	msg.Header.Add("op", "SET")
	msg.Header.Add("key", "key01")
	msg.Header.Add("value", "value01")

	r, err := nc.RequestMsg(msg, 1*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Set cache with key:", r.Header.Get("key"), "value:", r.Header.Get("value"))
}

func Del(nc *nats.Conn) {
	// Del
	msg := &nats.Msg{
		Subject: "cache",
		Header:  nats.Header{},
	}

	msg.Header.Add("op", "DEL")
	msg.Header.Add("key", "key01")

	r, err := nc.RequestMsg(msg, 1*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Del cache with key:", r.Header.Get("key"))
}
