package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
	"github.com/teezzan/ffmpeg-worker/pkg/controller"
)

var port = os.Getenv("PORT")

func main() {
	app := iris.New()

	app.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "pong",
		})
	})

	app.Post("/convert", controller.GetMetaFromURL)

	app.Listen(":" + string(port))
}
