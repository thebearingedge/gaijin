package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/goccy/go-json"
)

type Message struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

func main() {

	projectID, ok := os.LookupEnv("GOOGLE_PROJECT_ID")

	if !ok {
		log.Fatal("GOOGLE_PROJECT_ID env var not set")
	}

	topicID, ok := os.LookupEnv("GOOGLE_PUB_SUB_TOPIC")

	if !ok {
		log.Fatal("GOOGLE_PUB_SUB_TOPIC env var not set")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer client.Close()

	topic := client.Topic(topicID)

	for {
		time.Sleep(3 * time.Second)
		msg := Message{
			Text: "hello, pub/sub!",
			Time: time.Now(),
		}
		data, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("pubsub: result.Get: %v", err)
		}
		result := topic.Publish(ctx, &pubsub.Message{
			Data: data,
		})
		if id, err := result.Get(ctx); err != nil {
			log.Fatalf("pubsub: result.Get: %v", err)
		} else {
			log.Printf("Published a message; msg ID: %v\n", id)
		}
	}
}
