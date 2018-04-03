package autoscaler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"k8s.io/client-go/kubernetes"
)

const (
	// AMQP channel to query and polling interval
	queuename = "fib-work"
	// AMQP channel polling interval
	pollint = 5 * time.Second
	// Threshold when to start scaling up/down workers
	upper = 100
	lower = 10
)

// Logger is the interface for passing a logger to AutoScaler
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// AutoScaler implements auto-scaling functionality based on watching a AMQP queue and scaling a Kubernetes deployment. Must be initialized with New()
type AutoScaler struct {
	// allow custom logger
	log Logger
	// Kubernetes namespace, deployment and max workers to scale
	ns  string
	dpl string
	w   int
	// Kubernetes clientset for interacting with the API
	kclient *kubernetes.Clientset
	// AMQP broker channel and queue threshold
	amqp *amqp.Channel
	// cleanup handline
	wg sync.WaitGroup
}

// New creates a ready to use AutoScaler. It requires a valid logger, AMQP broker address (e.g. "amqp://guest:guest@rabbitmq:5672/") and a Kubernetes namespace and deployment.
func New(kclient *kubernetes.Clientset, log Logger, amqp, ns, dpl string, mw int) (*AutoScaler, error) {
	var err error

	as := AutoScaler{
		// set max workers
		w: mw,
	}

	// sanity checks
	if kclient == nil {
		return nil, errors.New("Kubernetes clientset cannot be nil")
	}
	as.kclient = kclient

	if log == nil {
		return nil, errors.New("logger cannot be nil")
	}
	as.log = log

	if amqp == "" {
		return nil, errors.New("broker cannot be empty")
	}

	as.amqp, err = as.dial(amqp)
	if err != nil {
		return nil, fmt.Errorf("could not connect to broker: %v", err)
	}

	if ns == "" {
		return nil, errors.New("namespace cannot be empty")
	}
	as.ns = ns

	if dpl == "" {
		return nil, errors.New("deployment cannot be empty")
	}
	as.dpl = dpl

	return &as, nil

}

// Run starts the autoscaler and blocks until an error orrcurs or stopCh is closed
func (as *AutoScaler) Run(stopCh chan struct{}) error {
	var err error
	as.log.Println("Starting autoscaler...")

	// allows for cancellation of subroutines if stopCh is closed
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := as.work(ctx)

	// cleanup and first error handling - whichever fires first
	select {
	case <-stopCh:
		cancel()
	case err = <-errCh:
	}

	// wait for all internal Goroutines to finish
	as.wg.Wait()
	return err
}

// encapsulates AMQP polling and Kubernetes scaling logic
// starts an anonymous goroutine and returns a read-only buffered (1) error channel to not block sending goroutine
func (as *AutoScaler) work(ctx context.Context) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		as.wg.Add(1)
		defer as.wg.Done()

		// capture last seen queue depth
		// var last int

		for {
			select {
			case <-time.Tick(pollint):
				q, err := as.amqp.QueueInspect(queuename)
				if err != nil {
					errCh <- fmt.Errorf("could not inspect queue %q: %v", queuename, err)
					return
				}
				// as.log.Printf("Messages in queue %q: %d (last: %d)", q.Name, q.Messages, last)
				err = as.scale(q.Messages)
				if err != nil {
					errCh <- fmt.Errorf("could not scale deployment %q: %v", as.dpl, err)
					return
				}
				// last = q.Messages

			case <-ctx.Done():
				as.log.Println("Got cancelled")
				return
			}
		}
	}()

	return errCh

}
