package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/aristorinjuang/go-rabbitmq/pkg"
	"github.com/streadway/amqp"
)

func fibonacciRPC(n int) (res int, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	pkg.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	pkg.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	pkg.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	pkg.FailOnError(err, "Failed to register a consumer")

	corrId := pkg.RandomString(32)

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	pkg.FailOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			pkg.FailOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	nString := pkg.BodyFrom(os.Args)
	n, err := strconv.Atoi(nString)
	pkg.FailOnError(err, "Failed to convert n from string to int")

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := fibonacciRPC(n)
	pkg.FailOnError(err, "Failed to handle RPC request")

	log.Printf(" [.] Got %d", res)
}
