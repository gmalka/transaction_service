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
	nc, err := nats.Connect(messageBroker)
	for err != nil {
		nc, err = nats.Connect(messageBroker)
		log.Println("connecting to nats...")
		time.Sleep(time.Second)
	}
	defer nc.Close()
	log.Println("Connected to nats")

	js, err := nc.JetStream()
	for err != nil {
		js, err = nc.JetStream()
		log.Println("connecting to jetstream...")
		time.Sleep(time.Second)
	}

	log.Println("Waiting for messages...")
	tr, err := js.Subscribe("USERS.*", func(msg *nats.Msg) {
		fmt.Printf("%s address on %s\n", msg.Subject, msg.Data)
	})
	for err != nil {
		tr, err = js.Subscribe("USERS.*", func(msg *nats.Msg) {
		fmt.Printf("%s address on %s\n", msg.Subject, msg.Data)
	})
		log.Println("connecting to stream...")
		time.Sleep(time.Second)
	}
	defer tr.Unsubscribe()
	time.Sleep(10 * time.Minute)
}
