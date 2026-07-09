package milesight_processor

import (
	"github.com/influxdata/telegraf"
)

type DeviceMapping struct {
	Name     string            `toml:"name"`
	Mappings map[string]string `toml:"mappings"`
}

type MilesightProcessor struct {
	Devices map[string]DeviceMapping `toml:"devices"`
}

func (p *MilesightProcessor) SampleConfig() string {
	return `
  [devices]
    # [devices.<devEUI>]
    #   name = "Device Alias"
    #   mappings = { "old_field" = "new_field" }
`
}

func (p *MilesightProcessor) Description() string {
	return "Mapea canales Modbus genéricos de dispositivos Milesight a nombres legibles por humanos"
}

func (p *MilesightProcessor) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, m := range in {
		devEUI, has := m.GetTag("devEUI")
		if !has {
			// Intento alternativo en campos
			if val, ok := m.GetField("devEUI"); ok {
				if str, ok := val.(string); ok {
					devEUI = str
					has = true
				}
			}
		}

		if !has {
			continue
		}

		deviceConfig, ok := p.Devices[devEUI]
		if !ok {
			continue
		}

		// Cambiar tag de nombre del dispositivo si se especifica
		if deviceConfig.Name != "" {
			m.AddTag("deviceName", deviceConfig.Name)
		}

		// Renombrar campos dinámicamente
		for oldField, newField := range deviceConfig.Mappings {
			if val, ok := m.GetField(oldField); ok {
				m.AddField(newField, val)
				m.RemoveField(oldField)
			}
		}
	}
	return in
}
