package main

import (
	"context"
	"flag"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/retry"
	"k8s.io/kubernetes/pkg/util/taints"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	nodeName := flag.String("node", "", "Node name")
	taintEffect := flag.String("taint", "NoSchedule", "Taint Effect Options: NoSchedule, PreferNoSchedule, NoExecute")
	removeTaint := flag.Bool("remove", false, "Boolean to remove taint from node")
	flag.Parse()

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	client, err := ctrlclient.New(cfg, ctrlclient.Options{})
	if err != nil {
		panic(err)
	}

	if (*nodeName != "") {
		taint := corev1.Taint{
			Key: "dedicated",
			Value: "groupName",
			Effect: corev1.TaintEffect(*taintEffect),
		}

		var node corev1.Node
		if err :=client.Get(context.TODO(), ctrlclient.ObjectKey{Name: *nodeName}, &node); err != nil {
			panic(err)
		}

		fmt.Printf("Node before: %s: %v\n", node.Name, node.Spec.Taints)

		var (
			updatedNode *corev1.Node
			updated bool
		)

		if (*removeTaint) {
			updatedNode, updated, _ = taints.RemoveTaint(&node, &taint)
		} else {
			updatedNode, updated, _ = taints.AddOrUpdateTaint(&node, &taint)
		}
		fmt.Printf("Action OK: %v\n", updated)

		fmt.Printf("Node after: %s: %v\n", updatedNode.Name, updatedNode.Spec.Taints)

		err := retry.RetryOnConflict(
			retry.DefaultBackoff, func() error {
				return client.Update(context.TODO(), updatedNode)
			},
		)
		if err != nil {
			panic(err)
		}

		fmt.Println("Updated")
	}
}
