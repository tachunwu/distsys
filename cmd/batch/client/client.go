package main

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	corev1 "github.com/tachunwu/distsys/pkg/proto/core/v1"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	js, _ := nc.JetStream()

	start := time.Now()
	UnbatchPublish(js)
	log.Println("Unbatch txn publish:", time.Since(start), "TPS:", 100000.0/time.Since(start).Seconds())

	start = time.Now()
	BatchPublish(js)
	log.Println("Batch txn publish:", time.Since(start), "TPS:", 100000.0/time.Since(start).Seconds())

}

func BatchPublish(js nats.JetStreamContext) {
	msgCh := make(chan *nats.Msg, 1)
	go BatchPublishProducer(msgCh)
	BatchPublishConsumer(msgCh, js)
}

func BatchPublishProducer(msgCh chan *nats.Msg) {
	for b := 0; b < 100; b++ {
		batch := &corev1.TransactionBatch{}
		for i := 0; i < 1000; i++ {

			t := &corev1.Transaction{
				TxnId: uint64(b*1000 + i),
				Key:   []byte(string("BYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTE")),
				Value: []byte(string("BYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTE")),
			}
			batch.Txns = append(batch.Txns, t)
		}

		batchBuf, _ := proto.Marshal(batch)
		msg := &nats.Msg{
			Subject: "Sequencer",
			Header:  nats.Header{},
			Data:    batchBuf,
		}
		// msg.Header.Add("Id", fmt.Sprint(b))
		msgCh <- msg
	}
}

func BatchPublishConsumer(msgCh chan *nats.Msg, js nats.JetStreamContext) {
	batchCount := 0
	for {
		b := <-msgCh
		js.PublishMsg(b)

		batchBuf := &corev1.TransactionBatch{}
		proto.Unmarshal(b.Data, batchBuf)

		// for _, txn := range batchBuf.GetTxns() {
		// 	log.Println("batch txn:", txn.TxnId)
		// }

		batchCount++
		// log.Println(batchCount)
		if batchCount == 100 {
			break
		}
	}
}

func UnbatchPublish(js nats.JetStreamContext) {

	msgCh := make(chan *nats.Msg, 1)

	// Producer
	go UnbatchPublishProducer(msgCh)

	// Consumer
	UnbatchPublishConsumer(msgCh, js)

}

func UnbatchPublishProducer(msgCh chan *nats.Msg) {
	for i := 0; i < 100000; i++ {
		// log.Println("unbatch txn:", i)
		t := &corev1.Transaction{
			TxnId: uint64(i),
			Key:   []byte(string("BYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTE")),
			Value: []byte(string("BYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTEBYTE")),
		}
		tBuf, _ := proto.Marshal(t)
		msg := &nats.Msg{
			Subject: "Sequencer",
			Header:  nats.Header{},
			Data:    tBuf,
		}
		// msg.Header.Add("Id", fmt.Sprint(i))
		msgCh <- msg
	}
}

func UnbatchPublishConsumer(msgCh chan *nats.Msg, js nats.JetStreamContext) {
	txnCount := 0
	for {
		js.PublishMsg(<-msgCh)
		txnCount++
		if txnCount == 99999 {
			break
		}
	}
}
