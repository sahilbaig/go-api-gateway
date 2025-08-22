package discovery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Service struct {
	Name      string
	Namespace string
}

func getClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		kubeconfig := filepath.Join(homeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(config)
}

// Returns a map of service name -> pod IPs
func ServiceDiscovery() map[string][]string {
	fmt.Println("Discovering services...")
	clientset, err := getClientSet()
	if err != nil {
		fmt.Println("Could not get clientset:", err)
		return nil
	}

	services, err := clientset.CoreV1().Services("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing services:", err)
		return nil
	}

	servicePods := make(map[string][]string)

	for _, svc := range services.Items {
		selector := svc.Spec.Selector
		if len(selector) == 0 {
			continue
		}

		labelSelector := ""
		for k, v := range selector {
			labelSelector += fmt.Sprintf("%s=%s,", k, v)
		}
		labelSelector = labelSelector[:len(labelSelector)-1]

		pods, err := clientset.CoreV1().Pods(svc.Namespace).List(context.Background(), v1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			fmt.Println("  Error listing pods:", err)
			continue
		}

		for _, pod := range pods.Items {
			servicePods[svc.Name] = append(servicePods[svc.Name], pod.Status.PodIP)
		}
	}

	fmt.Println("\nService -> Pod IPs mapping:")
	for svc, ips := range servicePods {
		fmt.Printf("%s: %v\n", svc, ips)
	}

	return servicePods
}
