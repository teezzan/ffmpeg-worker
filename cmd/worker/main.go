package main

import (
	// "fmt"
	"context"
	"encoding/json"
	"log"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"

	// "os"
	// "os/signal"
	// "syscall"
	// amqp "github.com/rabbitmq/amqp091-go"
	// rabbitmq "github.com/wagslane/go-rabbitmq"
	"github.com/teezzan/ffmpeg-worker/pkg/metadata"
)

// var consumerName = "example"
// var amqpUrl = "amqps://sztwqfjl:la9uwaS03--T93hv0JuJsoiUgxxexMhw@rattlesnake.rmq.cloudamqp.com/sztwqfjl"

// func main() {
// 	consumer, err := rabbitmq.NewConsumer(
// 		amqpUrl, amqp.Config{},
// 		rabbitmq.WithConsumerOptionsLogging,
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = consumer.StartConsuming(
// 		func(d rabbitmq.Delivery) rabbitmq.Action {
// 			log.Printf("consumed: %v", string(d.Body))
// 			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
// 			return rabbitmq.Ack
// 		},
// 		"my_queue",
// 		[]string{"routing_key", "routing_key_2"},
// 		rabbitmq.WithConsumeOptionsConcurrency(10),
// 		rabbitmq.WithConsumeOptionsQueueDurable,
// 		rabbitmq.WithConsumeOptionsQuorum,
// 		rabbitmq.WithConsumeOptionsBindingExchangeName("events"),
// 		rabbitmq.WithConsumeOptionsBindingExchangeKind("topic"),
// 		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
// 		rabbitmq.WithConsumeOptionsConsumerName(consumerName),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// block main thread - wait for shutdown signal
// 	sigs := make(chan os.Signal, 1)
// 	done := make(chan bool, 1)

// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		sig := <-sigs
// 		fmt.Println()
// 		fmt.Println(sig)
// 		done <- true
// 	}()

// 	fmt.Println("awaiting signal")
// 	<-done
// 	fmt.Println("stopping consumer")

// 	// wait for server to acknowledge the cancel
// 	noWait := false
// 	consumer.StopConsuming(consumerName, noWait)
// 	consumer.Disconnect()
// }

func main() {
	// getMetadata("https://vibesmediastorage.s3.amazonaws.com/uploads/61d05316d5d1d2000f61f2d0.mp3")
	metadata.GetMetadata("https://vibesmediastorage.s3.amazonaws.com/uploads/61d05316d5d1d2000f61f2d0.mp3")
}
func getMetadata(url string) string {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, url)
	if err != nil {
		log.Panicf("Error getting data: %v", err)
	}

	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Panicf("Error unmarshalling: %v", err)
	}
	log.Print(string(buf))
	return string(buf)
}
