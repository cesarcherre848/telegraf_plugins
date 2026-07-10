# Milesight Parser Processor Plugin

The `milesight_parser` is an external Telegraf processor plugin running via the `execd` shim. It splits complex MQTT payloads containing multiple generic Modbus channel fields (`modbus_chn_xxx`) into multiple individual Telegraf metrics.

## Behavior

For every incoming metric (e.g. `milesight_mqtt`):
1. **Extraction of devEUI**: The plugin extracts the `devEUI` tag or field. If not present, the metric is dropped.
2. **Extraction of Timestamp**: The timestamp is extracted from the `timestamp` field in the payload (expected to be in seconds since epoch) and converted to a proper time value.
3. **Splitting Channels**: For every field starting with `modbus_chn_`, the plugin creates a new individual metric:
   - **Measurement name**: `{devEUI}_modbus_chn_xxx` (e.g., `24e124454f508602_modbus_chn_1`)
   - **Fields**: A single `value` field containing the channel's value.
   - **Tags**: All tags are omitted.
   - **Timestamp**: Sourced from the parsed `timestamp` field.
4. **Original metric is discarded**.

## Configuration

No external configuration files are needed for this plugin. It executes dynamically.

### Telegraf Configuration Snippet

```toml
[[processors.execd]]
  command = ["/usr/local/bin/milesight_parser"]
  namepass = ["milesight_mqtt"]
```

## Example Transformation

### Input Metric (Line Protocol)
```
milesight_mqtt,devEUI=24e124454f508602,deviceName=MCA,timeZone=America/Lima,topic=Uplink1 modbus_chn_1=0,modbus_chn_2=0,modbus_chn_3=0,modbus_chn_5=35596,applicationID=1,modbus_chn_4=0,timestamp=1783032609 1783719496705504222
```

### Output Metrics (Line Protocol)
```
24e124454f508602_modbus_chn_1 value=0 1783032609000000000
24e124454f508602_modbus_chn_2 value=0 1783032609000000000
24e124454f508602_modbus_chn_3 value=0 1783032609000000000
24e124454f508602_modbus_chn_4 value=0 1783032609000000000
24e124454f508602_modbus_chn_5 value=35596 1783032609000000000
```
