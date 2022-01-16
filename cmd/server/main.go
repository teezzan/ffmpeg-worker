package main

import (
	"fmt"
	"log"
	"os"

	"github.com/iris-contrib/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
	"github.com/teezzan/ffmpeg-worker/pkg/controller"
	"github.com/teezzan/ffmpeg-worker/pkg/redis"
	socketio "github.com/teezzan/ffmpeg-worker/pkg/socketio"
)

var port = os.Getenv("PORT")

type Body struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

func main() {
	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	app.UseRouter(crs)

	app.HandleMany("GET POST", "/socket.io/{any:path}", iris.FromStd(socketio.Server))

	app.Get("/fetch/{uuid}", controller.GetResult)

	app.Get("/total", controller.GetTotalSeconds)

	app.Post("/convert", LookupCache, controller.GetMetaFromURL)

	defer socketio.Server.Close()

	if err := app.Run(
		iris.Addr(":"+port),
		iris.WithoutPathCorrection,
		iris.WithoutServerError(iris.ErrServerClosed),
	); err != nil {
		log.Fatal("failed run app: ", err)
	}
}

func LookupCache(ctx iris.Context) {
	var body Body
	err := ctx.ReadJSON(&body)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	ctx.Values().Set("url", body.Url)
	ctx.Values().Set("type", body.Type)
	result, found := redis.FetchResultFromCache(body.Url)
	if found {
		if result.Error != "" {
			ctx.JSON(iris.Map{
				"message": "Failed",
				"uuid":    result.UUID,
				"error":   result.Error,
				"data":    result,
			})
			return
		} else {
			ctx.JSON(iris.Map{
				"message": "Success",
				"error":   nil,
				"data":    result,
			})
			return
		}

	}
	fmt.Println("cACHED mISSED")
	ctx.Next()

}
