package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	js, _ := nc.JetStream()

	for i := 0; i < 100000; i++ {
		log.Println(i)
		msg := &nats.Msg{
			Subject: "Materialize",
			Header:  nats.Header{},
		}
		msg.Header.Add("Id", fmt.Sprint(i))
		js.PublishMsgAsync(msg)
	}

}
