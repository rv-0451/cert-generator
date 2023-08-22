package kubeapi

import (
	"bytes"
	"context"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var KClient *KubeClient

type KubeClient struct {
	Clientset *kubernetes.Clientset
}

func init() {
	KClient = NewKubeClient()
}

func NewKubeClient() *KubeClient {
	return &KubeClient{
		Clientset: getClientset(),
	}
}

func (kc *KubeClient) InjectCAtoMutatingWebhook(webhookName string, caCert *bytes.Buffer) error {
	webhookConf, err := kc.Clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(
		context.TODO(),
		webhookName,
		metav1.GetOptions{},
	)
	if err != nil {
		return err
	}

	for i := range webhookConf.Webhooks {
		webhookConf.Webhooks[i].ClientConfig.CABundle = caCert.Bytes()
	}

	log.Printf("Updating mutating webhook '%s' with new CA...", webhookName)
	_, err = kc.Clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().Update(
		context.TODO(),
		webhookConf,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return err
	}

	log.Printf("Mutating webhook '%s' was updated successfully.", webhookName)
	return nil
}

func (kc *KubeClient) InjectCAtoValidatingWebhook(webhookName string, caCert *bytes.Buffer) error {
	webhookConf, err := kc.Clientset.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(
		context.TODO(),
		webhookName,
		metav1.GetOptions{},
	)
	if err != nil {
		return err
	}

	for i := range webhookConf.Webhooks {
		webhookConf.Webhooks[i].ClientConfig.CABundle = caCert.Bytes()
	}

	log.Printf("Updating validating webhook '%s' with new CA...", webhookName)
	_, err = kc.Clientset.AdmissionregistrationV1().ValidatingWebhookConfigurations().Update(
		context.TODO(),
		webhookConf,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return err
	}

	log.Printf("Validating webhook '%s' was updated successfully.", webhookName)
	return nil
}

func (kc *KubeClient) CreateSecret(secretName string, secretNamespace string, data map[string][]byte) error {
	_, err := kc.Clientset.CoreV1().Secrets(secretNamespace).Create(
		context.TODO(),
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: secretNamespace,
			},
			Data: data,
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		return err
	}
	return nil
}

func (kc *KubeClient) UpdateSecret(secretName string, secretNamespace string, data map[string][]byte) error {
	_, err := kc.Clientset.CoreV1().Secrets(secretNamespace).Update(
		context.TODO(),
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: secretNamespace,
			},
			Data: data,
		},
		metav1.UpdateOptions{},
	)
	if err != nil {
		return err
	}
	return nil
}

func getClientset() *kubernetes.Clientset {
	log.Println("Building config from in-cluster sa...")
	kubeconfig, err := rest.InClusterConfig()
	if err != nil && err == rest.ErrNotInCluster {
		log.Println("Failed to read in-cluster config. Reading KUBECONIFG env var...")
		kubeconfigPath, ok := os.LookupEnv("KUBECONFIG")
		if !ok {
			log.Panicln("Failed to get KUBECONFIG env var.")
		}
		log.Println("Building config from env var...")
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			log.Panicln("Error getting Kubernetes config.")
		}
	} else if err != nil {
		log.Panicln("Error building config.")
	}

	log.Println("Creating clientset...")
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		log.Panicln(err, "Error getting Kubernetes clientset.")
	}

	return clientset
}
