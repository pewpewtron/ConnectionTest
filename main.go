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

	// Extract hostnames from connection strings
	couchbaseHost := strings.Split(strings.TrimPrefix(couchbaseConnStr, "couchbase://"), ",")[0]
	confluentHosts := strings.Split(confluentConnStr, ",")

	// Check DNS resolution
	fmt.Println("Checking DNS resolution...")
	if _, err := net.LookupHost(couchbaseHost); err != nil {
		log.Printf("Failed to resolve Couchbase host: %v", err)
	} else {
		fmt.Println("Successfully resolved Couchbase host")
	}

	for _, host := range confluentHosts {
		hostname := strings.Split(host, ":")[0]
		if _, err := net.LookupHost(hostname); err != nil {
			log.Printf("Failed to resolve Confluent host %s: %v", hostname, err)
		} else {
			fmt.Printf("Successfully resolved Confluent host %s\n", hostname)
		}
	}

	// Ping hosts
	fmt.Println("Pinging hosts...")
	couchbaseHostWithPort := couchbaseHost + ":8091"
	if err := pingHost(couchbaseHostWithPort); err != nil {
		log.Printf("Failed to ping Couchbase host: %v", err)
	} else {
		fmt.Println("Successfully pinged Couchbase host")
	}

	for _, host := range confluentHosts {
		if err := pingHost(host); err != nil {
			log.Printf("Failed to ping Confluent host %s: %v", host, err)
		} else {
			fmt.Printf("Successfully pinged Confluent host %s\n", host)
		}
	}

	// Check Couchbase connectivity
	fmt.Println("Checking Couchbase connectivity...")
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
	fmt.Println("Checking Confluent connectivity...")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": confluentConnStr})
	if err != nil {
		log.Printf("Failed to create Confluent producer: %v", err)
	} else {
		defer producer.Close()

		// Produce a test message to check connectivity
		testTopic := "test_topic"
		testMessage := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &testTopic, Partition: kafka.PartitionAny},
			Value:          []byte("test message"),
		}

		err = producer.Produce(testMessage, nil)
		if err != nil {
			log.Printf("Failed to produce test message to Confluent: %v", err)
		} else {
			// Wait for message deliveries
			producer.Flush(15 * 1000)
			fmt.Println("Successfully connected to Confluent")
		}
	}
}
