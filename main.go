package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

var client_ *redis.Client

func main() {

	if err := initialize(); err != nil {
		// initialize error: dial tcp [::1]:6379: connect: connection refused
		fmt.Println("initialize error:", err)
		return
	}

	// set key value のサンプル
	sample_set()
	// リスト rpush key value のサンプル
	sample_rpush()
	// セット sadd key member のサンプル
	sample_sadd()
	// ソート済みセット zadd のサンプル
	sample_zadd()
	// ハッシュ型 hmset のサンプル
	sample_hmset()
}

func initialize() error {

	client_ = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client_.Ping().Result()
	return err
}

func sample_set() {

	fmt.Println("sample_set start ----------------")

	key := "Key001"

	err := client_.Set(key, "value001", 0).Err()
	if err != nil {
		fmt.Print(err)
		return
	}

	val, err := client_.Get(key + "---").Result()
	if err == redis.Nil {
		fmt.Println("key001 does not exist")
	}

	val, err = client_.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("key001 does not exist")
		return
	} else if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(val)
}

func sample_rpush() {
	fmt.Println("sample rpush start -------------")
	key := "rpush_key"
	err := client_.RPush(key, "A").Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	client_.RPush(key, "B")
	client_.RPush(key, "C")
	client_.RPush(key, "D")

	val, err := client_.LRange(key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(val)

}

func sample_sadd() {
	fmt.Println("sample sadd ----------------")
	key := "sadd_key"
	err := client_.SAdd(key, "1", "2", "3", "4", "5").Err()
	if err != nil {
		fmt.Println(err)
		return
	}

	vals, err := client_.SMembers(key).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	// [1 2 3 4 5]
	fmt.Println(vals)

	smap, err := client_.SMembersMap(key).Result()
	if _, ok := smap["2"]; ok {
		fmt.Println("contains 2")
	}
	if _, ok := smap["10"]; !ok {
		fmt.Println("not contains 10")
	}
}

func sample_zadd() {
	fmt.Println("sample zadd --------------")
	key := "zadd_key"

	err := client_.ZAdd(
		key,
		&redis.Z{Score: 10, Member: "A"},
		&redis.Z{Score: 1, Member: "B"},
	).Err()

	if err != nil {
		fmt.Println(err)
		return
	}

	val, err := client_.ZRange(key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	// [B A]
	fmt.Println(val)
}

func sample_hmset() {
	fmt.Println("sample hmset -----------")
	key := "hmset_key"

	m := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	client_.HMSet(key, m)

	val, err := client_.HMGet(key, "a").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	// [1]
	fmt.Println(val)

	val, err = client_.HMGet(key, "a", "b", "c", "d").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	// [1 2 3 <nil>]
	fmt.Println(val)
}
