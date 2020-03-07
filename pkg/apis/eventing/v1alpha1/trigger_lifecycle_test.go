/*
Copyright 2019 The Knative Authors

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
	"testing"

	messagingv1alpha1 "knative.dev/eventing/pkg/apis/messaging/v1alpha1"
	"knative.dev/pkg/apis"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var (
	triggerConditionReady = apis.Condition{
		Type:   TriggerConditionReady,
		Status: corev1.ConditionTrue,
	}

	triggerConditionBroker = apis.Condition{
		Type:   TriggerConditionBroker,
		Status: corev1.ConditionTrue,
	}

	triggerConditionDependency = apis.Condition{
		Type:   TriggerConditionDependency,
		Status: corev1.ConditionTrue,
	}

	triggerConditionSubscriberResolved = apis.Condition{
		Type:   TriggerConditionSubscriberResolved,
		Status: corev1.ConditionTrue,
	}

	triggerConditionSubscribed = apis.Condition{
		Type:   TriggerConditionSubscribed,
		Status: corev1.ConditionFalse,
	}
)

func TestTriggerGetCondition(t *testing.T) {
	tests := []struct {
		name      string
		ts        *TriggerStatus
		condQuery apis.ConditionType
		want      *apis.Condition
	}{{
		name: "single condition",
		ts: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{
					triggerConditionReady,
				},
			},
		},
		condQuery: apis.ConditionReady,
		want:      &triggerConditionReady,
	}, {
		name: "multiple conditions",
		ts: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{
					triggerConditionBroker,
					triggerConditionSubscribed,
					triggerConditionDependency,
					triggerConditionSubscriberResolved,
				},
			},
		},
		condQuery: TriggerConditionSubscribed,
		want:      &triggerConditionSubscribed,
	}, {
		name: "multiple conditions, condition false",
		ts: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{
					triggerConditionBroker,
					triggerConditionSubscribed,
					triggerConditionDependency,
					triggerConditionSubscriberResolved,
				},
			},
		},
		condQuery: TriggerConditionSubscribed,
		want:      &triggerConditionSubscribed,
	}, {
		name: "unknown condition",
		ts: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{
					triggerConditionSubscribed,
				},
			},
		},
		condQuery: apis.ConditionType("foo"),
		want:      nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.ts.GetCondition(test.condQuery)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("unexpected condition (-want, +got) = %v", diff)
			}
		})
	}
}

func TestTriggerInitializeConditions(t *testing.T) {
	tests := []struct {
		name string
		ts   *TriggerStatus
		want *TriggerStatus
	}{{
		name: "empty",
		ts:   &TriggerStatus{},
		want: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   TriggerConditionBroker,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   TriggerConditionDependency,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   TriggerConditionReady,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   TriggerConditionSubscribed,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   TriggerConditionSubscriberResolved,
					Status: corev1.ConditionUnknown,
				},
				},
			},
		},
	}, {
		name: "one false",
		ts: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   TriggerConditionBroker,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		want: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   TriggerConditionBroker,
					Status: corev1.ConditionFalse,
				}, {
					Type:   TriggerConditionDependency,
					Status: corev1.ConditionUnknown,
				},
					{
						Type:   TriggerConditionReady,
						Status: corev1.ConditionUnknown,
					}, {
						Type:   TriggerConditionSubscribed,
						Status: corev1.ConditionUnknown,
					}, {
						Type:   TriggerConditionSubscriberResolved,
						Status: corev1.ConditionUnknown,
					},
				},
			},
		},
	}, {
		name: "one true",
		ts: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   TriggerConditionSubscribed,
					Status: corev1.ConditionTrue,
				}},
			},
		},
		want: &TriggerStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   TriggerConditionBroker,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   TriggerConditionDependency,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   TriggerConditionReady,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   TriggerConditionSubscribed,
					Status: corev1.ConditionTrue,
				}, {
					Type:   TriggerConditionSubscriberResolved,
					Status: corev1.ConditionUnknown,
				},
				},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.ts.InitializeConditions()
			if diff := cmp.Diff(test.want, test.ts, ignoreAllButTypeAndStatus); diff != "" {
				t.Errorf("unexpected conditions (-want, +got) = %v", diff)
			}
		})
	}
}

func TestTriggerConditionStatus(t *testing.T) {
	tests := []struct {
		name                        string
		brokerStatus                *BrokerStatus
		markKubernetesServiceExists bool
		markVirtualServiceExists    bool
		subscriptionOwned           bool
		subscriptionStatus          *messagingv1alpha1.SubscriptionStatus
		subscriberResolvedStatus    bool
		dependencyAnnotationExists  bool
		dependencyStatus            corev1.ConditionStatus
		wantConditionStatus         corev1.ConditionStatus
	}{{
		name:                        "all happy",
		brokerStatus:                TestHelper.ReadyBrokerStatus(),
		markKubernetesServiceExists: true,
		markVirtualServiceExists:    true,
		subscriptionOwned:           true,
		subscriptionStatus:          TestHelper.ReadySubscriptionStatus(),
		subscriberResolvedStatus:    true,
		dependencyAnnotationExists:  false,
		wantConditionStatus:         corev1.ConditionTrue,
	}, {
		name:                        "broker status unknown",
		brokerStatus:                TestHelper.UnknownBrokerStatus(),
		markKubernetesServiceExists: true,
		markVirtualServiceExists:    true,
		subscriptionOwned:           true,
		subscriptionStatus:          TestHelper.ReadySubscriptionStatus(),
		subscriberResolvedStatus:    true,
		dependencyAnnotationExists:  false,
		wantConditionStatus:         corev1.ConditionUnknown,
	}, {
		name:                        "broker status false",
		brokerStatus:                TestHelper.FalseBrokerStatus(),
		markKubernetesServiceExists: true,
		markVirtualServiceExists:    true,
		subscriptionOwned:           true,
		subscriptionStatus:          TestHelper.ReadySubscriptionStatus(),
		subscriberResolvedStatus:    true,
		dependencyAnnotationExists:  false,
		wantConditionStatus:         corev1.ConditionFalse,
	}, {
		name:                        "subscribed sad",
		brokerStatus:                TestHelper.ReadyBrokerStatus(),
		markKubernetesServiceExists: true,
		markVirtualServiceExists:    true,
		subscriptionOwned:           true,
		subscriptionStatus:          TestHelper.FalseSubscriptionStatus(),
		subscriberResolvedStatus:    true,
		dependencyAnnotationExists:  false,
		wantConditionStatus:         corev1.ConditionFalse,
	}, {
		name:                        "subscription not owned",
		brokerStatus:                TestHelper.ReadyBrokerStatus(),
		markKubernetesServiceExists: true,
		markVirtualServiceExists:    true,
		subscriptionOwned:           false,
		subscriptionStatus:          TestHelper.ReadySubscriptionStatus(),
		subscriberResolvedStatus:    true,
		dependencyAnnotationExists:  false,
		wantConditionStatus:         corev1.ConditionFalse,
	}, {
		name:                        "failed to resolve subscriber",
		brokerStatus:                TestHelper.ReadyBrokerStatus(),
		markKubernetesServiceExists: true,
		markVirtualServiceExists:    true,
		subscriptionOwned:           true,
		subscriptionStatus:          TestHelper.ReadySubscriptionStatus(),
		subscriberResolvedStatus:    false,
		dependencyAnnotationExists:  true,
		dependencyStatus:            corev1.ConditionTrue,
		wantConditionStatus:         corev1.ConditionFalse,
	}, {
		name:                        "dependency unknown",
		brokerStatus:                TestHelper.ReadyBrokerStatus(),
		markKubernetesServiceExists: true,
		markVirtualServiceExists:    true,
		subscriptionOwned:           true,
		subscriptionStatus:          TestHelper.ReadySubscriptionStatus(),
		subscriberResolvedStatus:    true,
		dependencyAnnotationExists:  true,
		dependencyStatus:            corev1.ConditionUnknown,
		wantConditionStatus:         corev1.ConditionUnknown,
	}, {
		name:                        "dependency false",
		brokerStatus:                TestHelper.ReadyBrokerStatus(),
		markKubernetesServiceExists: true,
		markVirtualServiceExists:    true,
		subscriptionOwned:           true,
		subscriptionStatus:          TestHelper.ReadySubscriptionStatus(),
		subscriberResolvedStatus:    true,
		dependencyAnnotationExists:  true,
		dependencyStatus:            corev1.ConditionFalse,
		wantConditionStatus:         corev1.ConditionFalse,
	}, {
		name:                        "all sad",
		brokerStatus:                TestHelper.FalseBrokerStatus(),
		markKubernetesServiceExists: false,
		markVirtualServiceExists:    false,
		subscriptionOwned:           false,
		subscriptionStatus:          TestHelper.FalseSubscriptionStatus(),
		subscriberResolvedStatus:    false,
		dependencyAnnotationExists:  true,
		dependencyStatus:            corev1.ConditionFalse,
		wantConditionStatus:         corev1.ConditionFalse,
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := &TriggerStatus{}
			if test.brokerStatus != nil {
				ts.PropagateBrokerStatus(test.brokerStatus)
			}
			if !test.subscriptionOwned {
				ts.MarkSubscriptionNotOwned(&messagingv1alpha1.Subscription{})
			} else if test.subscriptionStatus != nil {
				ts.PropagateSubscriptionStatus(test.subscriptionStatus)
			}
			if test.subscriberResolvedStatus {
				ts.MarkSubscriberResolvedSucceeded()
			} else {
				ts.MarkSubscriberResolvedFailed("Unable to get the Subscriber's URI", "subscriber not found")
			}
			if !test.dependencyAnnotationExists {
				ts.MarkDependencySucceeded()
			} else {
				if test.dependencyStatus == corev1.ConditionTrue {
					ts.MarkDependencySucceeded()
				} else if test.dependencyStatus == corev1.ConditionUnknown {
					ts.MarkDependencyUnknown("The status of dependency is unknown", "The status of dependency is unknown: nil")
				} else {
					ts.MarkDependencyFailed("The status of dependency is false", "The status of dependency is unknown: nil")
				}
			}
			got := ts.GetTopLevelCondition().Status
			if test.wantConditionStatus != got {
				t.Errorf("unexpected readiness: want %v, got %v", test.wantConditionStatus, got)
			}
		})
	}
}