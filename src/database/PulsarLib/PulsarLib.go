package PulsarLib

import (
    "log"
    "time"
	"fmt"
	"context"
    "github.com/apache/pulsar-client-go/pulsar"
)

//Default configuration URL: "pulsar://localhost:6650"
func InitClient(URL string) *pulsar.Client{
	client, err := pulsar.NewClient(pulsar.ClientOptions{
        URL:               URL,
        OperationTimeout:  30 * time.Second,
        ConnectionTimeout: 30 * time.Second,
    })
    if err != nil {
        log.Fatalf("Could not instantiate Pulsar client: %v", err)
    }
	return &client
}

func CreateConsumer(client *pulsar.Client, topic string, sub_name string) *pulsar.Consumer{
	consumer, err := (*client).Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: sub_name,
		Type:             pulsar.Shared,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &consumer
}

func CreateProducer(client *pulsar.Client, topic string) *pulsar.Producer{
	producer, err := (*client).CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	
	if err != nil {
		log.Fatal(err)
	}
	return &producer
}

func SendMessage(producer *pulsar.Producer, message []byte){
	_, err := (*producer).Send(context.Background(), &pulsar.ProducerMessage{
		Payload: message,
	})

	if err != nil {
		fmt.Println("Failed to publish message", err)
	}
	fmt.Println("Published message")
}



