package termination

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-logr/logr"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	alibabacloudTerminationEndpointURL = "http://100.100.100.200/latest/meta-data/spot/termination-time"
)

// Handler represents a handler that will run to check the termination
// notice endpoint and delete Machine's if the instance termination notice is fulfilled.
type Handler interface {
	Run(stop <-chan struct{}) error
}

// NewHandler constructs a new Handler
func NewHandler(logger logr.Logger, cfg *rest.Config, pollInterval time.Duration, namespace, nodeName string) (Handler, error) {
	machinev1.AddToScheme(scheme.Scheme)
	c, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	pollURL, err := url.Parse(alibabacloudTerminationEndpointURL)
	if err != nil {
		// This should never happen
		panic(err)
	}

	logger = logger.WithValues("node", nodeName, "namespace", namespace)

	return &handler{
		client:       c,
		pollURL:      pollURL,
		pollInterval: pollInterval,
		nodeName:     nodeName,
		namespace:    namespace,
		log:          logger,
	}, nil
}

// handler implements the logic to check the termination endpoint and delete the
// machine associated with the node
type handler struct {
	client       client.Client
	pollURL      *url.URL
	pollInterval time.Duration
	nodeName     string
	namespace    string
	log          logr.Logger
}

// Run starts the handler and runs the termination logic
func (h *handler) Run(stop <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())

	errs := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		errs <- h.run(ctx, wg)
	}()

	select {
	case <-stop:
		cancel()
		// Wait for run to stop
		wg.Wait()
		return nil
	case err := <-errs:
		cancel()
		return err
	}
}

func (h *handler) run(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	machine, err := h.getMachineForNode(ctx)
	if err != nil {
		return fmt.Errorf("error fetching machine for node (%q): %v", h.nodeName, err)
	}

	logger := h.log.WithValues("machine", machine.Name)
	logger.V(1).Info("Monitoring node for machine")

	if err := wait.PollImmediateUntil(h.pollInterval, func() (bool, error) {
		resp, err := http.Get(h.pollURL.String())
		if err != nil {
			return false, fmt.Errorf("could not get URL %q: %v", h.pollURL.String(), err)
		}
		switch resp.StatusCode {
		case http.StatusNotFound:
			// Instance not terminated yet
			logger.V(2).Info("Instance not marked for termination")
			return false, nil
		case http.StatusOK:
			// Instance marked for termination
			return true, nil
		default:
			// Unknown case, return an error
			return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
	}, ctx.Done()); err != nil {
		return fmt.Errorf("error polling termination endpoint: %v", err)
	}

	// Will only get here if the termination endpoint returned 200
	logger.V(1).Info("Instance marked for termination, deleting Machine")
	if err := h.client.Delete(ctx, machine); err != nil {
		return fmt.Errorf("error deleting machine: %v", err)
	}

	return nil
}

// getMachineForNodeName finds the Machine associated with the Node name given
func (h *handler) getMachineForNode(ctx context.Context) (*machinev1.Machine, error) {
	machineList := &machinev1.MachineList{}
	err := h.client.List(ctx, machineList, client.InNamespace(h.namespace))
	if err != nil {
		return nil, fmt.Errorf("error listing machines: %v", err)
	}

	for _, machine := range machineList.Items {
		if machine.Status.NodeRef != nil && machine.Status.NodeRef.Name == h.nodeName {
			return &machine, nil
		}
	}

	return nil, fmt.Errorf("machine not found for node %q", h.nodeName)
}
