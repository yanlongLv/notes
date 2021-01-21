package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	// err = ch.ExchangeDeclare(
	// 	"logs",   // name
	// 	"fanout", // type
	// 	true,     // durable
	// 	false,    // auto-deleted
	// 	false,    // internal
	// 	false,    // no-wait
	// 	nil,      // arguments)
	// )
	q, err := ch.QueueDeclare("rabbit_mq_test", true, false, false, false, nil)
	//q1, err = ch.QueueDeclare("rabbit_mq_test2", false, false, false, false, nil)
	// if err != nil {
	// 	panic(err)
	// }
	body := "Hello World"
	for i := 0; i < 10; i++ {
		body = fmt.Sprintf("%s,%d", body, i)
		print(i)
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "text/plain", Body: []byte(body)})
		if err != nil {
			fmt.Errorf("error %v", err)
		}

	}

	select {}
}
