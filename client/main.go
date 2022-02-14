package main

import (
	"encoding/json"
	"fmt"
	"golang_task_client/models"
	"net/http"
	"sync"

	"github.com/streadway/amqp"
)

type Action struct {
	data string
	cond *sync.Cond
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}

func GetAction(w http.ResponseWriter, req *http.Request) {
	action.cond.L.Lock()
	for action.data == "" {
		action.cond.Wait()
	}
	fmt.Fprintf(w, "%v", action.data)
	action.cond.L.Unlock()
}

var action Action

func main() {
	action = Action{
		data: "",
		cond: sync.NewCond(&sync.Mutex{}),
	}

	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"action", // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		"action", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			var entity models.Response

			_ = json.Unmarshal(d.Body, &entity)
			action.data = entity.Action
			action.cond.Broadcast()
			fmt.Println("Recieved", entity.Action)
		}
	}()

	http.HandleFunc("/action", GetAction)
	http.ListenAndServe(":80", nil)
}
