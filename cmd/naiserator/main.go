package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nais/naiserator"
	"github.com/nais/naiserator/pkg/apis/naiserator/v1alpha1"
	clientV1Alpha1 "github.com/nais/naiserator/pkg/client/clientset/versioned"
	clientset "github.com/nais/naiserator/pkg/client/clientset/versioned"
	informers "github.com/nais/naiserator/pkg/client/informers/externalversions"
	"github.com/nais/naiserator/pkg/metrics"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
	bindAddr   string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.StringVar(&bindAddr, "bind-address", ":8080", "ip:port where http requests are served")
	flag.Parse()
}

func main() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	log.Info("Naiserator starting up")

	// register custom types
	err := v1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatal("unable to add scheme")
	}

	stopCh := StopCh()

	kubeconfig, err := getK8sConfig()
	if err != nil {
		log.Fatal("unable to initialize kubernetes config")
	}

	// serve metrics
	go metrics.Serve(
		bindAddr,
		"/metrics",
		"/ready",
		"/alive",
	)

	applicationInformerFactory := createApplicationInformerFactory(kubeconfig)
	n := naiserator.NewNaiserator(
		createGenericClientset(kubeconfig),
		createApplicationClientset(kubeconfig),
		applicationInformerFactory.Naiserator().V1alpha1().Applications())

	applicationInformerFactory.Start(stopCh)
	n.Run(stopCh)
	<-stopCh

	log.Info("Naiserator has shut down")
}

func createApplicationInformerFactory(kubeconfig *rest.Config) informers.SharedInformerFactory {
	config, err := clientset.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatal("unable to create naiserator clientset")
	}
	return informers.NewSharedInformerFactory(config, time.Second*30)
}

func createApplicationClientset(kubeconfig *rest.Config) *clientV1Alpha1.Clientset {
	clientSet, err := clientV1Alpha1.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatalf("unable to create new clientset")
	}

	return clientSet
}

func createGenericClientset(kubeconfig *rest.Config) *kubernetes.Clientset {
	cs, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return cs
}

func getK8sConfig() (*rest.Config, error) {
	if kubeconfig == "" {
		log.Infof("using in-cluster configuration")
		return rest.InClusterConfig()
	} else {
		log.Infof("using configuration from '%s'", kubeconfig)
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
}

func StopCh() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
