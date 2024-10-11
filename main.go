package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/util/taints"
)

var (
	kubeconfig	*string
	ok 					bool
)

func main() {
	nodeName := flag.String("node", "", "Node name")
	taintEffect := flag.String("taint", "NoSchedule", "Taint Effect Options: NoSchedule, PreferNoSchedule, NoExecute")
	removeTaint := flag.Bool("remove", false, "Boolean to remove taint from node")

	if home := os.Getenv("HOME"); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	if (*nodeName != "") {
		taint := v1.Taint{
			Key: "dedicated",
			Value: "groupName",
			Effect: v1.TaintEffect(*taintEffect),
		}

		node, err := client.CoreV1().Nodes().Get(context.TODO(), *nodeName, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Node before: %s: %v\n", node.Name, node.Spec.Taints)

		if (!*removeTaint) {
			node, ok, err = taints.AddOrUpdateTaint(node, &taint)
		} else {
			node, ok, err = taints.RemoveTaint(node, &taint)
		}
		fmt.Printf("Action OK: %v\n", ok)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Node after: %s: %v\n", node.Name, node.Spec.Taints)

		if _, err := client.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{}); err != nil {
			panic(err)
		}
		fmt.Println("Updated")
	}

	listNodes(client)
}

func listNodes(client *kubernetes.Clientset) {
	fmt.Println("List nodes:")
	nodes, _ := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	for _, node := range nodes.Items {
		fmt.Printf("%s: %v\n", node.Name, node.Spec.Taints)
	}
}
