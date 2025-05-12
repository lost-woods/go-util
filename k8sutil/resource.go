package k8sutil

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

func GetResource(client *dynamic.DynamicClient, target schema.GroupVersionResource, namespace string, name string) *unstructured.Unstructured {
	var err error
	var resource *unstructured.Unstructured

	if namespace != "" {
		resource, err = client.Resource(target).Namespace(namespace).Get(
			context.Background(),
			name,
			metav1.GetOptions{},
		)
	} else {
		resource, err = client.Resource(target).Get(
			context.Background(),
			name,
			metav1.GetOptions{},
		)
	}

	if err != nil {
		log.Fatalf("Failed to get %s '%s': %v", target.Resource, name, err)
	}

	return resource
}

func PatchResource(client *dynamic.DynamicClient, target schema.GroupVersionResource, namespace string, name string, patchBytes []byte) {
	var err error

	if namespace != "" {
		_, err = client.Resource(target).Namespace(namespace).Patch(
			context.Background(),
			name,
			types.JSONPatchType,
			patchBytes,
			metav1.PatchOptions{},
		)
	} else {
		_, err = client.Resource(target).Patch(
			context.Background(),
			name,
			types.JSONPatchType,
			patchBytes,
			metav1.PatchOptions{},
		)
	}

	if err != nil {
		log.Fatalf("Failed to patch %s '%s': %v", target.Resource, name, err)
	}
}

func WatchResource(client *dynamic.DynamicClient, target schema.GroupVersionResource, namespace string, name string) cache.SharedIndexInformer {
	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		client,
		0,
		namespace,
		func(options *metav1.ListOptions) {
			options.FieldSelector = fmt.Sprintf("metadata.name=%s", name)
		},
	)

	return informerFactory.ForResource(target).Informer()
}
