# platformctl

**platformctl** is a proof-of-concept internal CLI written in Go that demonstrates how a
company-wide developer platform tool could provide secure, extensible, self-service workflows
for engineers working with infrastructure, identity, and connected device fleets.

This project is intentionally scoped as a **platform engineering POC**, not a production system.

---

## What This Is

- A **single-entry CLI** for internal workflows
- A **framework** for building internal tools
- A demonstration of **developer experience design**
- A model for **secure-by-default self-service**

---

## What This Is Not

- A real IAM system
- A real infrastructure automation tool
- A robot control system
- A complete internal platform

Mocks and stubs are intentional.

---

## Example Usage

```bash
platformctl help
platformctl doctor
platformctl auth status
platformctl fleet ls
platformctl env bootstrap dev
```

---

## Build and Install

```bash
# From repo root
go build -o ./bin/platformctl ./cmd/platformctl
# Or install into GOPATH/bin
go install ./cmd/platformctl
```

## Configuration (optional)

Config is read from `~/.config/platformctl/config.yaml` (or `.yml` / `.json`).
If missing, commands still run in unauthenticated mode.

Example YAML:
```yaml
principal: jward448
scopes:
  - fleet:read
  - infra:write
```

## Manual Test Script

```bash
go build -o ./bin/platformctl ./cmd/platformctl

# Help and doctor
./bin/platformctl help
./bin/platformctl doctor

# Auth status (with or without config)
./bin/platformctl auth status

# Fleet listing (requires scope fleet:read)
./bin/platformctl fleet ls

# Env bootstrap dry-run (requires scope infra:write)
./bin/platformctl env bootstrap dev

# Plugin demo
echo -e '#!/usr/bin/env bash\necho plugin works' > /tmp/platformctl-hello
chmod +x /tmp/platformctl-hello
PATH="/tmp:$PATH" ./bin/platformctl hello
```

## Automated UA Script

Run the bundled checks end-to-end (expects the binary at ./bin/platformctl):

```bash
go build -o ./bin/platformctl ./cmd/platformctl
./scripts/ua.sh
```

The script asserts help/doctor/auth/fleet/env flows, scope denials, plugin fallback, and reports PASS/FAIL.

## Version

platformctl v1.0.0

## License

MIT License © 2026 James Ward

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## Roadmap / Next

- Add richer doctor checks (network, PATH hints)
- Add configurable output formats (json for scripts)
- Ship install convenience (Makefile targets, possible brew tap)
