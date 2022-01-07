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
	var body Body
	err := ctx.ReadJSON(&body)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	id, _ := gonanoid.New()

	var payload = &redis.Payload{
		Type: body.Type,
		Url:  body.Url,
		UUID: id,
	}

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
		go process(*payload)
		ctx.StatusCode(202)
		ctx.JSON(iris.Map{
			"message": "processing",
			"uuid":    payload.UUID,
		})
		return
	}
}

func process(payload redis.Payload) {
	result := metadata.GetMetadata(payload.Url)
	if result != "" {
		if redis.SaveResult(payload, result) {
			fmt.Println("Success")
		} else {
			fmt.Println("Failed to save!")
		}

	} else {
		fmt.Println("Failed to Convert!")
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
	result := redis.FetchResult(uuid)
	ctx.JSON(iris.Map{
		"message": "Success",
		"uuid":    uuid,
		"result":  result,
	})
}
