package controller

import (
	"fmt"
	"os"

	"github.com/kataras/iris/v12"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/teezzan/ffmpeg-worker/pkg/metadata"
	"github.com/teezzan/ffmpeg-worker/pkg/redis"
)

type Body struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

var queueRequest = os.Getenv("QUEUE_REQUEST") == "true"

func GetMetaFromURL(ctx iris.Context) {
	var body Body
	err := ctx.ReadJSON(&body)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	if queueRequest {
		fmt.Println("Queued")

	} else {
		fmt.Println("Processing")

		var payload redis.Payload
		payload.Type = body.Type
		payload.Url = body.Url
		id, _ := gonanoid.New()
		payload.UUID = id
		go process(payload)
		fmt.Println(payload)
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
