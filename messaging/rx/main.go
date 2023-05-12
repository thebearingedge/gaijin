package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
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
	subscriptionID, ok := os.LookupEnv("GOOGLE_PUB_SUB_SUBSCRIPTION")

	if !ok {
		log.Fatal("GOOGLE_PUB_SUB_SUBSCRIPTION env var not set")
	}

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer client.Close()

	subscription := client.Subscription(subscriptionID)

	subscription.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		var message Message

		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Fatalf("could not deserialize message data %v", string(msg.Data))
		} else {
			log.Printf("Got message: %v\n", message)
			msg.Ack()
		}
	})

}
