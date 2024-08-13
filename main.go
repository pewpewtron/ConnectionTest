package main

import (
	"fmt"
	"log"
	"os"

	"github.com/couchbase/gocb/v2"
)

// getEnv retrieves the value of the environment variable named by the key.
// It returns the value, or the default value if the variable is not present.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	// Retrieve environment variables
	bucketName := getEnv("COUCHBASE_BUCKET", "default_bucket")
	username := getEnv("COUCHBASE_USERNAME", "default_username")
	password := getEnv("COUCHBASE_PASSWORD", "default_password")
	serverAddress := getEnv("COUCHBASE_SERVER_ADDRESS", "localhost")

	// Connect to the Couchbase c    go mod download github.com/couchbase/gocb/v2luster
	cluster, err := gocb.Connect(fmt.Sprintf("couchbase://%s", serverAddress), gocb.ClusterOptions{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Couchbase cluster: %v", err)
	}

	// Open the bucket
	bucket := cluster.Bucket(bucketName)
	err = bucket.WaitUntilReady(5, nil)
	if err != nil {
		log.Fatalf("Failed to open bucket: %v", err)
	}

	// Create a N1QL query to fetch the first 3 documents
	query := fmt.Sprintf("SELECT * FROM `%s` LIMIT 3", bucketName)

	// Execute the query
	rows, err := cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	// Iterate through the rows and print the results
	for rows.Next() {
		var doc interface{}
		err := rows.Row(&doc)
		if err != nil {
			log.Fatalf("Failed to read row: %v", err)
		}
		fmt.Printf("Document: %v\n", doc)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Fatalf("Query execution error: %v", err)
	}
}
