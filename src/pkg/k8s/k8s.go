package k8s

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1r "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//"k8s.io/client-go/tools/clientcmd"
)

func connectKubeApi() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("❌ Cannot create config from incluster: " + err.Error() + "\n")
		os.Exit(1)
	}

	// creates the clientset
	//home := os.Getenv("HOME")
	//kubeConfigPath := home + "/.kube/configs/config_pov-sddc"
	// fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	//config, _ := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	//if err != nil {
	// 	fmt.Printf("error getting Kubernetes config: %v\n", err)
	// 	os.Exit(1)
	// }
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("❌ Cannot create clientset: " + err.Error() + "\n")
		os.Exit(1)
	}
	return clientset
}

func isPodRunning(c kubernetes.Interface, podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := c.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1r.GetOptions{})
		if err != nil {
			return false, err
		}

		switch pod.Status.Phase {
		case v1.PodRunning:
			return true, nil
		case v1.PodFailed, v1.PodSucceeded:
			return false, nil
		}
		return false, nil
	}
}

func waitForPodRunning(c kubernetes.Interface, namespace, podName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, isPodRunning(c, podName, namespace))
}

// Returns the list of currently scheduled or running pods in `namespace` with the given selector
func ListPods(c kubernetes.Interface, namespace, selector string) (*v1.PodList, error) {
	listOptions := metav1r.ListOptions{LabelSelector: selector}
	podList, err := c.CoreV1().Pods(namespace).List(context.TODO(), listOptions)

	if err != nil {
		return nil, err
	}
	return podList, nil
}

// Wait up to timeout seconds for all pods in 'namespace' with given 'selector' to enter running state.
// Returns an error if no pods are found or not all discovered pods enter running state.
func WaitForPodBySelectorRunning(namespace, selector string, timeout int) {
	c := connectKubeApi()

	podList, err := ListPods(c, namespace, selector)
	if err != nil {
		fmt.Printf("❌ Cannot list pods: " + err.Error() + "\n")
		os.Exit(1)
	}
	if len(podList.Items) == 0 {
		fmt.Printf("❌ no pods in namespace " + namespace + " matching selector " + selector + "\n")
		os.Exit(1)
	}

	for _, pod := range podList.Items {
		if strings.Contains(pod.Name, "upgrade-agent-cp") {
			continue
		}
		fmt.Printf("Checking: " + pod.Name + "\n")
		if err := waitForPodRunning(c, namespace, pod.Name, time.Duration(timeout)*time.Second); err != nil {
			fmt.Printf("❌ Cannot create clientset: " + err.Error() + "\n")
			os.Exit(1)
		}
	}
	time.Sleep(10 * time.Second)
}
