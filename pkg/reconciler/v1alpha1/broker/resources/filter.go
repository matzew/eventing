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

package resources

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type FilterArgs struct {
	Broker             *eventingv1alpha1.Broker
	Image              string
	ServiceAccountName string
}

func MakeFilterDeployment(args *FilterArgs) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: args.Broker.Namespace,
			Name:      fmt.Sprintf("%s-broker-filter", args.Broker.Name),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(args.Broker, schema.GroupVersionKind{
					Group:   eventingv1alpha1.SchemeGroupVersion.Group,
					Version: eventingv1alpha1.SchemeGroupVersion.Version,
					Kind:    "Broker",
				}),
			},
			Labels: filterLabels(args.Broker),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: filterLabels(args.Broker),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: filterLabels(args.Broker),
					Annotations: map[string]string{
						"sidecar.istio.io/inject": "true",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: args.ServiceAccountName,
					Containers: []corev1.Container{
						{
							Name:  "filter",
							Image: args.Image,
							Env: []corev1.EnvVar{
								{
									Name: "NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func MakeFilterService(b *eventingv1alpha1.Broker) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: b.Namespace,
			Name:      fmt.Sprintf("%s-broker-filter", b.Name),
			Labels:    filterLabels(b),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(b, schema.GroupVersionKind{
					Group:   eventingv1alpha1.SchemeGroupVersion.Group,
					Version: eventingv1alpha1.SchemeGroupVersion.Version,
					Kind:    "Broker",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: filterLabels(b),
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}
}

func filterLabels(b *eventingv1alpha1.Broker) map[string]string {
	return map[string]string{
		"eventing.knative.dev/broker":     b.Name,
		"eventing.knative.dev/brokerRole": "filter",
	}
}
