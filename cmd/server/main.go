package main

import (
	"github.com/go-redis/redis/v8"
	
	"fmt"
	"context"
	"time"
	"encoding/json"
	"reflect"
)

// Payload struct
type Payload struct {
	ID		int
	Data    interface{} `json:"data"`
}

var ctx = context.Background()

func main() {
	fmt.Println("Go - Redis Project")

	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    
	// basic(rdb)
	testTypes(rdb)
	// lookAndFeel(rdb)
	// publisMsg(rdb)
}

func basic(rdb *redis.Client) {
	err := rdb.Set(ctx, "name", "bond", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := rdb.Get(ctx, "name").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := rdb.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
	}
}

func testTypes(rdb *redis.Client) {
	// TEST INTEGER
	err := rdb.Set(ctx, "int", 1, 0).Err()
    if err != nil {
        panic(err)
    }

    val1, _ := rdb.Get(ctx, "int").Result()
	fmt.Println("Get 1")
	fmt.Println("int", reflect.TypeOf(val1), val1)

	// TEST STRING
	err = rdb.Set(ctx, "string", "qwe", 0).Err()
    if err != nil {
        panic(err)
    }

    val2, _ := rdb.Get(ctx, "string").Result()
	fmt.Println("Get 2")
	fmt.Println("string", reflect.TypeOf(val2), val2 + " asd")

	// TEST ARRAY
	err = rdb.RPush(ctx, "list", 1, 2, 3, 4).Err()
    if err != nil {
        panic(err)
    }

	val3, _ := rdb.LRange(ctx, "list", 0, -1).Result()
	p, _ := rdb.RPop(ctx, "list").Result()
	fmt.Println("Get 3")
	fmt.Println("list", reflect.TypeOf(val3), val3)
	fmt.Println("list", reflect.TypeOf(p), p)

	// TEST HASH
	err = rdb.HMSet(ctx, "hash", "one", 1, "two", 2).Err()
    if err != nil {
        panic(err)
    }

	val41, _ := rdb.HMGet(ctx, "hash", "one").Result()
	val42, _ := rdb.HMGet(ctx, "hash", "two").Result()
	fmt.Println("Get 4")
	fmt.Println("hash", reflect.TypeOf(val41), val41[0])
	fmt.Println("hash", reflect.TypeOf(val42), val42[0])

}

func lookAndFeel(rdb *redis.Client) {
	// SET key value EX 10 NX
	set1, err := rdb.SetNX(ctx, "key", "value", 5*time.Second).Result()
	fmt.Println("Set 1: ")
	fmt.Println(set1, err)

	time.Sleep(4 * time.Second)
	get1, err:= rdb.Get(ctx, "key").Result()
	fmt.Println("Get 1: ")
	fmt.Println(get1, err)

	fmt.Println("")

	// SET key value keepttl NX
	// redis.KeepTTL will make the key cannot expire
	set2, err := rdb.SetNX(ctx, "key2", "value2", redis.KeepTTL).Result()
	fmt.Println("Set 2: ")
	fmt.Println(set2, err)

	get2, err:= rdb.Get(ctx, "key2").Result()
	fmt.Println("Get 2: ")
	fmt.Println(get2, err)

	fmt.Println("")

	// SORT list LIMIT 0 2 ASC
	vals1, err := rdb.Sort(ctx, "list", &redis.Sort{ Offset: 0, Count: 2, Order: "ASC" }).Result()
	fmt.Println("Set 3: ")
	fmt.Println(vals1, err)

	get3, err:= rdb.Get(ctx, "list").Result()
	fmt.Println("Get 3: ")
	fmt.Println(get3, err)

	fmt.Println("")

	// ZRANGEBYSCORE zset -inf +inf WITHSCORES LIMIT 0 2
	vals2, err := rdb.ZRangeByScoreWithScores(ctx, "zset", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
		Offset: 0,
		Count: 2,
	}).Result()
	fmt.Println("Set 4: ")
	fmt.Println(vals2, err)

	get4, err:= rdb.Get(ctx, "zset").Result()
	fmt.Println("Get 4: ")
	fmt.Println(get4, err)

	fmt.Println("")

	// ZINTERSTORE out 2 zset1 zset2 WEIGHTS 2 3 AGGREGATE SUM
	vals3, err := rdb.ZInterStore(ctx, "out", &redis.ZStore{
		Keys: []string{"zset1", "zset2"},
		Weights: []float64{2, 3},
	}).Result()
	fmt.Println("Set 5: ")
	fmt.Println(vals3, err)

	get5, err:= rdb.Get(ctx, "out").Result()
	fmt.Println("Get 5: ")
	fmt.Println(get5, err)

	fmt.Println("")

	// EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
	vals4, err := rdb.Eval(ctx, "return {KEYS[1],ARGV[1]}", []string{"key"}, "hello").Result()
	fmt.Println("Set 6: ")
	fmt.Println(vals4, err)

	fmt.Println("")

	// custom command
	res, err := rdb.Do(ctx, "set", "asd", "qwe").Result()
	fmt.Println("Set 7: ")
	fmt.Println(res, err)

	get6, err:= rdb.Get(ctx, "asd").Result()
	fmt.Println("Get 6: ")
	fmt.Println(get6, err)
}

func publisMsg(rdb *redis.Client) {
	payload1, _ := json.Marshal(Payload{
		ID: 1,
		Data: []int{1,2,3,4,5},
	})

	err := rdb.Publish(ctx, "mychannel1", payload1).Err()
	if err != nil {
		panic(err)
	}
}