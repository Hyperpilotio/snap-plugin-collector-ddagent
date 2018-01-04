package instance

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

	"github.com/hyperpilotio/snap-plugin-collector-ddagent/pkg/dogstatsd"
	"github.com/hyperpilotio/snap-plugin-collector-ddagent/pkg/dogstatsd/message"
	"github.com/hyperpilotio/snap-plugin-collector-ddagent/pkg/dogstatsd/server"
)

type DogStatsd struct {
	server    *server.Server
	metrics   chan *[]plugin.Metric
	done      chan struct{}
	isStarted bool
}

func (d *DogStatsd) Start() (err error) {
	if d.isStarted {
		return errors.New("Server has been started")
	}
	errChan := make(chan error)
	go d.server.Run(errChan)

	d.isStarted = true

	go func() {
		for {
			select {
			case rawData := <-message.Data():
				var statsdMetrics dogstatsd.Metrics
				var snapMetrics []plugin.Metric
				if err := json.Unmarshal(rawData, &statsdMetrics); err == nil {
					parseMetrics(statsdMetrics, &snapMetrics)
					d.metrics <- &snapMetrics
				}
            case err = <-errChan:
				break
			case <-d.done:
				break
			}
		}
	}()

	return nil
}

func parseMetrics(dogstatsdMetrics dogstatsd.Metrics, snapMetrics *[]plugin.Metric) {
	for _, ddMetric := range dogstatsdMetrics.Series {
		tags := map[string]string{}
		tags["host"] = ddMetric.Host
		tags["metric"] = ddMetric.MetricName
		tags["type"] = ddMetric.Type

		ns := plugin.NewNamespace("dogstatsd")
		ns = ns.AddStaticElements(strings.Split(ddMetric.MetricName, ".")...)
		*snapMetrics = append(*snapMetrics, plugin.Metric{
			Namespace: ns,
			Tags:      tags,
			Timestamp: time.Now(),
			Data:      ddMetric.Points[0][1],
		})
	}
}

func (d *DogStatsd) Data() (metrics <-chan *[]plugin.Metric) {
	return d.metrics
}

func (d *DogStatsd) Stop() {
	d.isStarted = false
	close(d.done)
}

// NewDogStatsd create an instance of DogStatsd
func NewDogStatsd() *DogStatsd {
	return &DogStatsd{
		server:    server.NewServer(),
		metrics:   make(chan *[]plugin.Metric, 1000),
		done:      make(chan struct{}),
		isStarted: false,
	}
}
