package resources

import (
	"fmt"
	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/eventing/pkg/apis/sinks/v1alpha1"
)

func MakeCertificate(sink *v1alpha1.IntegrationSink) *cmv1.Certificate {
	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeploymentName(sink) + "-cert",
			Namespace: sink.Namespace,
		},
		Spec: cmv1.CertificateSpec{
			SecretName: DeploymentName(sink) + "-cert",
			DNSNames: []string{
				fmt.Sprintf("%s.%s.svc", DeploymentName(sink), sink.Namespace),
			},
			IssuerRef: cmmeta.ObjectReference{
				Name:  "knative-eventing-ca-issuer",
				Kind:  "ClusterIssuer",
				Group: "cert-manager.io",
			},
		},
	}
}
