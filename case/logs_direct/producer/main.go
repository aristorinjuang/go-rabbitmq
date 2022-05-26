package main

import (
	"log"
	"os"

	"github.com/aristorinjuang/go-rabbitmq/pkg"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	pkg.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	pkg.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	pkg.FailOnError(err, "Failed to declare an exchange")

	body := pkg.BodyFrom(os.Args)
	err = ch.Publish(
		"logs_direct",             // exchange
		pkg.SeverityFrom(os.Args), // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	pkg.FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
