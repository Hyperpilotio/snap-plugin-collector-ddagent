package main

import (
	"github.com/hyperpilotio/snap-plugin-collector-ddagent/cmd"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	pluginName    = "dogstatsd"
	pluginVersion = 1
)

func main() {
	plugin.StartStreamCollector(cmd.NewCollector(), pluginName, pluginVersion)
}
