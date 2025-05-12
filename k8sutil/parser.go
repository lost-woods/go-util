package k8sutil

import (
	"encoding/json"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func ParseJsonPatch(jsonPatch string) []JsonPatch {
	var patch []JsonPatch
	err := json.Unmarshal([]byte(jsonPatch), &patch)

	if err != nil {
		log.Fatalf("Error unmarshalling json patch '%s': %v", jsonPatch, err)
	}

	return patch
}

func ParseResourceType(resourceType string) schema.GroupVersionResource {
	group := ""
	version := ""
	resource := ""

	parts := strings.Split(strings.Trim(resourceType, "/"), "/")

	switch len(parts) {

	case 2:
		// core group: /v1/Pod
		version = parts[0]
		resource = parts[1]

	case 3:
		// named group: /apps/v1/Deployment
		group = parts[0]
		version = parts[1]
		resource = parts[2]

	default:
		log.Fatalf("Invalid resource type '%s'", resourceType)
	}

	return schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
}
