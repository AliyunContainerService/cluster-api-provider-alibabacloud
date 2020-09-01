package machine

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	timeout = 10 * time.Second
)

var (
	k8sClient     client.Client
	eventRecorder record.EventRecorder
)

func TestMain(m *testing.M) {
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crds")},
	}

	cfg, err := testEnv.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := testEnv.Stop(); err != nil {
			log.Fatal(err)
		}
	}()

	mgr, err := manager.New(cfg, manager.Options{
		Scheme:             scheme.Scheme,
		MetricsBindAddress: "0",
	})
	if err != nil {
		log.Fatal(err)
	}

	doneMgr := make(chan struct{})
	go func() {
		if err := mgr.Start(doneMgr); err != nil {
			log.Fatal(err)
		}
	}()
	defer close(doneMgr)

	k8sClient = mgr.GetClient()
	eventRecorder = mgr.GetEventRecorderFor("alicloud-controller")

	code := m.Run()
	os.Exit(code)
}
