package redis

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
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
	Error  string             `json:"error"`
	Result *ffprobe.ProbeData `json:"result"`
}

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDRESS"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0,
})
var totalKeyName = "total"

func SaveResult(payload Payload, result *ffprobe.ProbeData, error_message string) bool {
	var resp Response
	if error_message != "" {
		resp.Error = error_message
	} else {
		resp.Result = result
		SetTotal(int(result.Format.DurationSeconds))

	}
	resp.Type = payload.Type
	resp.UUID = payload.UUID
	resp.Url = payload.Url
	data, _ := json.Marshal(resp)

	err := rdb.Set(ctx, payload.UUID, data, 5*time.Hour).Err()
	if error_message == "" {
		rdb.Set(ctx, payload.Url, payload.UUID, 5*time.Hour).Err()
	}
	return err == nil

}

func FetchResult(uuid string) (Response, bool) {
	var resp Response

	data, err := rdb.Get(ctx, uuid).Result()

	if err != nil {
		return resp, false
	}

	json.Unmarshal([]byte(data), &resp)

	return resp, true
}

func FetchTotalDuration() string {
	data, err := rdb.Get(ctx, totalKeyName).Result()

	if err != nil {
		return "25000"
	}

	return data
}

func SetTotal(duration int) bool {
	data, err := rdb.Get(ctx, totalKeyName).Result()
	prev_duration := 0
	if err == nil {
		prev_duration, _ = strconv.Atoi(data)
	}

	new_duration := prev_duration + duration

	err = rdb.Set(ctx, totalKeyName, strconv.Itoa(new_duration), redis.KeepTTL).Err()
	return err == nil
}
func FetchResultFromCache(url string) (Response, bool) {
	var resp Response

	data, err := rdb.Get(ctx, url).Result()

	if err != nil {
		return resp, false
	}
	data, err = rdb.Get(ctx, data).Result()

	if err != nil {
		return resp, false
	}
	json.Unmarshal([]byte(data), &resp)

	return resp, true
}
