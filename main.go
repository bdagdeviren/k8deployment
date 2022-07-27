package main

import "C"
import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

type Deployments struct {
	WaitTime    int64        `yaml:"waitTime"`
	Deployments []Deployment `yaml:"deployments"`
}

type Deployment struct {
	Name       string       `yaml:"name"`
	Enable     bool         `yaml:"enable"`
	Path       string       `yaml:"path"`
	Monitoring []Monitoring `yaml:"monitoring"`
	Wait       Wait         `yaml:"wait"`
}

type Monitoring struct {
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}

type Wait struct {
	Name      string
	Namespace string
}

func DecodeYAML(data []byte) (<-chan *unstructured.Unstructured, <-chan error) {

	var (
		chanErr        = make(chan error)
		chanObj        = make(chan *unstructured.Unstructured)
		multidocReader = utilyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(data)))
	)

	go func() {
		defer close(chanErr)
		defer close(chanObj)

		for {
			buf, err := multidocReader.Read()
			if err != nil {
				if err == io.EOF {
					return
				}
				chanErr <- errors.Wrap(err, "failed to read yaml data")
				return
			}

			var typeMeta runtime.TypeMeta
			if err := yaml.Unmarshal(buf, &typeMeta); err != nil {
				continue
			}
			if typeMeta.Kind == "" {
				continue
			}

			obj := &unstructured.Unstructured{
				Object: map[string]interface{}{},
			}

			if err := yaml.Unmarshal(buf, &obj.Object); err != nil {
				chanErr <- errors.Wrap(err, "failed to unmarshal yaml data")
				return
			}

			chanObj <- obj
		}
	}()

	return chanObj, chanErr
}

func PostRequestKubernetes(ctx context.Context, cfg *rest.Config, obj *unstructured.Unstructured) error {

	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return err
	}

	gvk := obj.GroupVersionKind()

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		dr = dyn.Resource(mapping.Resource)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "application/apply-patch",
	})

	return err
}

func ApplyYaml(config *rest.Config, yamlFile string) error {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	chanObj, chanErr := DecodeYAML(data)
	for {
		select {
		case obj := <-chanObj:
			if obj == nil {
				return nil
			}

			if obj.GroupVersionKind().Kind == "Service" {
				println("Prometheus -> " + obj.GetName())
			}

			err := PostRequestKubernetes(context.Background(), config, obj)
			if err != nil {
				return err
			}
		case err := <-chanErr:
			if err == nil {
				return nil
			}
			return errors.Wrap(err, "received error while decoding yaml")
		}
	}

}

func Client() (*rest.Config, *kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return config, clientset, nil
}

func CheckPodStatus(namespace string, podName string) error {
	_, clientset, err := Client()
	if err != nil {
		return err
	}

	isDeployed := false

	pods, err := clientset.CoreV1().Pods(C.GoString(namespace)).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, element := range pods.Items {
		if strings.Contains(element.Name, C.GoString(podName)) {
			if element.Status.Phase == "Running" {
				isDeployed = true
			} else {
				isDeployed = false
			}
		}
	}

	resultGo := ""
	if isDeployed != false {
		resultGo = "Deployment Successfully - Name: " + podName + ", Namespace: " + namespace
	} else {
		resultGo = "Waiting Pod Deployment - Name: " + podName + ", Namespace: " + namespace
	}

	fmt.Println(resultGo)

	return nil
}

func main() {

	//kubeconfig := filepath.Join(
	//	os.Getenv("HOME"), ".kube", "config",
	//)
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//

	data, err := ioutil.ReadFile("deployment.yaml")
	if err != nil {
		panic(err)
	}

	var deployments Deployments

	err = yaml.Unmarshal(data, &deployments)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v", deployments)

	for _, element := range deployments.Deployments {
		fmt.Printf("%+v \n", element)
		if element.Name != "" && element.Path != "" && element.Enable != false {

			//err = ApplyYaml(config, element.Path)
			//if err != nil {
			//	println(err.Error())
			//}
		}
	}

}
