package kubeclient

import (
	"context"
	"fmt"
	"hashing/client"
	ph "hashing/hashing"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func handleAddPod(podName string, podIp string) {
	c := client.GetClient()
	c.AddNodeForRequest(context.Background(), &ph.NodeRequest{Node: podName, Ip: podIp})
}

func handleRemovePod(podName string) {
	c := client.GetClient()
	c.RemoveNodeForRequest(context.Background(), &ph.NodeRequest{Node: podName})
}

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
			handleRemovePod(pod.Name)
		case watch.Modified:
			fmt.Printf("✏️ Pod modified: %s/%s %s \n", pod.Namespace, pod.Name, pod.Status.PodIP)
			var ch chan string = make(chan string)
			if pod.Status.Phase == corev1.PodRunning {
				go func() {
					for x := range ch {
						handleAddPod(pod.Name, x)
					}
				}()
				go func() {
					ch <- pod.Status.PodIP
				}()
			}
		default:
			fmt.Println("Unknown event type")
		}
	}
}
