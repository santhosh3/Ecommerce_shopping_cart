package kafka

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	config "github.com/santhosh3/ECOM/Config"
)

func KafaConsumer() {
	topic := "product"
	port := config.Envs.KafkaPort
	worker, err := ConnectConsumer([]string{fmt.Sprintf("localhost:%s", port)})
	if err != nil {
		panic(err)
	}

	// Calling ConsumePartition. It will open one connection per broker
	// and share it for all partitions that live on it.
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintf("Kafka connected on port %s", config.Envs.KafkaPort);
	log.Println(msg)

	// signal hndling for gracefull shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Count how many message processed
	msgCount := 0

	// Channel to signal the goroutine to stop
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))
			case <-sigchan:
				fmt.Println("Interrupt is detected. Shutting down consumer...")
				if err := consumer.Close(); err != nil {
					fmt.Printf("Error closing consumer: %v\n", err)
				}
				return
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}

}

func ConnectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create new consumer
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}