package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
	"log"
)

func main() {
	const (
		messageBroker = "nats://nats:8222"
	)
	log.Println("Connecting to nats")
	nc, _ := nats.Connect(messageBroker)
	defer nc.Close()
	log.Println("Connected to nats")

	js, _ := nc.JetStream()

	log.Println("Waiting for messages...")
	tr, _ := js.Subscribe("USERS.*", func(msg *nats.Msg) {
		fmt.Printf("%s address on %s\n", msg.Subject, msg.Data)
	})
	defer tr.Unsubscribe()
	time.Sleep(10 * time.Minute)
}
