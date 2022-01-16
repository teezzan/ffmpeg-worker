package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kataras/iris/v12"
	gonanoid "github.com/matoous/go-nanoid/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/teezzan/ffmpeg-worker/pkg/metadata"
	"github.com/teezzan/ffmpeg-worker/pkg/redis"
	"github.com/teezzan/ffmpeg-worker/pkg/socketio"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

type Body struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

var queueRequest = os.Getenv("QUEUE_REQUEST") == "true"
var amqpUrl = os.Getenv("AMQPURL")

var glob_publisher *rabbitmq.Publisher

func init() {

	publisher, init_err := rabbitmq.NewPublisher(
		amqpUrl, amqp.Config{},
		rabbitmq.WithPublisherOptionsLogging,
	)
	if init_err != nil {
		log.Fatal(init_err)
	}
	glob_publisher = publisher
}

func GetMetaFromURL(ctx iris.Context) {
	id, _ := gonanoid.New()
	var payload = &redis.Payload{
		Type: ctx.Values().GetString("type"),
		Url:  ctx.Values().GetString("url"),
		UUID: id,
	}
	socketID := ctx.Values().GetString("socket_id")

	if queueRequest {
		fmt.Println("Queued")
		result := enqueue(*payload)
		if result {
			ctx.StatusCode(202)
			ctx.JSON(iris.Map{
				"message": "processing",
				"uuid":    payload.UUID,
			})
			return
		} else {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"message": "Error",
			})
			return
		}
	} else {
		ctx.StatusCode(202)
		ctx.JSON(iris.Map{
			"message": "processing",
			"uuid":    payload.UUID,
		})
		go process(*payload, socketID)
		return
	}
}

func process(payload redis.Payload, socketID string) {
	result, error_message := metadata.GetMetadata(payload.Url)

	if redis.SaveResult(payload, result, error_message) {
		fmt.Println("Result Saved Successfully")

	} else {
		fmt.Println("Failed to save!")
	}
	if socketID != "" {
		socketio.EmitCompleted(payload.UUID, socketID)
	}

}
func enqueue(payload redis.Payload) bool {
	data, _ := json.Marshal(payload)

	err := glob_publisher.Publish(
		[]byte(data),
		[]string{"key"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("direct_xch"),
	)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GetResult(ctx iris.Context) {
	uuid := ctx.Params().Get("uuid")
	result, found := redis.FetchResult(uuid)
	if found {
		if result.Error != "" {
			ctx.JSON(iris.Map{
				"message": "Failed",
				"uuid":    uuid,
				"error":   result.Error,
				"data":    result,
			})
		} else {
			ctx.JSON(iris.Map{
				"message": "Success",
				"error":   nil,
				"data":    result,
			})
		}

	} else {
		ctx.StatusCode(404)
		ctx.JSON(iris.Map{
			"message": "Not Found",
			"uuid":    uuid,
			"error":   "uuid not found",
			"data":    nil,
		})
	}
}
func GetTotalSeconds(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"message":        "Success",
		"total_duration": redis.FetchTotalDuration(),
	})

}
