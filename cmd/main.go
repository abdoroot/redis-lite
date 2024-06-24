package main

import (
	"github.com/abdoroot/lite-redis/internal/redis"
)

func main() {
	s := redis.NewServer(redis.Options{})
	s.Run()
}
