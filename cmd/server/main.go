package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to RabbitMQ successfully.")

	newChannel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer newChannel.Close()

	fmt.Println("Channel created successfully.")

	err = pubsub.PublishJSON(
		newChannel, 
		routing.ExchangePerilDirect, 
		string(routing.PauseKey),
		routing.PlayingState{
			IsPaused: true,
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Published initial paused state.")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Shutting down Peril server...")
	conn.Close()
	fmt.Println("Peril server stopped.")
}
