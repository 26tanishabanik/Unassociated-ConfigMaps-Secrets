package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/aquasecurity/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func clientSetup() *kubernetes.Clientset {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Printf("Error in new client config: %s\n", err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	return clientset

}

func configMap(clientset kubernetes.Clientset) {
	configMaps, err := clientset.CoreV1().ConfigMaps("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error in getting configMap list %s\n", err)
	}
	t := table.New(os.Stdout)
	t.SetHeaders("Name", "Namespace")
	for _, cm := range configMaps.Items {
		if len(cm.OwnerReferences) > 0 {
		} else {
			t.AddRow(cm.Name, cm.Namespace)
		}
	}
	t.Render()
}

func secrets(clientset kubernetes.Clientset) {
	secrets, err := clientset.CoreV1().Secrets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error in getting secrets list %s\n", err)
	}
	t := table.New(os.Stdout)
	t.SetHeaders("Name", "Namespace")
	for _, sc := range secrets.Items {
		if len(sc.OwnerReferences) > 0 {
		} else {
			t.AddRow(sc.Name, sc.Namespace)
		}
	}
	t.Render()
}

func main() {
	clientset := clientSetup()
	fmt.Println("Unreferenced Config Maps")
	configMap(*clientset)
	fmt.Println("Unreferenced Secrets")
	secrets(*clientset)

}
