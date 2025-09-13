package client_k8s

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"ingress-test-suite/internal/consts"
	networkingv1 "k8s.io/api/networking/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Logger interface {
	Infof(format string, args ...interface{})
}

type K8SClient struct {
	logger       Logger
	k8sClientSet *kubernetes.Clientset
}

func NewK8SClient(logger Logger) (*K8SClient, error) {
	client := &K8SClient{
		logger: logger,
	}

	k8sClientSet, err := client.makeClientSet()
	if err != nil {
		return nil, err
	}
	client.k8sClientSet = k8sClientSet
	client.logger.Infof(consts.MessageCreteK8SClient, client.k8sClientSet.AppsV1())

	return client, nil
}

func (c *K8SClient) IngressCreate(ctx context.Context, namespace string, ingress *networkingv1.Ingress, opts metav1.CreateOptions) (*networkingv1.Ingress, error) {
	return c.k8sClientSet.NetworkingV1().Ingresses(namespace).
		Create(ctx, ingress, opts)
}

func (c *K8SClient) IngressDelete(ctx context.Context, namespace string, ingressName string, opts metav1.DeleteOptions) error {
	return c.k8sClientSet.NetworkingV1().Ingresses(namespace).
		Delete(ctx, ingressName, opts)
}

func (c *K8SClient) IngressGet(ctx context.Context, namespace string, ingressName string, opts metav1.GetOptions) (*networkingv1.Ingress, error) {
	return c.k8sClientSet.NetworkingV1().Ingresses(namespace).
		Get(ctx, ingressName, opts)
}

func (c *K8SClient) makeClientSet() (*kubernetes.Clientset, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		var kubeConfig string
		if home := homedir.HomeDir(); home != "" {
			kubeConfig = filepath.Join(home, ".kube", "config")
		} else if os.Getenv("KUBECONFIG") != "" {
			kubeConfig = os.Getenv("KUBECONFIG")
		} else {
			return nil, errors.Wrap(err, consts.ErrLoadKubeConfig.Error())
		}

		cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, errors.Wrap(err, consts.ErrBuildKubeConfig.Error())
		}
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, consts.ErrFailedCreateKubeClient.Error())
	}

	return clientSet, nil
}
