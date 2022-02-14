package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}

func main() {
	forever := make(chan bool)

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

	go Worker("https://novasite.su/test1.php", ch)
	go Worker("https://novasite.su/test2.php", ch)

	<-forever
}

func Worker(url string, ch *amqp.Channel) {
	for {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(2000)
		time.Sleep(time.Duration(r)*time.Millisecond + time.Second)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()

		err = ch.Publish(
			"action", // exchange
			"",       // routing key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        data,
			})
		failOnError(err, "Failed to publish a message")
	}
}
