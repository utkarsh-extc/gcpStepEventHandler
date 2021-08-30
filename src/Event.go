// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func (e PubSubMessage) UnMarshallEvent() MyCloudEvent {
	var myEvent MyCloudEvent
	err := json.Unmarshal(e.Data, &myEvent)
	if err != nil {
		log.Println("Failed to unmarshall PubSubMessage")
		return MyCloudEvent{}
	}

	return myEvent
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
func Eventhandler(ctx context.Context, m PubSubMessage) error {

	msg := m.UnMarshallEvent()

	log.Printf("%s received by event handler", m)
	return EventRouter(ctx, msg)
}

func EventRouter(ctx context.Context, msg MyCloudEvent) error {
	switch msg.Type {
	case Step1MessageType:
		topic := c.Topic(Step2Topic)
		// use same correlationID
		msg.Type = Step2MessageType
		msg.Data = "Step2MessagePayload"
		res := topic.Publish(ctx, &pubsub.Message{Data: msg.MashallEvent()})
		id, e := res.Get(ctx)
		if e != nil {
			log.Println("err while creating pubsub client", e.Error())
			return e
		}
		log.Println("message publish on Step2Topic", id)
	case Step2MessageType:
		topic := c.Topic(Step3Topic)
		// use same correlationID
		msg.Type = Step3MessageType
		msg.Data = "Step3MessagePayload"
		res := topic.Publish(ctx, &pubsub.Message{Data: msg.MashallEvent()})
		id, e := res.Get(ctx)
		if e != nil {
			log.Println("err while creating pubsub client", e.Error())
			return e
		}
		log.Println("message publish on Step3Topic", id)
	case Step3MessageType:
		topic := c.Topic(Step4Topic)
		// use same correlationID
		msg.Type = Step4MessageType
		msg.Data = "Step4MessagePayload"
		res := topic.Publish(ctx, &pubsub.Message{Data: msg.MashallEvent()})
		id, e := res.Get(ctx)
		if e != nil {
			log.Println("err while creating pubsub client", e.Error())
			return e
		}
		log.Println("message publish on Step3Topic", id)
	case Step4MessageType:
		log.Println("msg recived of type", Step4MessageType, "has", string(msg.MashallEvent()))

	default:
		log.Printf("%s received, but not handled", msg)
	}

	return nil
}
