package main

import (
	"log"
	"math"
	"strconv"

	"github.com/streadway/amqp"
)

func main() {
	strconv.Atoi()
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
	// 	nil,      // arguments
	// )

	q, err := ch.QueueDeclare("rabbit_mq_test", true, false, false, false, nil)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	forever := make(chan bool)
	// i := 0
	go func() {
		for d := range msgs {
			// i++
			// if i == 5 {
			// 	panic("err")
			// }
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func levelOrder(root *TreeNode) [][]int {
	quene := []*TreeNode{root}
	result := [][]int{}
	for len(quene) > 0 {
		ret := []int{}
		end := len(quene)
		for i := 0; i < end; i++ {
			ret = append(ret, quene[i].Val)
			if quene[i].Left != nil {
				quene = append(quene, quene[i].Left)
			}
			if quene[i].Right != nil {
				quene = append(quene, quene[i].Right)
			}
		}
		result = append(result, ret)
	}
	return result
}

func strToInt(str string) int {
	if len(str) == 0 {
		return 0
	}
	bytes := []byte(str)
	i := 0
	ch := bytes[i]
	for i < len(bytes) && ch == ' ' {
		i++
		ch = bytes[i]
	}
	bytes = bytes[i:]
	if len(bytes) == 0 {
		return 0
	}
	na := false
	if bytes[0] == '-' {
		na = true
		bytes = bytes[1:]
	} else if bts[0] == '+' {
		bts = bts[1:]
	}
	n := 0
	for _, ch := range bytes {
		ch -= '0'
		if ch > 9 {
			break
		}
		if n > (1<<31)-1 {
			if na {
				return math.MinInt32
			} else {
				return (1 << 31) - 1
			}
		}
		n = n*10 + int(ch)
	}
	if n > (1<<31)-1 {
		if na {
			return math.MinInt32
		} else {
			return (1 << 31) - 1
		}
	}
	if na {
		n = -n
	}
	return n

}

func levelOrder1(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	quene := []*TreeNode{root}
	result := []int{}
	for len(quene) > 0 {
		length := len(quene)
		for i := 0; i < length; i++ {
			result = append(result, quene[i].Val)
			if quene[i].Left != nil {
				quene = append(quene, quene[i].Left)
			}
			if quene[i].Right != nil {
				quene = append(quene, quene[i].Right)
			}
		}
	}
	return result
}
