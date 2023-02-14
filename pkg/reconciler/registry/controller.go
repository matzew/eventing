package registry

import (
	"context"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func NewController(ctx context.Context, cmw configmap.Watcher, ) *controller.Impl {
	logger := logging.FromContext(ctx)
	logger.Infow("hello.")

	return nil
}
