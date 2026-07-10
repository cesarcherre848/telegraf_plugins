package milesight_parser

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf/metric"
)

func TestMilesightParserApply(t *testing.T) {
	p := &MilesightParser{}

	// Sample payload representation
	tags := map[string]string{
		"devEUI":     "24e124454f508602",
		"deviceName": "MCA",
		"timeZone":   "America/Lima",
		"topic":      "Uplink1",
	}
	fields := map[string]interface{}{
		"modbus_chn_1":  float64(0),
		"modbus_chn_2":  float64(0),
		"modbus_chn_3":  float64(0),
		"modbus_chn_4":  float64(0),
		"modbus_chn_5":  float64(35596),
		"applicationID": int64(1),
		"timestamp":     int64(1783032609),
	}
	// The original metric timestamp in the line protocol is different
	originalTime := time.Unix(1783719496, 705504222)
	m := metric.New("milesight_mqtt", tags, fields, originalTime)

	results := p.Apply(m)

	// We expect 5 separate metrics (for modbus_chn_1 to modbus_chn_5)
	if len(results) != 5 {
		t.Fatalf("Expected 5 metrics, got %d", len(results))
	}

	// The expected timestamp is the one from the payload (1783032609)
	expectedTime := time.Unix(1783032609, 0)

	expectedMeasurements := map[string]float64{
		"24e124454f508602_modbus_chn_1": 0,
		"24e124454f508602_modbus_chn_2": 0,
		"24e124454f508602_modbus_chn_3": 0,
		"24e124454f508602_modbus_chn_4": 0,
		"24e124454f508602_modbus_chn_5": 35596,
	}

	for _, result := range results {
		meas := result.Name()
		expectedValue, ok := expectedMeasurements[meas]
		if !ok {
			t.Fatalf("Unexpected metric name: %s", meas)
		}

		// Verify fields
		val, ok := result.GetField("value")
		if !ok {
			t.Errorf("Metric %s is missing the 'value' field", meas)
		} else if val.(float64) != expectedValue {
			t.Errorf("Metric %s value expected %f, got %v", meas, expectedValue, val)
		}

		// Verify tags are omitted
		if len(result.TagList()) != 0 {
			t.Errorf("Metric %s has tags, but all tags should be omitted. Got: %v", meas, result.TagList())
		}

		// Verify timestamp is the one from the payload field
		if !result.Time().Equal(expectedTime) {
			t.Errorf("Metric %s time expected %v (from payload), got %v", meas, expectedTime, result.Time())
		}
	}
}
