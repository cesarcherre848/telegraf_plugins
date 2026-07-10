package milesight_parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
)

type MilesightParser struct{}

func (p *MilesightParser) SampleConfig() string {
	return ""
}

func (p *MilesightParser) Description() string {
	return "Splits generic Modbus channels from Milesight devices into separate metrics with payload timestamp and no tags"
}

func (p *MilesightParser) Apply(in ...telegraf.Metric) []telegraf.Metric {
	var results []telegraf.Metric
	for _, m := range in {
		devEUI, has := m.GetTag("devEUI")
		if !has {
			if val, ok := m.GetField("devEUI"); ok {
				if str, ok := val.(string); ok {
					devEUI = str
					has = true
				}
			}
		}

		if !has {
			// Without a devEUI, we cannot form the measurement name, so drop the metric
			continue
		}

		// Extract timestamp from the payload (seconds since epoch)
		var metricTime time.Time
		if tsVal, ok := m.GetField("timestamp"); ok {
			switch v := tsVal.(type) {
			case int64:
				metricTime = time.Unix(v, 0)
			case float64:
				metricTime = time.Unix(int64(v), 0)
			case int:
				metricTime = time.Unix(int64(v), 0)
			default:
				metricTime = m.Time()
			}
		} else {
			metricTime = m.Time()
		}

		// Find and map each modbus_chn_xxx field
		for _, field := range m.FieldList() {
			if strings.HasPrefix(field.Key, "modbus_chn_") {
				measurementName := fmt.Sprintf("%s_%s", devEUI, field.Key)

				// Omit all tags (empty tags map)
				tags := make(map[string]string)

				fields := map[string]interface{}{
					"value": field.Value,
				}

				newMetric := metric.New(measurementName, tags, fields, metricTime)
				results = append(results, newMetric)
			}
		}
	}
	return results
}
