package main

import (
	"fmt"
	"github.com/polskikiel/hb-test/internal/helm"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/client-go/rest"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {
	cfg, err := rest.InClusterConfig()
	fatalOnError(err)
	log := logrus.New()
	cli, err := helm.NewClient(cfg, "secrets", log)
	fatalOnError(err)

	log.Info("Loading dir")
	ch, err := loader.LoadDir("testing/")
	fatalOnError(err)

	port := os.Getenv("PORT")
	log.Info("Starting server")
	go func() {
		log.Warn(http.ListenAndServe(":"+port, nil))
	}()
	time.Sleep(time.Second * 10)
	log.Info("Starting provisioning")
	name := os.Getenv("NAME")
	for i := 0; i < 1000; i++ {
		_, err = cli.Install(ch, map[string]interface{}{}, fmt.Sprintf("%s-%d", name, i), "default")
		fatalOnError(err)
	}
	log.Info("Starting deprovisioning")
	for i := 0; i < 1000; i++ {
		err = cli.Delete(fmt.Sprintf("%s-%d", name, i), "default")
		fatalOnError(err)
	}
	log.Info("OK")
	time.Sleep(time.Minute)
}

func fatalOnError(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
