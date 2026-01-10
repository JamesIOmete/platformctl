# platformctl (V2)

# platformctl (V2)

> **Standardizes developer onboarding, service scaffolding, and day-to-day ops via a composable CLI.**

**platformctl** is a proof-of-concept internal developer platform (IDP) tool. It demonstrates a scalable, decentralized architecture focusing on developer experience (DX), plugin extensibility, and secure self-service workflows.

---

## Architecture

This project demonstrates a **Hub-and-Spoke** CLI architecture designed for organizational scale:

1.  **Core Binary (`platformctl`)**:
    -   Handles Identity (`auth`) and Onboarding (`init`).
    -   Manages local secrets (`secrets`).
    -   Dispatches subcommands to plugins.

2.  **Decoupled Plugins (`platformctl-sim`)**:
    -   **Discovery Contract**: The core CLI searches the system `$PATH` for executables named `platformctl-<name>`.
    -   **Execution**: Running `platformctl foo args...` invokes `platformctl-foo args...`.
    -   **Independence**: Plugins are standalone binaries. This proves how distinct teams (e.g., "Sim Team") can release updates independently of the Platform Team.

3.  **Mocked Backend**:
    -   State is persisted locally in `~/.config/platformctl/mock-state.json` to simulate real API interactions.

---

## Installation

### Prerequisites
-   Linux / WSL / macOS
-   `go` 1.22+ (for building)

### Quick Install (Go Users)
Zero-sudo installation if you have Go installed:
```bash
go install github.com/JamesIOmete/platformctl/cmd/platformctl@latest
go install github.com/JamesIOmete/platformctl/cmd/platformctl-sim@latest
```

### System Install (Recommended)
This installs binaries to `/usr/local/bin`, ensuring plugins are discoverable on your `$PATH` for all users.

```bash
# Build binaries
go build -o ./bin/platformctl ./cmd/platformctl
go build -o ./bin/platformctl-sim ./cmd/platformctl-sim

# Install (requires sudo)
chmod +x install_system.sh
./install_system.sh
```

### Local Build (Testing)
If you prefer not to use sudo, you must add the `./bin` directory to your path manually:
```bash
export PATH=$PWD/bin:$PATH
```

---

## Usage Guide

### 1. Onboarding
Start here to bootstrap your local environment and identity.
```bash
$ platformctl init
```
*Follow the interactive wizard to set your Principal and Scopes.*

### 2. Verify Environment
Check that dependencies (and plugins) are correctly answering.
```bash
$ platformctl doctor
```

### 3. Fleet Management (Core Feature)
Interact with the simulated robot fleet.
```bash
# List all devices
$ platformctl fleet ls

# Check specific device status
$ platformctl fleet status robot-001

# View logs via simulated stream
$ platformctl fleet logs robot-001

# Open a secure tunnel (Mock SSH)
$ platformctl fleet ssh robot-001
```

### 4. Simulations (Plugin Feature)
This command is provided by the external `platformctl-sim` plugin.
```bash
# Submit a new simulation job
$ platformctl sim run --scenario=warehouse_v2

# List all jobs
$ platformctl sim ls
```

### 5. Secrets Management
Securely store local credentials (mocked Vault integration).
```bash
$ platformctl secrets set api_key "s3cr3t"
$ platformctl secrets get api_key
```

---

## Development

### Running Tests
The project includes an End-to-End (E2E) test suite that sets up an isolated environment:
```bash
./tests/e2e.sh
```

### Project Structure
-   `cmd/platformctl`: Main entry point.
-   `cmd/platformctl-sim`: Source for the simulation plugin.
-   `internal/`: Shared libraries (Auth, Storage, Output).
-   `tests/`: E2E verification scripts.

---

## License
MIT License. See [LICENSE](./LICENSE) for details.
