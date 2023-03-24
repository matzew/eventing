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
	"knative.dev/pkg/apis"
	"strings"

	"log"
)

type Group struct {
	ID          string       `json:"id"`
	Version     int          `json:"version"`
	Definitions []Definition `json:"definitions"`
}

type Definition struct {
	ID          string   `json:"id"`
	Version     int      `json:"version"`
	Description string   `json:"description"`
	SchemaURL   string   `json:"schemaUrl"`
	Format      string   `json:"format"`
	Metadata    Metadata `json:"metadata"`
}

type Metadata struct {
	Attributes map[string]Attribute `json:"attributes"`
}

type Attribute struct {
	Required bool        `json:"required"`
	Value    interface{} `json:"value,omitempty"`
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
	data := &Group{}
	if err := event.DataAs(data); err != nil {
		log.Printf("Got error while unmarshalling data: %s", err.Error())
		return nil, http.NewResult(400, "got error while unmarshalling data: %w", err)
	}

	//TODO fix me: we should not hardcode the namespace
	systemNamespace := "knative-eventing"

	for _, definition := range data.Definitions {

		et := &v1beta1.EventType{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.ToLower(data.ID) + "-" + strings.ToLower(definition.ID),
				Namespace: systemNamespace,
			},
			Spec: v1beta1.EventTypeSpec{
				Type:        definition.ID,
				Source:      readURLValue(definition.Metadata.Attributes["source"]),
				Schema:      readURLValue(definition.Metadata.Attributes["dataschema"]),
				Broker:      "",
				Description: definition.Description,
			},
		}
		_, err := h.eventingClient.EventingV1beta1().EventTypes(systemNamespace).Create(ctx, et, metav1.CreateOptions{})
		if err != nil {
			log.Printf("Got error: %s", err.Error())
		}
	}
	return nil, nil
}

func readURLValue(attr Attribute) *apis.URL {
	url, err := apis.ParseURL(attr.Value.(string))
	if err != nil {
		return nil
	}
	return url
}
