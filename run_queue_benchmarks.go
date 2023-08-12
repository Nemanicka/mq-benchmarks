package main

import (
	"context"
	"fmt"
	"time"

	"sync"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client
var aof *redis.Client
var np *redis.Client
var bs *beanstalk.Conn

type Object struct {
	Value  string
	Expiry time.Time
	Delta  time.Duration
}

func init() {
	bs, _ = beanstalk.Dial("tcp", "localhost:11300")

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6400",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	aof = redis.NewClient(&redis.Options{
		Addr:     "localhost:6401",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	np = redis.NewClient(&redis.Options{
		Addr:     "localhost:6402",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func writeTest(c *redis.Client, volume int, agent int) {
	for i := 0; i < volume; i++ {
		c.LPush(ctx, fmt.Sprintf("agent%dKey%d", agent, i), fmt.Sprintf(tShevchenko, i))
	}
}

func readTest(c *redis.Client, volume int, agent int) {
	for i := 0; i < volume; i++ {
		c.RPop(ctx, fmt.Sprintf("agent%dKey%d", agent, i))
	}
}

func benchmarkRedis(test func(*redis.Client, int, int), c *redis.Client, concurrency, volume int) {
	start := time.Now()
	var wg sync.WaitGroup
	for agent := 0; agent < concurrency; agent++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			test(c, volume, agent)
		}()
	}
	wg.Wait()
	end := time.Now()
	fmt.Printf("%v\n", end.Sub(start))
}

func benchmarkBS(concurrency, volume int, mode string) {
	start := time.Now()
	var wg sync.WaitGroup
	for agent := 0; agent < concurrency; agent++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < volume; i++ {
				if mode == "write" {
					data := fmt.Sprintf(tShevchenko, i)
					bs.Put([]byte(data), 1, 0, 5*time.Second)
				} else if mode == "read" {
					id, _, err := bs.Reserve(1 * time.Second)
					// fmt.Println(id)
					if err != nil {
						panic(err)
					}

					bs.Delete(id)
				} else {
					panic("Unknown mode: only write or read allowed")
				}
			}
		}()
	}
	wg.Wait()
	end := time.Now()
	fmt.Printf("%v\n", end.Sub(start))
}

func main() {

	fmt.Println("Check rdb...")
	concurrency := 100
	volume := 1000

	fmt.Printf("RDB write: ")
	benchmarkRedis(writeTest, rdb, concurrency, volume)

	fmt.Printf("RDB read: ")
	benchmarkRedis(readTest, rdb, concurrency, volume)

	fmt.Printf("AOF write: ")
	benchmarkRedis(writeTest, aof, concurrency, volume)

	fmt.Printf("AOF read: ")
	benchmarkRedis(readTest, aof, concurrency, volume)

	fmt.Printf("NP write: ")
	benchmarkRedis(writeTest, np, concurrency, volume)

	fmt.Printf("NP read: ")
	benchmarkRedis(readTest, np, concurrency, volume)

	fmt.Printf("BS write: ")
	benchmarkBS(concurrency, volume, "write")

	fmt.Printf("BS read: ")
	benchmarkBS(concurrency, volume, "read")
}
