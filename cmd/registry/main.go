package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/eventing/pkg/apis/eventing/v1beta1"
	eventingclientset "knative.dev/eventing/pkg/client/clientset/versioned"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/signals"
	"log"

	"github.com/cloudevents/sdk-go/v2/event"
)

func main() {
	run(signals.NewContext())
}

type payload struct {
	Message string `json:"message"`
}

func run(ctx context.Context) {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	cfg := injection.ParseAndGetRESTConfigOrDie()
	eventingClient := eventingclientset.NewForConfigOrDie(cfg)

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, func(ctx context.Context, event cloudevents.Event) (*event.Event, protocol.Result) {
		data := &payload{}
		if err := event.DataAs(data); err != nil {
			log.Printf("Got error while unmarshalling data: %s", err.Error())
			return nil, http.NewResult(400, "got error while unmarshalling data: %w", err)
		}

		knativeeventing := "knative-eventing"

		et := &v1beta1.EventType{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.New().String(),
				Namespace: knativeeventing,
			},
			Spec: v1beta1.EventTypeSpec{
				Type:        data.Message,
				Source:      nil,
				Schema:      nil,
				SchemaData:  "",
				Broker:      "my-default-broker",
				Description: "",
			},
		}

		_, err := eventingClient.EventingV1beta1().EventTypes(knativeeventing).Create(ctx, et, metav1.CreateOptions{})

		if err != nil {
			log.Printf("Got error: %s", err.Error())
		}

		return nil, nil
	}))

}
