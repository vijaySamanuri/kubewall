package config

import (
	"fmt"
	"github.com/charmbracelet/log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionsinformers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeConfigInfo struct {
	Name         string              `json:"name"`
	AbsolutePath string              `json:"absolutePath"`
	FileExists   bool                `json:"fileExists"`
	Clusters     map[string]*Cluster `json:"clusters"`
}

type Cluster struct {
	Name                     string                                       `json:"name"`
	Namespace                string                                       `json:"namespace"`
	AuthInfo                 string                                       `json:"authInfo"`
	Connected                bool                                         `json:"connected"`
	RestConfig               *rest.Config                                 `json:"-"`
	ClientSet                *kubernetes.Clientset                        `json:"-"`
	DynamicClient            *dynamic.DynamicClient                       `json:"-"`
	DiscoveryClient          *discovery.DiscoveryClient                   `json:"-"`
	SharedInformerFactory    informers.SharedInformerFactory              `json:"-"`
	ExtensionInformerFactory apiextensionsinformers.SharedInformerFactory `json:"-"`
	DynamicInformerFactory   dynamicinformer.DynamicSharedInformerFactory `json:"-"`
}

func (c *Cluster) MarkAsConnected() *Cluster {
	c.Connected = true
	return c
}

func LoadInClusterConfig() (KubeConfigInfo, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return KubeConfigInfo{}, err
	}
	kubeConfig, err := loadClientConfig(config)
	if err != nil {
		return KubeConfigInfo{}, err
	}
	kubeConfig.Name = "incluster"
	newConfig := KubeConfigInfo{
		Name:         "incluster",
		AbsolutePath: "",
		FileExists:   true,
		Clusters: map[string]*Cluster{
			"incluster": kubeConfig,
		},
	}
	return newConfig, nil
}

func LoadK8ConfigFromFile(path string) (map[string]*Cluster, error) {
	cmdConfig, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading cmdConfig from file: %v", err)
	}

	// save from nil map
	clusters := make(map[string]*Cluster)

	for key, cluster := range cmdConfig.Contexts {
		rc, err := restConfig(*cmdConfig, key)
		if err != nil {
			// here we will ignore any invalid context and continue for next
			continue
		}

		kubeConfig, err := loadClientConfig(rc)
		if err != nil {
			// here we will ignore any invalid context and continue for next
			continue
		}

		cfg := &Cluster{
			Name:                     key,
			Namespace:                cluster.Namespace,
			AuthInfo:                 cluster.AuthInfo,
			RestConfig:               kubeConfig.RestConfig,
			ClientSet:                kubeConfig.ClientSet,
			DynamicClient:            kubeConfig.DynamicClient,
			DiscoveryClient:          kubeConfig.DiscoveryClient,
			SharedInformerFactory:    kubeConfig.SharedInformerFactory,
			ExtensionInformerFactory: kubeConfig.ExtensionInformerFactory,
			DynamicInformerFactory:   kubeConfig.DynamicInformerFactory,
		}

		clusters[key] = cfg
	}

	return clusters, nil
}

func restConfig(config api.Config, key string) (*rest.Config, error) {
	config.CurrentContext = key
	cc := clientcmd.NewDefaultClientConfig(config, &clientcmd.ConfigOverrides{CurrentContext: key})

	restConfig, err := cc.ClientConfig()
	if err != nil {
		log.Error("err", "err", err)
		return nil, fmt.Errorf("error creating Kubernetes ClientConfig: %w", err)
	}
	restConfig.ContentType = runtime.ContentTypeProtobuf
	restConfig.QPS = float32(K8SQPS)
	restConfig.Burst = K8SBURST
	if restConfig.BearerToken != "" {
		restConfig.Insecure = true
	}
	return restConfig, nil
}

func loadClientConfig(restConfig *rest.Config) (*Cluster, error) {
	if restConfig == nil {
		return nil, fmt.Errorf("restConfig is nil")
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes clientset: %w", err)
	}

	sharedInformerFactory := informers.NewSharedInformerFactory(clientSet, 0)

	clientset, err := apiextensionsclientset.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes clientset: %w", err)
	}
	externalInformer := apiextensionsinformers.NewSharedInformerFactory(clientset, 0)

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes dynamicClient: %w", err)
	}
	dynamicinformer := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, 0)

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes NewDiscoveryClientForConfig: %w", err)
	}

	return &Cluster{
		RestConfig:               restConfig,
		ClientSet:                clientSet,
		DynamicClient:            dynamicClient,
		DiscoveryClient:          discoveryClient,
		SharedInformerFactory:    sharedInformerFactory,
		ExtensionInformerFactory: externalInformer,
		DynamicInformerFactory:   dynamicinformer,
	}, nil
}
