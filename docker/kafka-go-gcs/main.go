package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Kafka consumer configuration
	fmt.Println("NewConsumer")
	log.Printf("New consumer logs")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka1:19092",
		"group.id":          "test_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}
	defer consumer.Close()

	// Subscribe to Kafka topi
	fmt.Println("Subscript topic")
	consumer.SubscribeTopics([]string{"test_topic"}, nil)

	// Google Cloud Storage client
	ctx := context.Background()

	fmt.Println("Connecting to storage")
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket("tf_gke_bucket2")

	fmt.Println("Connected to bucket")

	for {
		fmt.Println("Read Message")
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			// Save message to GCS
			fmt.Println("Getting time now")
			currentDate := time.Now().Format("2006-01-02")
			fmt.Println("Sending to bucket")
			wc := bucket.Object(fmt.Sprintf("messages/%s/%s.json", currentDate, msg.Timestamp)).NewWriter(ctx)
			//wc := bucket.Object(fmt.Sprintf("messages/%s.json", msg.Timestamp)).NewWriter(ctx)
			if _, err := wc.Write(msg.Value); err != nil {
				log.Printf("Failed to write message to GCS: %v", err)
			}
			if err := wc.Close(); err != nil {
				log.Printf("Failed to close writer: %v", err)
			}
		} else {
			log.Printf("Error while consuming message: %v\n", err)
		}
	}
}
