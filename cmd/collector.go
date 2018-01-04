package cmd

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/hyperpilotio/snap-plugin-collector-ddagent/pkg/dogstatsd/instance"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	log "github.com/sirupsen/logrus"
)

func NewCollector() *Collector {
	return &Collector{
		server: instance.NewDogStatsd(),
	}
}

// Collector a collector
type Collector struct {
	metrics []plugin.Metric
	server  *instance.DogStatsd
}

// StreamMetrics takes both an in and out channel of []plugin.Metric
//
// The metrics_in channel is used to set/update the metrics that Snap is
// currently requesting to be collected by the plugin.
//
// The metrics_out channel is used by the plugin to send the collected metrics
// to Snap.
func (c *Collector) StreamMetrics(
	ctx context.Context,
	metrics_in chan []plugin.Metric,
	metrics_out chan []plugin.Metric,
	err chan string) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Debug("No metadata")
	}
	taskID := "no-set"
	if tempVal, ok := md["task-id"]; ok {
		if len(tempVal) == 1 {
			taskID = tempVal[0]
		} else {
			log.Debug("Skipping assignment of metadata")
		}
	}

	dogstatsdServerStarted := false

	for metrics := range metrics_in {
		log.WithFields(
			log.Fields{
				"len(metrics)": len(metrics),
				"task-id":      taskID,
			},
		).Error("Received metrics")
		for _, metric := range metrics {
			log.WithFields(
				log.Fields{
					"metric":  metric.Namespace.String(),
					"task-id": taskID,
				},
			).Debug("Received request")

			if !dogstatsdServerStarted && strings.Contains(metric.Namespace.String(), "dogstatsd") {
				c.server.Start()
				go dispatchMetrics(ctx, c.server.Data(), metrics_out)
			}
		}
	}
	return nil
}

func dispatchMetrics(ctx context.Context, in <-chan *[]plugin.Metric, out chan []plugin.Metric) {
	for {
		select {
		case mt := <-in:
			log.Debugf("Processed metrics: len(*mts) %v", len(*mt))
			out <- *mt
		case <-ctx.Done():
			return
		}
	}
}

/*
	GetMetricTypes returns metric types for testing.
	GetMetricTypes() will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
	available metrics).

	Config info is passed in. This config information would come from global config snap settings.

	The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
*/
func (Collector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}

	vals := []string{"dogstatsd"}
	for _, val := range vals {
		metric := plugin.Metric{
			Namespace: plugin.NewNamespace("hyperpilot", "ddagent", val),
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

/*
	GetConfigPolicy() returns the configPolicy for your plugin.

	A config policy is how users can provide configuration info to
	plugin. Here you define what sorts of config info your plugin
	needs and/or requires.
*/
func (Collector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	return *policy, nil
}
