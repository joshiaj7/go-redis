package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	
	"fmt"
	"context"
)

var ctx = context.Background()

// Payload struct
type Payload struct {
	ID		int
	Data    interface{} `json:"data"`
}

func main() {
	fmt.Println("Go - Redis Subsciber")

	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
	})
	
	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, "mychannel1")

	// first method to receive message
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
	
		fmt.Println(msg.Channel, msg.Payload)

		res := Payload{}

		err = json.Unmarshal([]byte(msg.Payload), &res)
		if err != nil {
			panic(err)
		}

		fmt.Println(res)
	}

	// second method to receive message
	// ch := pubsub.Channel()

	// for msg := range ch {
	// 	fmt.Println(msg.Channel, msg.Payload)
	// }
}