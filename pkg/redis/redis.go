package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type Payload struct {
	Url  string
	Type string
	UUID string
}
type Response struct {
	Url    string
	Type   string
	UUID   string
	Result string
}

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func SaveResult(payload Payload, result string) bool {
	var resp Response
	resp.Result = result
	data, _ := json.Marshal(resp)

	err := rdb.Set(ctx, payload.UUID, data, 0).Err()
	return err == nil
}
