package metadata

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func GetMetadata(url string) string {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, url)
	if err != nil {
		return ""
	}

	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Panicf("Error unmarshalling: %v", err)
	}
	return string(buf)
}
