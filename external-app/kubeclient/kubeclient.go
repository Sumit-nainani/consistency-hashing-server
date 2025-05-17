package kubeclient

import (
	"context"
	"fmt"
	"hashing/hashring"
	websocketserver "hashing/websocket-server"
	"log"
	"time"

	pb "hashing/hashing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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

	var resourceVersion string

	for {
		// Get latest pod list and resourceVersion
		podList, err := clientset.CoreV1().Pods("demo").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error listing pods: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		resourceVersion = podList.ResourceVersion

		watcher, err := clientset.CoreV1().Pods("demo").Watch(context.TODO(), metav1.ListOptions{
			ResourceVersion: resourceVersion,
		})
		if err != nil {
			log.Printf("Failed to start watcher: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("🔁 Watching for pod events...")

	watchLoop:
		for event := range watcher.ResultChan() {
			switch obj := event.Object.(type) {
			case *corev1.Pod:
				pod := obj
				resourceVersion = pod.ResourceVersion

				switch event.Type {
				case watch.Added:
					fmt.Printf("🟢 Pod Added: %s/%s\n", pod.Namespace, pod.Name)

				case watch.Deleted:
					fmt.Printf("🔴 Pod Deleted: %s/%s\n", pod.Namespace, pod.Name)
					if nodeMeta, exists := hashring.GetRingInstance().NodeNameToNodeMetaData[pod.Name]; exists {
						websocketserver.Broadcast(&pb.WebSocketMetadata{
							Type:   "pod",
							Action: "remove",
							Data: &pb.WebSocketMetadata_NodeMetaData{
								NodeMetaData: &pb.NodeMetaData{
									NodeIp:   nodeMeta.NodeIP,
									NodeName: pod.Name,
									NodeHash: nodeMeta.NodeHash,
								},
							},
						})
					}
					hashring.GetRingInstance().RemoveNode(pod.Name)

				case watch.Modified:
					fmt.Printf("✏️ Pod Modified: %s/%s %s\n", pod.Namespace, pod.Name, pod.Status.PodIP)
					if pod.Status.Phase == corev1.PodRunning {

						hashring.GetRingInstance().AddNode(pod.Name, pod.Status.PodIP)
						websocketserver.Broadcast(&pb.WebSocketMetadata{
							Type:   "pod",
							Action: "add",
							Data: &pb.WebSocketMetadata_NodeMetaData{
								NodeMetaData: &pb.NodeMetaData{
									NodeHash: hashring.GetRingInstance().NodeNameToNodeMetaData[pod.Name].NodeHash,
									NodeIp:   pod.Status.PodIP,
									NodeName: pod.Name,
								},
							},
						})

					}
				default:
					fmt.Println("Unknown pod event type")
				}

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
