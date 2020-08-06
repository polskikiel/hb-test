package main

import (
	"fmt"
	"github.com/polskikiel/hb-test/internal/helm"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/client-go/rest"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"
)

type Config struct {
	Port      string `default:"8080"`
	Name      string `default:"helm-test"`
	Namespace string `default:"default"`
}

func main() {
	clusterCfg, err := rest.InClusterConfig()
	fatalOnError(err)
	log := logrus.New()

	cfg := &Config{}
	err = envconfig.Init(cfg)
	fatalOnError(err)

	log.Info("Loading dir")
	ch, err := loader.LoadDir("testing/")
	fatalOnError(err)

	log.Info("Starting server")
	go func() {
		log.Warn(http.ListenAndServe(":"+cfg.Port, nil))
	}()
	log.Info("Starting provisioning")


	for i := 0; i < 4000; i++ {
		n := strings.ToLower(fmt.Sprintf("%s-%d", cfg.Name, i))

		_, err = helm.NewClient(clusterCfg, "secrets", log).Install(ch, map[string]interface{}{}, n, cfg.Namespace)
		fatalOnError(err)

		err = helm.NewClient(clusterCfg, "secrets", log).Delete(n, cfg.Namespace)
		fatalOnError(err)
	}

	log.Info("OK")
	time.Sleep(time.Hour)
}

func fatalOnError(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
