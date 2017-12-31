package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/client"
	"github.com/intelsdi-x/snap/core"

	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	if len(os.Args) != 2 {
		log.Fatal("this program expects a single arg for the host and port example: localhost:1234")
	}
	c, err := client.NewStreamCollectorGrpcClient(
		os.Args[1],
		5*time.Second,
		client.SecurityTLSOff(),
	)
	if err != nil {
		panic(err)
	}

	cfg := cdata.NewNode()
	cfg.AddItem("MaxCollectDuration", ctypes.ConfigValueInt{Value: 5000000000})
	cfg.AddItem("MaxMetricsBuffer", ctypes.ConfigValueInt{Value: 2})
	requested_metrics := []core.Metric{
		plugin.MetricType{
			Namespace_: core.NewNamespace("hyperpilot", "ddagent", "dogstatsd"),
			Config_:    cfg,
		},
	}

	metricsOut, errOut, err := c.StreamMetrics("test-TaskID", requested_metrics)

	if err != nil {
		panic(err)
	}
	go func() {
		for metrics := range metricsOut {
			for _, metric := range metrics {
				log.WithFields(
					log.Fields{
						"metric": metric,
					},
				).Debug("received metric")
			}
		}
	}()

	go func() {
		for err := range errOut {
			log.WithFields(
				log.Fields{
					"error": err,
				},
			).Debugf("error received")
		}
	}()
	<-done
}
