package metadata

import (
	"context"
	"log"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

var Dummy *ffprobe.ProbeData

func GetMetadata(url string) *ffprobe.ProbeData {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, url)
	if err != nil {
		log.Println(err)
		return Dummy
	}
	return data
}
