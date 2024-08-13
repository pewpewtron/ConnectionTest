package main

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/couchbase/gocb/v2"
)

func main() {
	// Retrieve connection strings from environment variables
	couchbaseConnStr := os.Getenv("COUCHBASE_CONN_STR")
	confluentConnStr := os.Getenv("CONFLUENT_CONN_STR")

	if couchbaseConnStr == "" || confluentConnStr == "" {
		log.Fatal("Environment variables COUCHBASE_CONN_STR and CONFLUENT_CONN_STR must be set")
	}

	// Check Couchbase connectivity
	cluster, err := gocb.Connect(couchbaseConnStr, gocb.ClusterOptions{
		Username: os.Getenv("COUCHBASE_USER"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("Failed to connect to Couchbase: %v", err)
	}
	defer cluster.Close(nil)

	err = cluster.WaitUntilReady(0, nil)
	if err != nil {
		log.Fatalf("Couchbase cluster not ready: %v", err)
	}
	fmt.Println("Successfully connected to Couchbase")

	// Check Confluent connectivity
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": confluentConnStr})
	if err != nil {
		log.Fatalf("Failed to create Confluent producer: %v", err)
	}
	defer producer.Close()

	// Produce a test message to check connectivity
	testTopic := "test_topic"
	testMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &testTopic, Partition: kafka.PartitionAny},
		Value:          []byte("test message"),
	}

	err = producer.Produce(testMessage, nil)
	if err != nil {
		log.Fatalf("Failed to produce test message to Confluent: %v", err)
	}

	// Wait for message deliveries
	producer.Flush(15 * 1000)
	fmt.Println("Successfully connected to Confluent")
}
