package redis

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type Payload struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	UUID string `json:"uuid"`
}
type Response struct {
	Url    string             `json:"url"`
	Type   string             `json:"type"`
	UUID   string             `json:"uuid"`
	Result *ffprobe.ProbeData `result:"uuid"`
}

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDRESS"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0,
})

func SaveResult(payload Payload, result *ffprobe.ProbeData) bool {
	var resp Response
	resp.Result = result
	resp.Type = payload.Type
	resp.UUID = payload.UUID
	resp.Url = payload.Url
	data, _ := json.Marshal(resp)

	err := rdb.Set(ctx, payload.UUID, data, 3*time.Hour).Err()
	return err == nil
}

func FetchResult(uuid string) Response {
	var resp Response

	data, _ := rdb.Get(ctx, uuid).Result()

	// if err != nil {
	// 	return ""
	// }

	json.Unmarshal([]byte(data), &resp)

	return resp
	// data, _ := json.Marshal(resp)
	// return err == nil
}
