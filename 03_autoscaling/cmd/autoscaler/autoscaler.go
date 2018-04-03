package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/embano1/vmworld2017/03_autoscaling/internal/autoscaler"
	"github.com/embano1/vmworld2017/03_autoscaling/internal/kubeclient"
)

// build vars set at compile time
var (
	VERSION   string
	BUILDDATE string
	COMMIT    string
)

func main() {

	// flags
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) Absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	dpl := flag.String("d", "receiver", "Deployment to scale")
	ns := flag.String("n", "pubsub", "Namespace of deployment")
	w := flag.Int("w", 5, "Max workers (Kubernetes Pods) to scale up")
	amqp := flag.String("a", "amqp://guest:guest@rabbitmq:5672/", "Fully qualified address of the AMQP broker")
	incluster := flag.Bool("i", false, "Run autoscaller inside Kubernetes cluster.")
	flag.Parse()

	// set up logging
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// create Kubernetes clientset
	kclient, err := kubeclient.New(*incluster, *kubeconfig)
	if err != nil {
		logger.Fatalf("error creating Kubernetes client: %v", err)
	}

	// create autoscaler
	as, err := autoscaler.New(kclient, logger, *amqp, *ns, *dpl, *w)
	if err != nil {
		logger.Fatalf("error creating autoscaler: %v", err)
	}

	// signal handling
	stopCh := signals(logger)

	// blocks until Run() returns
	err = as.Run(stopCh)
	if err != nil {
		logger.Fatalf("error running autoscaler: %v", err)
	}

	// exit 0
	logger.Println("Shutdown complete")
}

// convinience function to get homeDir on multiple OSes
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// signal handling for graceful shutdown
func signals(logger autoscaler.Logger) chan struct{} {
	stopCh := make(chan struct{})

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
		s := <-sigCh
		logger.Printf("Received signal (%v)", s)
		close(stopCh)
	}()

	return stopCh
}
