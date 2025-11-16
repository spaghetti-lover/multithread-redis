package main

import (
	"fmt"
	"log"

	"github.com/spaghetti-lover/multithread-redis/client"
)

func main() {
	// Kết nối đến Redis server
	c, err := client.NewClient("localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Ping
	pong, err := c.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PING:", pong)

	// SET
	err = c.Set("mykey", "Hello Redis!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SET mykey successfully")

	// GET
	value, err := c.Get("mykey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET mykey:", value)

	// SET with expiration
	err = c.SetEx("tempkey", "temporary value", 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SET tempkey with 10s expiration")
}
