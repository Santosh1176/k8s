//The Source for this example is from
//https://www.youtube.com/watch?v=vlw1NYySbmQ&list=PLh4KH3LtJvRTb_J-8T--wZeOBV3l3uQhc

package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/informers"
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
	// clientset.GroupVersion().Resources(namespace).Verb()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {

		fmt.Printf("error %s Could not create clientset\n", err.Error())
	}
	informers := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	// if err != nil {
	// 	fmt.Println("Getting informers %s\n", err.Error())
	// }
	ch := make(chan struct{})
	c := newController(clientset, informers.Apps().V1().Deployments())
	informers.Start(ch)
	c.run(ch)
	fmt.Println(informers)
}
