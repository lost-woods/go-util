package k8sutil

import (
	"fmt"

	"github.com/lost-woods/go-util/osutil"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	log              = osutil.GetLogger()
	defaultEnvPrefix = "KUBERNETES"
)

func GetCurrentNamespace() string {
	return osutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
}

func CreateKubeconfig(prefix ...string) KubeConfig {
	_prefix := defaultEnvPrefix
	if len(prefix) > 0 {
		_prefix = prefix[0]
	}

	return KubeConfig{
		ApiServer: osutil.GetEnvStr(fmt.Sprintf("%s_APISERVER", _prefix), "https://kubernetes.default.svc"),
		Token:     osutil.GetEnvStr(fmt.Sprintf("%s_TOKEN", _prefix), osutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")),
		CaCert:    osutil.GetEnvStr(fmt.Sprintf("%s_CACERT", _prefix), osutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")),
	}
}

func CreateKubernetesClient(kubeConfig KubeConfig) *dynamic.DynamicClient {
	config := &rest.Config{
		Host:        kubeConfig.ApiServer,
		BearerToken: kubeConfig.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(kubeConfig.CaCert),
		},
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
