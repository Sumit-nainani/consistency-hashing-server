package kubeclient

import (
	"context"
	"fmt"
	"hashing/hashring"
	websocketserver "hashing/websocket-server"
	"log"

	pb "hashing/hashing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func StartKubeClient() {
	// Load kubeconfig (from ~/.kube/config or in-cluster config if deployed inside k8s)
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	podList, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	// Step 2: Start watching from that ResourceVersion
	watcher, err := clientset.CoreV1().Pods("demo").Watch(context.TODO(), metav1.ListOptions{
		ResourceVersion: podList.ResourceVersion,
	})

	if err != nil {
		log.Fatalf("Failed to watch pods: %v", err)
	}
	fmt.Println("🔁 Watching for pod events...")

	for event := range watcher.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			log.Println("Unexpected type")
			continue
		}

		switch event.Type {
		case watch.Added:
			fmt.Printf("🔴 Pod Added: %s/%s\n", pod.Namespace, pod.Name)

		case watch.Deleted:
			fmt.Printf("🔴 Pod deleted: %s/%s\n", pod.Namespace, pod.Name)
			websocketserver.Broadcast(&pb.WebSocketMetadata{
				Type:   "pod",
				Action: "remove",
				Data: &pb.WebSocketMetadata_NodeMetaData{
					NodeMetaData: &pb.NodeMetaData{
						NodeIp:   hashring.GetRingInstance().NodeNameToNodeMetaData[pod.Name].NodeIP,
						NodeName: pod.Name,
						NodeHash: hashring.GetRingInstance().NodeNameToNodeMetaData[pod.Name].NodeHash,
					},
				},
			})
			hashring.GetRingInstance().RemoveNode(pod.Name)

		case watch.Modified:
			fmt.Printf("✏️ Pod modified: %s/%s %s \n", pod.Namespace, pod.Name, pod.Status.PodIP)

			if pod.Status.Phase == corev1.PodRunning {
				go func() {
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
				}()
			}
		default:
			fmt.Println("Unknown event type")
		}
	}
}
