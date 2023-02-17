package main

import (
	"go.uber.org/zap"
	eventingclientset "knative.dev/eventing/pkg/client/clientset/versioned"
	"knative.dev/eventing/pkg/registry"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/signals"
	"log"
)

func main() {
	ctx := signals.NewContext()
	cfg := injection.ParseAndGetRESTConfigOrDie()
	eventingClient := eventingclientset.NewForConfigOrDie(cfg)

	handler, err := registry.NewHandler(ctx, eventingClient)

	err = handler.Start(ctx)
	if err != nil {
		log.Fatal("handler.Start() returned an error", zap.Error(err))
	}
}
