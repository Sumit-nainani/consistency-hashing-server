package kubeclient

import (
	"context"
	"fmt"
	"hashing/hashring"
	"hashing/utility"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	addNode    = "add"
	removeNode = "remove"
)

var hashringInstance *hashring.HashRing = hashring.GetRingInstance()

func StartKubeClient() {
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	
	for {
		// Getting latest pod list and resourceVersion.
		podList, err := clientset.CoreV1().Pods("demo").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error listing pods: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		watcher, err := clientset.CoreV1().Pods("demo").Watch(context.TODO(), metav1.ListOptions{
			ResourceVersion: podList.ResourceVersion,
		})
		
		if err != nil {
			log.Printf("Failed to start watcher: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("🔁 Watching for pod events...")

    // Go label used for breaking the watcher loop inside switch statement.
	watchLoop:
		for event := range watcher.ResultChan() {
			switch obj := event.Object.(type) {
			// if object type is k8s pod type then We will do upscaling/downscaling.
			case *corev1.Pod:
				pod := obj
                
				switch event.Type {
				// Pod downscaling event , We will broadcast the event only if pod exists.
				case watch.Deleted:
					fmt.Printf("🔴 Pod Deleted: %s/%s\n", pod.Namespace, pod.Name)
					if nodeMeta, exists := hashringInstance.NodeNameToNodeMetaData[pod.Name]; exists {
						utility.BroadcastNodeMetaData(removeNode, nodeMeta.NodeIP, nodeMeta.NodeName, nodeMeta.NodeHash)
					}
					hashringInstance.RemoveNode(pod.Name)

                // Pod modified event , We will broadcast the creation event only if the pod is assigned IP address.
				case watch.Modified:
					fmt.Printf("✏️ Pod Modified: %s/%s %s\n", pod.Namespace, pod.Name, pod.Status.PodIP)
					if pod.Status.Phase == corev1.PodRunning {
						hashringInstance.AddNode(pod.Name, pod.Status.PodIP)
						nodeMeta := hashringInstance.NodeNameToNodeMetaData[pod.Name]
						utility.BroadcastNodeMetaData(addNode, nodeMeta.NodeIP, nodeMeta.NodeName, nodeMeta.NodeHash)
					}

				default:
					fmt.Println("Unknown pod event type")
				}

            // When object type is k8s status, not pod then we will restart the watcher.
			case *metav1.Status:
				log.Printf("📴 Watch closed with status: %s", obj.Message)
				break watchLoop // exit the for loop cleanly

			default:
				log.Printf("❓ Unknown event object type: %T\n", event.Object)
			}
		}

		log.Println("🔁 Restarting watcher...")
		time.Sleep(2 * time.Second)
	}
}
