package registry

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/eventing/pkg/apis/eventing/v1beta1"
	eventingclientset "knative.dev/eventing/pkg/client/clientset/versioned"
	"strings"

	"log"
)

type payload struct {
	Sourcename  string   `json:"sourcename"`
	Description string   `json:"description"`
	Events      []string `json:"events"`
}

// Handler is the HTTP handler for the registry.
type Handler struct {
	client         client.Client
	eventingClient *eventingclientset.Clientset
}

// NewHandler creates a new registry handler.
func NewHandler(ctx context.Context, eventingClient *eventingclientset.Clientset) (*Handler, error) {

	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &Handler{
		client:         client,
		eventingClient: eventingClient,
	}, nil

}

func (h *Handler) Start(ctx context.Context) error {
	return h.client.StartReceiver(ctx, h.receive)
}

func (h *Handler) receive(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	data := &payload{}
	if err := event.DataAs(data); err != nil {
		log.Printf("Got error while unmarshalling data: %s", err.Error())
		return nil, http.NewResult(400, "got error while unmarshalling data: %w", err)
	}

	//TODO fix me: we should not hardcode the namespace
	systemNamespace := "knative-eventing"

	for _, event := range data.Events {
		et := &v1beta1.EventType{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.ToLower(data.Sourcename) + "-" + strings.ToLower(event),
				Namespace: systemNamespace,
			},
			Spec: v1beta1.EventTypeSpec{
				Type:        event,
				Source:      nil,
				Schema:      nil,
				SchemaData:  "",
				Broker:      "",
				Description: data.Description,
			},
		}
		_, err := h.eventingClient.EventingV1beta1().EventTypes(systemNamespace).Create(ctx, et, metav1.CreateOptions{})
		if err != nil {
			log.Printf("Got error: %s", err.Error())
		}
	}
	return nil, nil
}
