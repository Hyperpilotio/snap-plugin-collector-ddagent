package main

import (
	"github.com/hyperpilotio/snap-plugin-collector-ddagent/cmd"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	log "github.com/sirupsen/logrus"
)

const (
	pluginName    = "dogstatsd"
	pluginVersion = 1
)

// FIXME rewrite here, add flag
func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Debug("Start plugin")
	plugin.StartStreamCollector(cmd.NewCollector(), pluginName, pluginVersion)
}
