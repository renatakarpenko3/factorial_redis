// factorial_with_redis.go
package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// Initialize Redis client
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // Address of the Redis server
})

// Recursive factorial function
func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

// Redis-integrated factorial function
func getFactorialWithRedis(n int) int {
	cacheKey := fmt.Sprintf("factorial:%d", n)
	cachedValue, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil { // Cache hit
		value, _ := strconv.Atoi(cachedValue)
		fmt.Println("Retrieved from Redis cache!")
		return value
	}

	// If not cached, calculate and store the result
	result := factorial(n)
	err = rdb.Set(ctx, cacheKey, result, 0).Err() // Store in Redis
	if err != nil {
		log.Fatalf("Failed to cache result: %v", err)
	}
	return result
}

func main() {
	var n int
	fmt.Print("Enter a number: ")
	fmt.Scan(&n)

	// Measure time for Redis-integrated factorial calculation
	start := time.Now()
	resultWithRedis := getFactorialWithRedis(n)
	duration := time.Since(start)

	fmt.Printf("Factorial of %d is %d, took %v\n", n, resultWithRedis, duration)
}
