package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	ctx := context.Background()
	client(ctx, nc)
}

func client(ctx context.Context, nc *nats.Conn) {

	nc.PublishRequest("vectored-req", "vectored-rep", []byte("vectored"))
	sub, _ := nc.SubscribeSync("vectored-rep")

	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)
	res := []string{}
	msgCh := make(chan *nats.Msg, 1)
	go Fetch(ctx, sub, msgCh)

	for {
		select {
		case <-ctx.Done():
			fmt.Println(res)
			return
		case msg := <-msgCh:
			res = append(res, string(msg.Data))
		}
	}

}

func Fetch(ctx context.Context, sub *nats.Subscription, ch chan<- *nats.Msg) {
	for {
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			return
		}
		ch <- msg
	}
}
