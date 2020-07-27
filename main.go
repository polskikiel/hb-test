package main

import (
	"context"
	"fmt"
	"github.com/polskikiel/hb-test/internal/helm"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/client-go/rest"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := rest.InClusterConfig()
	fatalOnError(err)
	log := logrus.New()
	cli, err := helm.NewClient(cfg, "secrets", log)
	fatalOnError(err)

	log.Info("Loading dir")
	ch, err := loader.LoadDir("/testing")
	fatalOnError(err)

	log.Info("Starting provisioning")
	for i:=0; i < 1000; i++ {
		_, err = cli.Install(ch, map[string]interface{}{}, fmt.Sprintf("%s-%d", "test", i), "default")
		fatalOnError(err)
	}
	log.Info("Starting deprovisioning")
	for i:=0; i < 1000; i++ {
		err = cli.Delete(fmt.Sprintf("%s-%d", "test", i), "default")
		fatalOnError(err)
	}
	log.Info("OK")
}

func fatalOnError(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
	}
}

// cancelOnInterrupt calls cancel func when os.Interrupt or SIGTERM is received
func cancelOnInterrupt(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
		case <-c:
			cancel()
		}
	}()
}
