// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/cloudevents/sdk-go/v2/event"
)

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var c *pubsub.Client

func init() {
	var err error
	c, err = pubsub.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Println("err while creating pubsub client", err.Error())
	}
}

// Eventhandler consumes a CloudEvent message and extracts the Pub/Sub message.
func Eventhandler(ctx context.Context, e event.Event) error {

	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	msgBlob := string(msg.Message.Data) // Automatically decoded from base64.

	log.Printf("%s received by event handler", msgBlob)
	return EventRouter(ctx,msgBlob)
}

func EventRouter(ctx context.Context,msgBlob string)error{
	switch msgBlob {
	case Step1Message:
		topic := c.Topic(Step2Topic)
		res := topic.Publish(ctx, &pubsub.Message{Data: []byte(Step2Message)})
		id, e := res.Get(ctx)
		if e != nil {
			log.Println("err while creating pubsub client", e.Error())
			return e
		}
		log.Println("message publish on Step2Topic", id)
	case Step2Message:
		topic := c.Topic(Step3Topic)
		res := topic.Publish(ctx, &pubsub.Message{Data: []byte(Step3Message)})
		id, e := res.Get(ctx)
		if e != nil {
			log.Println("err while creating pubsub client", e.Error())
			return e
		}
		log.Println("message publish on Step3Topic", id)

	default:
		log.Printf("%s received, but not handled", msgBlob)
	}

	return nil
}
