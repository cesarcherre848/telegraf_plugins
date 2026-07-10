package main

import (
	"github.com/influxdata/telegraf/plugins/common/shim"
	"telegraf_plugins/plugins/processors/milesight_parser"
)

func main() {
	s := shim.New()
	s.AddProcessor(&milesight_parser.MilesightParser{})
	if err := s.Run(shim.PollIntervalDisabled); err != nil {
		panic(err)
	}
}
