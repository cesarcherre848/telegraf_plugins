# Telegraf Plugins - Active Session State

This file is automatically updated by the AI agent at the end of each session and read at the beginning of each session to retain progress and context.

---

## Current Project Status
- **Phase**: Development & Verification
- **Focus**: Refactored to `milesight_parser` with MQTT payload split logic, verified integration in Docker.

---

## Active Objectives
1. [x] Define steering rules in [.agents/AGENTS.md](file:///home/ccherre/Projects/telegraf_plugins/.agents/AGENTS.md).
2. [x] Create [session_state.md](file:///home/ccherre/Projects/telegraf_plugins/docs/session_state.md) to serve as a session log.
3. [x] Create [specs.md](file:///home/ccherre/Projects/telegraf_plugins/docs/specs.md) to define standard architecture/templates for plugins.
4. [x] Initialize Go module and Telegraf dependency.
5. [x] Implement the `milesight_parser` plugin (code, main, test, conf).
6. [x] Implement the transversal Docker environment (Dockerfile, docker-compose, telegraf.conf).
7. [x] Run unit tests and verify integration.

---

## Key Decisions & Architecture Choices
- **Steering and Rules**: Leveraging Workspace Customization Root `.agents/AGENTS.md` to guide the agent natively and read/write state to `docs/session_state.md` automatically.
- **Telegraf Plugin Model**: Favouring external plugins utilizing Telegraf's `execd` pattern for high modularity and ease of development.
- **Telegraf Go Shim**: Used `github.com/influxdata/telegraf/plugins/common/shim` to abstract Influx Line Protocol parsing and config loading.
- **Transversal Docker Setup**: The Dockerfile automatically compiles all plugins located under `plugins/.../cmd/main.go` and copies them to `/usr/local/bin`, making a single testing image.
- **Docker Compose**: Reduced to only run Telegraf, allowing connection to external brokers.

---

## Completed Milestones
- **2026-07-06**: Initialized the repo structure, steering configs, and plugin directories.
- **2026-07-06**: Fully developed the `milesight_processor` plugin and created the transversal Docker testing environment.
- **2026-07-10**: Refactored the plugin to `milesight_parser`, implemented dynamic payload splitting into separate metrics using the payload's timestamp, omitted all tags from output, and successfully verified both unit tests and Docker integration.

---

## Next Steps / Backlog
1. Plan and develop the next plugins (Inputs or Outputs) as requested by the user.
