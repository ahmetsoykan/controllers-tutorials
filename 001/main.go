package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ahmetsoykan/controllers-tutorials/001/pkg/runtime"
	"github.com/ahmetsoykan/controllers-tutorials/001/pkg/subscription"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog"
)

var (
	minWatchTimeout = 5 * time.Minute
	timeoutSeconds  = int64(minWatchTimeout.Seconds() * (rand.Float64() + 1.0))
	addr            = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func main() {

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	defaultKubernetesClientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building watcher clientset: %s", err.Error())
	}

	// Context
	context := context.TODO()

	if err := runtime.RunLoop([]subscription.ISubscription{
		&subscription.PodSubscription{
			ClientSet:  defaultKubernetesClientSet,
			Ctx:        context,
			Completion: make(chan bool),
		},
	}); err != nil {
		klog.Fatalf(err.Error())
	}
}
