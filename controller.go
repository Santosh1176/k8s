//The Source of this Custom controller is from Vivek Singhs YouTube video
// https://www.youtube.com/watch?v=lzoWSfvE2yA&list=PLh4KH3LtJvRQ43JAwwjvTnsVOMp0WKnJO

package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslister "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// The Main controller struct to for the controller, This controller usese Informers instead of `watch` verb
// which could be resource hungry on API servers. The Informers maintain a chache of all the objects in memory
// and the relevant fuctions are triggered whernever and Add, Update or Delete calls are made to the Kubernetes objects
// and updates them accordlinglyt from the incache memory.

// This controllers will attach a SVC and Ingress resource, to a Deployment whenever a Deployment is created or Updated
type controller struct {
	clientset         kubernetes.Interface
	deployLister      appslister.DeploymentLister
	deploycacheSynced cache.InformerSynced
	queue             workqueue.RateLimitingInterface
}

// Initializing the controller with relevant values
func newController(clientset kubernetes.Interface, depInformers appsinformer.DeploymentInformer) *controller {
	c := &controller{
		clientset:         clientset,
		deployLister:      depInformers.Lister(),
		deploycacheSynced: depInformers.Informer().HasSynced,
		queue:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekspose"),
	}

	//Add the Handeler funtions to our informers
	depInformers.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			DeleteFunc: handledel,
		},
	)
	return c

}

func (c *controller) run(ch <-chan struct{}) {
	fmt.Println("Starting the Controller")
	// Wait for cache is synced by passing it a chennel of empty struct and InformerSynced
	// Since, local cache needs to be synced beforehand we need to check for it.
	if !cache.WaitForCacheSync(ch, c.deploycacheSynced) {
		fmt.Println("Waiting for cache to sync")
	}
	// If the Cache is synced, Add the objects to a Que through a Goroutine
	go wait.Until(c.worker, 30*time.Second, ch)
	<-ch
}

func (c *controller) worker() {

}
func handleAdd(obj interface{}) {
	fmt.Println("Addfunc was Triggered")
}
func handledel(obj interface{}) {
	fmt.Println("delFunc was triggered")
}
