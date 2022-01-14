package main

import (
	"os"

	cors "github.com/iris-contrib/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
	"github.com/teezzan/ffmpeg-worker/pkg/controller"
)

var port = os.Getenv("PORT")

func main() {
	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	app.UseRouter(crs)

	app.Get("/fetch/{uuid}", controller.GetResult)

	app.Get("/total", controller.GetTotalSeconds)

	app.Post("/convert", controller.GetMetaFromURL)

	app.Listen(":" + string(port))
}
