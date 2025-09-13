package ingress_manager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"ingress-test-suite/internal/consts"
	ds "ingress-test-suite/internal/datastruct"
	networkingv1 "k8s.io/api/networking/v1"
	k8s_machinery "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sClient interface {
	IngressCreate(ctx context.Context, namespace string, ingress *networkingv1.Ingress, opts metav1.CreateOptions) (*networkingv1.Ingress, error)
	IngressDelete(ctx context.Context, namespace string, ingressName string, opts metav1.DeleteOptions) error
	IngressGet(ctx context.Context, namespace string, ingressName string, opts metav1.GetOptions) (*networkingv1.Ingress, error)
}

type K8sIngressManager struct {
	client    *kubernetes.Clientset
	k8sClient K8sClient
}

func NewManager(client K8sClient) *K8sIngressManager {
	return &K8sIngressManager{
		k8sClient: client,
	}
}

func (m *K8sIngressManager) Create(ctx context.Context, entry *ds.IngressTestEntry, c *ds.IngressTestsFile) error {
	ingress, err := createIngressRule(entry, c)
	if err != nil {
		return err
	}

	_, err = m.k8sClient.IngressCreate(ctx, entry.Namespace, ingress, metav1.CreateOptions{})
	return err
}

func (m *K8sIngressManager) Delete(ctx context.Context, entry *ds.IngressTestEntry) error {
	return m.k8sClient.IngressDelete(ctx, entry.Namespace, fmt.Sprintf("test-%s", entry.Host), metav1.DeleteOptions{})
}

func (m *K8sIngressManager) CheckExist(ctx context.Context, entry *ds.IngressTestEntry) (bool, error) {
	ingressName := fmt.Sprintf("test-%s", entry.Host)

	_, err := m.k8sClient.IngressGet(ctx, entry.Namespace, ingressName, metav1.GetOptions{})
	if err != nil {
		if k8s_machinery.IsNotFound(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, consts.ErrFailedCheckExistIngressRule.Error(), ingressName)
	}

	return true, nil
}

func createIngressRule(t *ds.IngressTestEntry, c *ds.IngressTestsFile) (*networkingv1.Ingress, error) {
	pathType, err := convertPathType(t.PathType)
	if err != nil {
		return nil, errors.Wrap(err, consts.ErrFailedConvertPathType.Error())
	}

	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("test-%s", t.Host),
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &c.IngressClassName,
			Rules: []networkingv1.IngressRule{
				{
					Host: t.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     t.Path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: t.Service,
											Port: networkingv1.ServiceBackendPort{
												Number: int32(t.Port),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func convertPathType(s string) (networkingv1.PathType, error) {
	switch s {
	case "Exact":
		return networkingv1.PathTypeExact, nil
	case "Prefix":
		return networkingv1.PathTypePrefix, nil
	case "ImplementationSpecific":
		return networkingv1.PathTypeImplementationSpecific, nil
	default:
		return "", errors.Errorf(consts.ErrInvalidPathType.Error(), s)
	}
}
