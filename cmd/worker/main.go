package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"

	// "github.com/teezzan/ffmpeg-worker/pkg/metadata"
	"github.com/teezzan/ffmpeg-worker/pkg/metadata"
	"github.com/teezzan/ffmpeg-worker/pkg/redis"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

var consumerName = "example"
var amqpUrl = "amqp://localhost" //os.Getenv("AMQPURL")

// var queueName = os.Getenv("QUEUE_NAME")

func main() {
	consumer, err := rabbitmq.NewConsumer(
		amqpUrl, amqp.Config{},
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = consumer.StartConsuming(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			var payload redis.Payload
			json.Unmarshal(d.Body, &payload)
			result := metadata.GetMetadata(payload.Url)
			redis.SaveResult(payload, result)
			return rabbitmq.Ack
		},
		"queuwnn",
		[]string{"key"},
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
		rabbitmq.WithConsumeOptionsBindingExchangeName("direct_xch"),
		rabbitmq.WithConsumeOptionsBindingExchangeKind("direct"),
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsConsumerName(consumerName),
	)
	if err != nil {
		log.Fatal(err)
	}

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("stopping consumer")

	// wait for server to acknowledge the cancel
	noWait := false
	consumer.StopConsuming(consumerName, noWait)
	consumer.Disconnect()
}
