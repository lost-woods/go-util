package k8sutil

type KubeConfig struct {
	ApiServer string
	Token     string
	CaCert    string
}

type GenericWatcher struct {
	Name            string
	InformerChannel chan struct{}
}

type JsonPatch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}
