package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/teezzan/ffmpeg-worker/pkg/metadata"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

var consumerName = "example"
var amqpUrl = os.Getenv("AMQPURL")

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
			log.Printf("consumed: %v", string(d.Body))
			metadata.GetMetadata("https://vibesmediastorage.s3.amazonaws.com/uploads/61d05316d5d1d2000f61f2d0.mp3")
			return rabbitmq.Ack
		},
		"ffprobeQueue",
		[]string{"testKey"},
		// rabbitmq.WithConsumeOptionsConcurrency(5),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
		rabbitmq.WithConsumeOptionsBindingExchangeName("events"),
		rabbitmq.WithConsumeOptionsBindingExchangeKind("topic"),
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
