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

func getClientSet()(*kubernetes.Clientset , error){
	config, err := rest.InClusterConfig()
	if err!=nil{
		homeDir,_ := os.UserHomeDir()
		kubeconfig := filepath.Join(homeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
		return nil, err
			}
	}
	return kubernetes.NewForConfig(config)
}
func ServiceDiscovery() int  {
	fmt.Println("Discovering services")
	clientset, err := getClientSet()
	if err != nil {
		fmt.Println("Could not get clientset:", err)
		return 0
	}

	services, err :=clientset.CoreV1().Services("").List(context.Background() , v1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing services:", err)
		return  0
	}
	for _,svc:= range services.Items{
		fmt.Println("Service:", svc.Name, "in namespace:", svc.Namespace)
	}
	return 2
}