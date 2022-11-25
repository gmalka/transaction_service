package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/nats.go"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Transaction struct {
	ReceiverId int   `json:"Receiver"`
	SenderId   int   `json:"Sender"`
	Count      int64 `json:"Count"`
}

type ResultTransaction struct {
	UserId  int   `json:"UserId"`
	Remains int64 `json:"Remains"`
}

var (
	js nats.JetStreamContext
	db dbConn
)

func main() {
	start()
}

func start() {
	const (
		transactionUrl = "/users/:id/transactions"
		messageBroker  = "nats://nats:8222"
		jetStreamName  = "USERS"
		serverURL      = "0.0.0.0:8080"
	)

	log.Println("Connecting to db")
	err := db.createConnection("postgres", "postgres", "db", "5432", "postgres")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connect to db")
	defer db.Close()

	log.Println("Connecting to nats")
	nc, err := nats.Connect(messageBroker)
	log.Println("Connect to nats")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	js, _ = nc.JetStream()
	log.Println("Get jetStream")
	_, err = js.StreamInfo(jetStreamName)
	if err != nil {
		createStream()
	}

	router := httprouter.New()
	log.Println("Created router")
	router.POST(transactionUrl, createTransaction)

	server := &http.Server{
		Addr:         serverURL,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Starting server")
	log.Fatal(server.ListenAndServe())
}

func createStream() {
	var streamConfig nats.StreamConfig

	log.Println("Cant found stream USERS")
	cfg, err := os.ReadFile("streamConfig.json")
	if err != nil {
		log.Fatal("Cant find configure file")
	}
	err = json.Unmarshal(cfg, &streamConfig)
	if err != nil {
		log.Fatal("Cant Unmarshal log file")
	}
	_, err = js.AddStream(&streamConfig)
	if err != nil {
		log.Fatal("Cant Unmarshal log file")
	}
	log.Println("Stream USERS created")
}

func createTransaction(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var tr Transaction

	writer.WriteHeader(204)
	log.Println("Get request")
	if request.Body != nil {
		defer request.Body.Close()
		result, _ := io.ReadAll(request.Body)
		err := json.Unmarshal(result, &tr)
		if err != nil {
			log.Println("Unmarshal error")
			return
		}
		if tr.SenderId == tr.ReceiverId {
			log.Println("Ids are equals")
			return
		}
		resultTr, err := db.makeTransaction(tr)
		if err != nil {
			log.Println(err)
			return
		}
		r, _ := json.Marshal(resultTr)
		log.Println("Publishing...")
		log.Printf("%s\n", r)
		_, err = js.Publish(fmt.Sprintf("USERS.%d", resultTr.UserId), r)
		if err != nil {
			log.Println("Publish error")
			return
		}
		log.Println("Publish success")
	}
	log.Println("End request")
}
