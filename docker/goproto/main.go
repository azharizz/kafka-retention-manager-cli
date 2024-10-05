package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
)

// Data represents the structure of your JSON message
type Data struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

var ctx = context.Background()

func generateRandomID() string {
	return strconv.Itoa(rand.Intn(1000)) // Generates a random ID between 0 and 999
}

func generateRandomName() string {
	names := []string{"Azhar", "Bob", "Charlie", "Diana", "Edward", "Fiona"}
	return names[rand.Intn(len(names))] // Randomly select a name from the slice
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Configure the Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	// Configure the Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
	defer rdb.Close()

	topic := "test_topic"

	// Set up a ticker to run every 5 minutes
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Prepare a JSON message with random ID and Name
			data := Data{
				ID:        generateRandomID(),
				Name:      generateRandomName(),
				Timestamp: time.Now().Unix(),
			}

			// Serialize the JSON message
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Failed to marshal JSON: %s", err)
				continue
			}

			// Produce the message
			deliveryChan := make(chan kafka.Event)
			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          jsonData,
			}, deliveryChan)

			if err != nil {
				log.Printf("Failed to produce message: %s", err)
				continue
			}

			go func() {
				for e := range deliveryChan {
					switch ev := e.(type) {
					case *kafka.Message:
						if ev.TopicPartition.Error != nil {
							log.Printf("Delivery failed: %v\n", ev.TopicPartition)
						} else {
							log.Printf("Successfully delivered message to %v\n", ev.TopicPartition)

							// Get the current date in YYYY-MM-DD format
							today := time.Now().Format("2006-01-02")

							// Increment the count in Redis using today's date as the key
							redisKey := fmt.Sprintf("datatest_ingestion_count:%s", today) //PLEASE SETTING OTHER THAN THIS
							fmt.Println(redisKey)
							err := rdb.Incr(ctx, redisKey).Err()
							if err != nil {
								log.Printf("Failed to increment Redis count: %s", err)
							} else {
								log.Printf("Ingestion count for %s updated successfully.", today)
							}
						}
					}
				}
			}()

			// Wait for the message to be delivered
			p.Flush(15 * 1000)
		}
	}
}
