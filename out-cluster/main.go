//The Source for this example is from
//https://www.youtube.com/watch?v=vlw1NYySbmQ&list=PLh4KH3LtJvRTb_J-8T--wZeOBV3l3uQhc

package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//Look for the Config file and save it in kubeconfig
	// This allows us to talk to various k8s apis
	kubeconfig := flag.String("kubeconfig", "/home/santosh/.kube/config", "Location to your kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("error %s Could not load the kubeconfig\n", err.Error())

		//For InCluster setup, we need to tell to look into Service Acount configurations
		// for finding the kubeconfig.
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Println("Could not locate In Cluster config %s\n", err.Error())
		}
	}
	//Create a clientset to talt to apis
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {

		fmt.Printf("error %s Could not create clientset\n", err.Error())
	}
	//using client set we can now talk to Corev1 api froup and get the list of pods listed in
	// a specific Namespace and list them.
	ctx := context.Background()
	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("error %s Could not List the Pods from default namespace\n", err.Error())
	}
	fmt.Println("Pds from derault Namespace")
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	// similarly we can get the deployments from the Appsv1 API group and
	//list the pods in different namespaces.

	deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("error %s Could not list the deployments from default namespace\n", err.Error())
	}
	for _, d := range deployments.Items {
		fmt.Printf("%s", d)
	}
}
