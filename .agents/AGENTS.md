# Telegraf Plugins Workspace Rules & Steering

This file defines the project-specific rules, design guidelines, and steering instructions. The agent loads these rules automatically at the beginning of each session.

## Context Retention (Automatic Steering)

To maintain context seamlessly between sessions, the agent must adhere to the following workflow:

1. **State Loading**: At the start of any new session or task, the agent must immediately inspect and read the session log file [session_state.md](file:///home/ccherre/Projects/telegraf_plugins/docs/session_state.md). This file contains the active state, historical context, decisions, and current goals.
2. **State Tracking**: Keep track of the active task checklist, current files being modified, and goals.
3. **State Persisting**: At the end of a session or after completing any significant milestone, the agent must update [session_state.md](file:///home/ccherre/Projects/telegraf_plugins/docs/session_state.md) with:
   - Completed items.
   - Newly discovered issues or decisions made.
   - Exact current state and immediate next steps.

---

## Directory & Architecture Standards

All custom plugins must reside under their respective directory in `plugins/`:
- **Input Plugins**: [plugins/inputs/](file:///home/ccherre/Projects/telegraf_plugins/plugins/inputs/)
- **Processor Plugins**: [plugins/processors/](file:///home/ccherre/Projects/telegraf_plugins/plugins/processors/)
- **Output Plugins**: [plugins/outputs/](file:///home/ccherre/Projects/telegraf_plugins/plugins/outputs/)

Each plugin should be self-contained in its own subdirectory containing:
1. `README.md` explaining config parameters, metrics gathered/modified, and usage.
2. The source code (e.g., in Go, Python, or another language using `execd` shim).
3. A sample `telegraf.conf` configuration file snippet.

---

## Technical Guidelines for Telegraf Plugins

- **Execd External Plugin Pattern**: Unless otherwise specified, prioritize building external plugins using the Telegraf `execd` shim pattern (supported in Go, Python, etc.). This allows modular development and execution outside Telegraf's main binary.
- **Go Shims**: If writing Go-based `execd` plugins, leverage Telegraf's official `github.com/influxdata/telegraf/plugins/common/shim` library to facilitate compilation as independent executables.
- **Error Handling**: Implement robust error handling. Plugins must handle missing inputs, network timeouts, and malformed input gracefully without crashing.
- **Performance**: Minimize resource footprint, especially for CPU/Memory, as Telegraf runs as an agent on host systems. Keep polling and parsing overhead low.
