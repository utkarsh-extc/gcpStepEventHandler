package api

import (
	"context"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/utkarsh-extc/gcpStepEventHandler/internal"
)

func Send(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	c, e := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}

	topic := c.Topic(internal.Step1Topic)

	res := topic.Publish(ctx, &pubsub.Message{Data: []byte(internal.Step1Message)})

	id, e := res.Get(ctx)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}
