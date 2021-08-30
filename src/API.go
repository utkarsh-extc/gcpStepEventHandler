package p

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
)

type MyCloudEvent struct {
	CorrelationID string `json:"correlationID"`
	Data          string `json:"data"`
	Type          string `json:"type"`
}

func (msg MyCloudEvent) MashallEvent() []byte {
	blob, err := json.Marshal(msg)
	if err != nil {
		log.Println("Failed to Mashall MyCloudEvent")
		return []byte{}
	}
	return blob
}

func Send(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	c, e := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}

	topic := c.Topic(Step1Topic)
	msg := MyCloudEvent{CorrelationID: uuid.NewString()}
	msg.Type = Step1MessageType
	msg.Data = "Step1MessagePayload"
	res := topic.Publish(ctx, &pubsub.Message{Data: msg.MashallEvent()})

	id, e := res.Get(ctx)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}
