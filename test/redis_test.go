package test

import (
	"OJ/models"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// 初始化连接
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // 没有密码，默认值
	DB:       0,  // 默认DB 0
})

func TestRedisSet(t *testing.T) {
	rdb.Set(ctx, "name", "mmc", time.Second*100)
}

func TestRedisGet(t *testing.T) {
	v, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}

func TestRedisGetByModels(t *testing.T) {
	v, err := models.RDB.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}
