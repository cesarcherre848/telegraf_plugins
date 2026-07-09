package milesight_processor

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf/metric"
)

func TestMilesightProcessorApply(t *testing.T) {
	// Initialize the processor with mappings config
	p := &MilesightProcessor{
		Devices: map[string]DeviceMapping{
			"24e124454e106951": {
				Name: "Transmisor_Caudal_2",
				Mappings: map[string]string{
					"modbus_chn_1": "estado_bomba",
					"modbus_chn_2": "volumen_acumulado_litros",
				},
			},
		},
	}

	// Create test metric mimicking MQTT consumer parsed output
	tags := map[string]string{
		"devEUI":     "24e124454e106951",
		"deviceName": "TRANSMISOR 2",
		"timeZone":   "America/Lima",
	}
	fields := map[string]interface{}{
		"applicationID": int64(3),
		"modbus_chn_1":  float64(0),
		"modbus_chn_2":  float64(25886.955078125),
	}
	m := metric.New("milesight_mqtt", tags, fields, time.Now())

	// Apply processor
	results := p.Apply(m)

	if len(results) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(results))
	}

	result := results[0]

	// Verify tags were updated
	deviceName, _ := result.GetTag("deviceName")
	if deviceName != "Transmisor_Caudal_2" {
		t.Errorf("Expected deviceName tag to be 'Transmisor_Caudal_2', got '%s'", deviceName)
	}

	// Verify old fields were removed and new fields added
	if _, ok := result.GetField("modbus_chn_1"); ok {
		t.Error("modbus_chn_1 should have been removed")
	}
	if _, ok := result.GetField("modbus_chn_2"); ok {
		t.Error("modbus_chn_2 should have been removed")
	}

	estadoBomba, ok := result.GetField("estado_bomba")
	if !ok || estadoBomba != float64(0) {
		t.Errorf("Expected field 'estado_bomba' to be 0, got %v", estadoBomba)
	}

	volumenAcumulado, ok := result.GetField("volumen_acumulado_litros")
	if !ok || volumenAcumulado != float64(25886.955078125) {
		t.Errorf("Expected field 'volumen_acumulado_litros' to be 25886.955078125, got %v", volumenAcumulado)
	}

	// Verify unrelated fields were preserved
	appID, ok := result.GetField("applicationID")
	if !ok || appID != int64(3) {
		t.Errorf("Expected field 'applicationID' to be preserved as 3, got %v", appID)
	}
}

func TestMilesightProcessorNoMapping(t *testing.T) {
	p := &MilesightProcessor{
		Devices: map[string]DeviceMapping{},
	}

	tags := map[string]string{
		"devEUI":     "unknown_device",
		"deviceName": "TRANSMISOR X",
	}
	fields := map[string]interface{}{
		"modbus_chn_1": float64(10),
	}
	m := metric.New("milesight_mqtt", tags, fields, time.Now())

	results := p.Apply(m)

	if len(results) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(results))
	}

	// Verify it remains unmodified
	result := results[0]
	if _, ok := result.GetField("modbus_chn_1"); !ok {
		t.Error("modbus_chn_1 should have been preserved")
	}
}
