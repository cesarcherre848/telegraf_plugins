package main

import (
	"github.com/influxdata/telegraf/plugins/common/shim"
	"telegraf_plugins/plugins/processors/milesight_processor"
)

func main() {
	s := shim.New()
	s.AddProcessor(&milesight_processor.MilesightProcessor{})
	if err := s.Run(shim.PollIntervalDisabled); err != nil {
		panic(err)
	}
}
