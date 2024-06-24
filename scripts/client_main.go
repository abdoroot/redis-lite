package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abdoroot/lite-redis/internal/redis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please type your command")
	}
	cmd := os.Args[1]
	options := redis.Options{Addr: "127.0.0.1:8080"}
	client := redis.NewClient(options)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	switch strings.ToLower(cmd) {
	case "set":
		k, v := os.Args[2], os.Args[3]
		fmt.Println(client.Set(ctx, k, v))
	case "get":
		k := os.Args[2]
		fmt.Println(client.Get(ctx, k))
	}
}
