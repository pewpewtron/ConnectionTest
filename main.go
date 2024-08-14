package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/couchbase/gocb/v2"
)

func pingHost(host string) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func main() {
	// Retrieve connection strings from environment variables
	couchbaseConnStr := os.Getenv("COUCHBASE_CONN_STR")
	confluentConnStr := os.Getenv("CONFLUENT_CONN_STR")

	if couchbaseConnStr == "" || confluentConnStr == "" {
		log.Fatal("Environment variables COUCHBASE_CONN_STR and CONFLUENT_CONN_STR must be set")
	}

	// Print connection strings
	fmt.Printf("Couchbase connection string: %s\n", couchbaseConnStr)
	fmt.Printf("Confluent connection string: %s\n", confluentConnStr)

	// Check if connection is resolve from dns
	// Extract hostnames from connection strings
	couchbaseHost := strings.Split(strings.TrimPrefix(couchbaseConnStr, "couchbase://"), ",")[0]
	confluentHost := strings.Split(confluentConnStr, ",")[0]

	// Check DNS resolution
	if _, err := net.LookupHost(couchbaseHost); err != nil {
		log.Fatalf("Failed to resolve Couchbase host: %v", err)
	}
	if _, err := net.LookupHost(confluentHost); err != nil {
		log.Fatalf("Failed to resolve Confluent host: %v", err)
	}

	// Ping hosts
	if err := pingHost(couchbaseHost); err != nil {
		log.Fatalf("Failed to ping Couchbase host: %v", err)
	}
	if err := pingHost(confluentHost); err != nil {
		log.Fatalf("Failed to ping Confluent host: %v", err)
	}

	// Check Couchbase connectivity
	cluster, err := gocb.Connect(couchbaseConnStr, gocb.ClusterOptions{
		Username: os.Getenv("COUCHBASE_USER"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	})
	if err != nil {
		log.Printf("Failed to connect to Couchbase: %v", err)
	} else {
		defer cluster.Close(nil)

		err = cluster.WaitUntilReady(0, nil)
		if err != nil {
			log.Printf("Couchbase cluster not ready: %v", err)
		} else {
			fmt.Println("Successfully connected to Couchbase")
		}
	}

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
