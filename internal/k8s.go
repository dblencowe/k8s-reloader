package internal

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func MakeK8Client(isDev bool) (*K8Client, error) {
	cs, err := makeK8ClientSet(isDev)
	if err != nil {
		return nil, err
	}
	return &K8Client{
		clientset: cs,
	}, nil
}

type K8Client struct {
	clientset *kubernetes.Clientset
}

func (c *K8Client) RestartDeployment(ctx context.Context, namespace, deployment string) error {
	d, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, deployment, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if d.Spec.Template.ObjectMeta.Annotations == nil {
		d.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	d.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	_, err = c.clientset.AppsV1().Deployments(namespace).Update(ctx, d, metav1.UpdateOptions{})
	return err
}

func makeK8sConfig(isDev bool) (*rest.Config, error) {
	var err error
	var config *rest.Config

	if isDev {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		config, err = kubeConfig.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	if config == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func makeK8ClientSet(isDev bool) (*kubernetes.Clientset, error) {
	config, err := makeK8sConfig(isDev)
	if err != nil {
		return nil, err
	}
	theClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return theClient, nil
}
