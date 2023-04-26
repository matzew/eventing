package v1alpha1

import "knative.dev/pkg/apis"

var eventTypeCondSet = apis.NewLivingConditionSet(EventTypeDefinitionConditionReady)

const (
	EventTypeDefinitionConditionReady = apis.ConditionReady
)

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*EventTypeDefinition) GetConditionSet() apis.ConditionSet {
	return eventTypeCondSet
}
