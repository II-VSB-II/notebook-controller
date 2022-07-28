package main

import (
	"context"
	"flag"
	"fmt"
	_ "os"
	"path/filepath"

	_ "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/client-go/util/retry"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	namespace := "default"

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentRes := schema.GroupVersionResource{Group: "api.sandatasystem.com", Version: "v1", Resource: "notebooks"}

	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "api.sandatasystem.com/v1",
			"kind":       "Notebook",
			"metadata": map[string]interface{}{
				"name": "demo-deployment",
			},
			"spec": map[string]interface{}{
				"replicas": 2,
				"template": map[string]interface{}{},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":  "test",
							"image": "public.ecr.aws/j1r0q0g6/notebooks/notebook-servers/jupyter:v1.5.0",
						},
					},
				},
				"user":    "abc",
				"project": "def",
				"access":  []string{"def", "ghi"},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := client.Resource(deploymentRes).Namespace(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())

	// List Deployments
	/*
		list, err := client.Resource(deploymentRes).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, d := range list.Items {
			replicas, found, err := unstructured.NestedInt64(d.Object, "spec", "replicas")
			if err != nil || !found {
				fmt.Printf("Replicas not found for deployment %s: error=%s", d.GetName(), err)
				continue
			}
			fmt.Printf(" * %s (%d replicas)\n", d.GetName(), replicas)
		}
	*/

}
