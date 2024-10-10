/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	v1 "knative.dev/eventing/pkg/apis/sources/v1"
	"knative.dev/pkg/apis"
)

const (
	// IntegrationSourceConditionReady has status True when the IntegrationSource is ready to send events.
	IntegrationSourceConditionReady = apis.ConditionReady

	// IntegrationSourceConditionSinkProvided has status True when the ApiServerSource has been configured with a sink target.
	IntegrationSourceConditionSinkProvided apis.ConditionType = "SinkProvided"

	// IntegrationSourceConditionContainerSourceReady has status True when the IntegrationSource's ContainerSource is ready.
	IntegrationSourceConditionContainerSourceReady apis.ConditionType = "ContainerSourceReady"
)

var IntegrationCondSet = apis.NewLivingConditionSet(
	IntegrationSourceConditionContainerSourceReady,
	//	IntegrationSourceConditionSinkProvided,
)

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*IntegrationSource) GetConditionSet() apis.ConditionSet {
	return IntegrationCondSet
}

// / GetCondition returns the condition currently associated with the given type, or nil.
func (i *IntegrationSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return IntegrationCondSet.Manage(i).GetCondition(t)
}

// GetTopLevelCondition returns the top level condition.
func (i *IntegrationSourceStatus) GetTopLevelCondition() *apis.Condition {
	return IntegrationCondSet.Manage(i).GetTopLevelCondition()
}

// IsReady returns true if the resource is ready overall.
func (i *IntegrationSourceStatus) IsReady() bool {
	return IntegrationCondSet.Manage(i).IsHappy()
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (i *IntegrationSourceStatus) InitializeConditions() {
	IntegrationCondSet.Manage(i).InitializeConditions()
}

func (i *IntegrationSourceStatus) MarkSink(uri *apis.URL) {
	i.SinkURI = uri
	if len(uri.String()) > 0 {
		IntegrationCondSet.Manage(i).MarkTrue(IntegrationSourceConditionSinkProvided)
	} else {
		IntegrationCondSet.Manage(i).MarkFalse(IntegrationSourceConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.%s", "")
	}
}

func (iss *IntegrationSourceStatus) MarkNoSink(reason, messageFormat string, messageA ...interface{}) {
	IntegrationCondSet.Manage(iss).MarkFalse(IntegrationSourceConditionSinkProvided, reason, messageFormat, messageA...)
}

func (i *IntegrationSourceStatus) PropagateContainerSourcueStatus(status *v1.ContainerSourceStatus) {
	// Do not copy conditions nor observedGeneration
	conditions := i.Conditions
	observedGeneration := i.ObservedGeneration
	i.SourceStatus = status.SourceStatus
	i.Conditions = conditions
	i.ObservedGeneration = observedGeneration

	cond := status.GetCondition(apis.ConditionReady)
	switch {
	case cond == nil:
		IntegrationCondSet.Manage(i).MarkUnknown(IntegrationSourceConditionContainerSourceReady, "", "")
	case cond.Status == corev1.ConditionTrue:
		IntegrationCondSet.Manage(i).MarkTrue(IntegrationSourceConditionContainerSourceReady)
	case cond.Status == corev1.ConditionFalse:
		IntegrationCondSet.Manage(i).MarkFalse(IntegrationSourceConditionContainerSourceReady, cond.Reason, cond.Message)
	case cond.Status == corev1.ConditionUnknown:
		IntegrationCondSet.Manage(i).MarkUnknown(IntegrationSourceConditionContainerSourceReady, cond.Reason, cond.Message)
	default:
		IntegrationCondSet.Manage(i).MarkUnknown(IntegrationSourceConditionContainerSourceReady, cond.Reason, cond.Message)
	}

	// Propagate containersources AuthStatus to integrationrsources AuthStatus
	i.Auth = status.Auth
}
