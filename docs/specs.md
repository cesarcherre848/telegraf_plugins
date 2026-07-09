# Telegraf Plugins Specification & Architecture

This document outlines the standard architecture, interface requirements, and development guidelines for custom Telegraf plugins developed in this repository.

---

## 1. Architecture Overview

Telegraf works as a pipeline with four main plugin types:
1. **Inputs**: Collect metrics from various sources (poll or listen).
2. **Processors**: Transform, filter, or decorate metrics in-flight.
3. **Aggregators**: Group metrics (e.g., average, sum) over time.
4. **Outputs**: Send metrics to various destinations.

In this repository, we focus on **Inputs**, **Processors**, and **Outputs**.

### External Plugins via `execd`

To keep development decoupled from the main Telegraf codebase, we prioritize the `execd` plugin mechanism. The `execd` runner spawns the plugin as a long-running subprocess and communicates with it via standard I/O (stdin/stdout).

- **Inputs (execd)**: The plugin writes metrics to stdout in a supported format (usually Influx Line Protocol) at regular intervals or continuously.
- **Processors (execd)**: Telegraf writes metrics to the plugin's stdin, the plugin processes them, and writes the output metrics back to stdout.
- **Outputs (execd)**: Telegraf writes metrics to the plugin's stdin, and the plugin writes/sends them to the target destination.

---

## 2. Directory Structure Conventions

Every plugin must have its own directory containing:
```
plugins/<type>/<plugin_name>/
├── README.md             # Plugin-specific documentation (configuration details and metric formats)
├── telegraf.conf         # Example telegraf configuration block
├── main.go               # Entry point of the plugin (or main.py, main.js, etc.)
└── <plugin_name>_test.go # Unit tests
```

---

## 3. Plugin Templates

Use the templates below when designing and writing new plugins.

### 3.1 Input Plugin (Go/execd) Template
An input plugin gathers data and writes it to stdout using Influx Line Protocol.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Gather metrics
		metricName := "custom_metric"
		tags := "host=myhost,env=prod"
		fields := "value=42.0,status=\"OK\""
		timestamp := time.Now().UnixNano()

		// Write to stdout in Influx Line Protocol format
		fmt.Printf("%s,%s %s %d\n", metricName, tags, fields, timestamp)
	}
}
```

### 3.2 Processor Plugin (Go/execd) Template
A processor plugin reads line protocol from stdin, processes or modifies the metric, and writes it back to stdout.

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Example transformation: Uppercase the entire line, or inject a tag
		processedLine := strings.Replace(line, " ", ",processed=true ", 1)

		fmt.Println(processedLine)
	}
}
```

### 3.3 Output Plugin (Go/execd) Template
An output plugin reads line protocol from stdin and pushes it to a external system or API.

```go
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse metric and send to third-party API or database
		err := sendToDestination(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending metric: %v\n", err)
		}
	}
}

func sendToDestination(metric string) error {
	// Custom implementation here
	return nil
}
```

---

## 4. Configuration Standards

Each plugin directory must have a `telegraf.conf` file specifying how it is configured. Below are examples of how the `execd` runner is configured inside Telegraf:

### Input execd Configuration
```toml
[[inputs.execd]]
  ## Program path and arguments
  command = ["/usr/local/bin/my_custom_input", "--config", "/etc/my_custom_input.json"]

  ## Data format to consume.
  data_format = "influx"

  ## Interval to run command if it runs once and exits, or for parsing line by line
  # interval = "10s"
```

### Processor execd Configuration
```toml
[[processors.execd]]
  ## Program path and arguments
  command = ["/usr/local/bin/my_custom_processor"]
```

### Output execd Configuration
```toml
[[outputs.execd]]
  ## Program path and arguments
  command = ["/usr/local/bin/my_custom_output"]
  
  ## Data format to send to stdin
  data_format = "influx"
```
