# Datasource Service

Simulates network devices generating SNMP traps, syslog events, and device metadata for the NOC Platform. Events are validated using the shared Event model and forwarded to Ingestor Core via HTTP with retry logic.

> **Note:** The datasource does not persist events to a database by default. It generates and forwards events to the ingestor pipeline.

## Architecture

```
                          Datasource Service
    ┌─────────────────────────────────────────────────────────┐
    │                                                         │
    │  main.go (entry point)                                  │
    │    ├── mapper/     → Maps raw JSON to shared Event model│
    │    ├── client/     → HTTP client with retry logic       │
    │    └── config/     → YAML device configuration loader   │
    │                                                         │
    │  Standalone Simulators (cmd/)                           │
    │    ├── cmd/snmp-trap-sim/    → SNMP trap generator      │
    │    ├── cmd/syslog-sim/      → Syslog event generator    │
    │    ├── cmd/metadata-pub/    → Device metadata publisher  │
    │    └── cmd/snmp-trap-listener/ → UDP trap listener      │
    │                                                         │
    │  Simulation Libraries (pkg/)                            │
    │    ├── pkg/snmptrap/    → Trap generation & persistence │
    │    ├── pkg/syslogsim/   → RFC 5424 syslog simulation    │
    │    └── pkg/metadatasim/ → Device metadata generation    │
    │                                                         │
    │  Device Framework (simulator/)                          │
    │    ├── Device interface + Manager                       │
    │    ├── Router  (SNMP trap stub)                         │
    │    └── Switch  (Syslog event stub)                      │
    │                                                         │
    └────────────┬────────────────────────────────────────────┘
                 │ HTTP POST /ingest/event
                 ▼
         Ingestor Core (:8001)
                 │
                 ▼
         Event Router (:8082) → Kafka → AI Core → UI
```

## Data Flow

```
Datasource → Ingestor Core → Event Router → API Gateway → UI
  (HTTP)       :8001          :8082          :8080        :3000
```

## Project Structure

```
datasource/
├── main.go                     # Main entry point — event generation & forwarding
├── Dockerfile                  # Multi-stage Go build (no CGO)
├── docker-compose.yml          # Standalone deployment config
├── go.mod / go.sum             # Go modules (depends on ingestor/shared)
├── .env / .env.example         # Environment configuration
│
├── client/                     # Service clients
│   ├── ingestor_client.go      # HTTP client with retry + health check
│   ├── ingestor_client_test.go # Client tests
│   └── kafka.go                # Kafka producer (build tag: kafka)
│
├── mapper/                     # Event type mappers
│   ├── syslog.go               # Syslog JSON → shared Event
│   ├── snmp.go                 # SNMP JSON → shared Event
│   ├── metadata.go             # Metadata JSON → shared Event
│   ├── resolver.go             # IP resolution with TTL caching
│   ├── syslog_test.go          # Mapper tests
│   ├── snmp_test.go
│   └── metadata_test.go
│
├── config/                     # Simulator configuration
│   ├── config.go               # YAML config loader
│   └── sample.yml              # Reference device configuration
│
├── db/                         # Optional database layer
│   ├── db.go                   # PostgreSQL connection
│   └── event_repo.go           # Event repository (not used in default runtime)
│
├── simulator/                  # Device simulation framework
│   ├── device.go               # Device interface definition
│   ├── manager.go              # Concurrent device manager
│   ├── router.go               # Router simulator (SNMP trap stub)
│   └── switch.go               # Switch simulator (syslog stub)
│
├── pkg/
│   ├── snmptrap/               # SNMP trap simulation library
│   │   ├── generator.go        # Random trap generation
│   │   ├── sender.go           # UDP trap transmission
│   │   ├── store.go            # JSON file persistence
│   │   ├── templates.go        # OID templates (router/switch/firewall)
│   │   └── udplisten.go        # Development UDP listener
│   │
│   ├── syslogsim/              # Syslog simulation library
│   │   ├── generator.go        # RFC 5424 message generation
│   │   └── store.go            # JSON file persistence
│   │
│   └── metadatasim/            # Metadata simulation library
│       └── publisher.go        # Device metadata generation & publishing
│
├── cmd/                        # Standalone entry points
│   ├── snmp-trap-sim/main.go   # SNMP trap simulator CLI
│   ├── syslog-sim/main.go      # Syslog simulator CLI
│   ├── metadata-pub/main.go    # Metadata publisher CLI
│   └── snmp-trap-listener/main.go # UDP trap listener
│
├── data/                       # Sample data files
│   ├── devices-metadata.json   # Sample device inventory
│   ├── snmp-traps.json         # Generated SNMP trap samples
│   └── syslog-events.json      # Generated syslog event samples
│
└── sysylog-listener/           # Syslog UDP listener package
    └── udplisten.go            # StartUDPListener function
```

## Quick Start

### Prerequisites

- Go 1.23+
- Ingestor Core running at `http://localhost:8001`

### Run the Main Service

```bash
# Set required environment variables
export INGESTOR_CORE_URL=http://localhost:8001

# Run with default config
go run main.go

# Run with custom config path
go run main.go config/sample.yml
```

### Run with Docker

```bash
# Standalone (requires prod_org-network to exist)
docker-compose up --build

# Or via the infra orchestrator (recommended)
cd ../infra && python run_local.py
```

## Standalone Simulators

Independent CLI tools for generating specific telemetry types:

### SNMP Trap Simulator

Generates random SNMP v2c traps and sends them over UDP.

```bash
go run cmd/snmp-trap-sim/main.go \
  -addr localhost:5162 \
  -device router \
  -freq 3 \
  -file data/snmp-traps.json
```

| Flag | Default | Description |
|------|---------|-------------|
| `-addr` | `localhost:5162` | UDP destination address |
| `-device` | `router` | Device type: router, switch, firewall |
| `-freq` | `3` | Seconds between traps |
| `-file` | `data/snmp-traps.json` | JSON persistence file |

### Syslog Event Simulator

Generates RFC 5424 syslog messages and sends them over UDP or TCP.

```bash
go run cmd/syslog-sim/main.go \
  -host localhost \
  -port 5140 \
  -protocol udp \
  -batch 5 \
  -batches 3
```

| Flag | Default | Description |
|------|---------|-------------|
| `-host` | `localhost` | Target host |
| `-port` | `5140` | Target port |
| `-protocol` | `udp` | Transport: udp or tcp |
| `-interval` | `2s` | Interval between batches |
| `-batch` | `5` | Messages per batch |
| `-batches` | `3` | Total batches (0 = infinite) |
| `-file` | `data/syslog-events.json` | JSON persistence file |

### Metadata Publisher

Generates simulated device inventory metadata.

```bash
go run cmd/metadata-pub/main.go \
  -devices 10 \
  -output data/devices-metadata.json \
  -updates 5 \
  -update-interval 30s
```

| Flag | Default | Description |
|------|---------|-------------|
| `-output` | `./data/devices-metadata.json` | Output file path |
| `-devices` | `10` | Number of devices to generate |
| `-updates` | `0` | Update cycles (0 = no updates) |
| `-update-interval` | `30s` | Time between metadata updates |

### SNMP Trap Listener

Listens for incoming UDP traps on port 5162 (development tool).

```bash
go run cmd/snmp-trap-listener/main.go
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `INGESTOR_CORE_URL` | Yes | — | Ingestor Core endpoint (e.g. `http://localhost:8001`) |
| `KAFKA_BROKER` | No | `kafka:9092` | Kafka broker (used by optional Kafka client) |
| `IP_RESOLVER_CACHE_TTL_SECONDS` | No | `300` | IP resolver cache TTL |
| `LOG_LEVEL` | No | `info` | Log verbosity |
| `ENV` | No | `dev` | Runtime environment |

## Event Model

Uses the unified `Event` model from `ingestor/shared/models`:

| Field | Type | Description |
|-------|------|-------------|
| `event_type` | string | syslog, snmp, metadata |
| `source_host` | string | Device hostname |
| `source_ip` | string | Resolved device IP address |
| `severity` | string | critical, high, medium, low, info |
| `category` | string | Event category |
| `message` | string | Event message |
| `raw_payload` | string | Original data |
| `event_timestamp` | time | When event occurred |

## Mappers

| Mapper | Input | Severity Normalization |
|--------|-------|------------------------|
| `mapper.MapSyslog()` | Syslog JSON | ERROR → critical, WARN → high |
| `mapper.MapSNMP()` | SNMP JSON | CRITICAL → critical |
| `mapper.MapMetadata()` | Metadata JSON | Defaults to info |

All mappers use `mapper/resolver.go` for IP resolution with TTL-based caching.

## HTTP Client

The `client.IngestorClient` provides:

- Retry logic with 3 attempts and exponential backoff
- Health check endpoint verification before sending
- Event validation using `shared/models.Event.Validate()`

## Kafka Client (Optional)

The `client.KafkaProducer` (in `client/kafka.go`) provides an alternative transport:

- Async event delivery to Kafka topics
- Requires `confluent-kafka-go` (CGO + librdkafka)
- Built only with `go build -tags kafka`
- Not used by default — main.go uses the HTTP client

## Shared Dependencies

This service depends on `ingestor/shared`:

```go
// go.mod
require github.com/ibm-live-project-interns/ingestor/shared v0.0.0
replace github.com/ibm-live-project-interns/ingestor/shared => ../ingestor/shared
```

Packages used from shared:
- `shared/models` — Event struct and validation
- `shared/constants` — EventType and Severity enums
- `shared/config` — Environment variable utilities

## Package Reference

| Package | Description |
|---------|-------------|
| `pkg/snmptrap` | SNMP trap generation, JSON persistence, UDP sender, OID templates |
| `pkg/syslogsim` | RFC 5424 syslog message generation with configurable batches |
| `pkg/metadatasim` | Device inventory metadata generation with periodic updates |
| `simulator` | Device simulation framework with Manager, Router, and Switch stubs |
| `client` | HTTP IngestorClient (default) and Kafka producer (optional) |
| `mapper` | Event type mappers with IP resolution and severity normalization |
| `config` | YAML configuration loader for simulator device definitions |
| `db` | Optional PostgreSQL event repository (not used in default runtime) |

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Test specific package
go test -v ./mapper/...
go test -v ./client/...
```

## Docker

### Dockerfile

Multi-stage build targeting the main entry point:

```dockerfile
FROM golang:1.23-alpine AS builder
# Copies ingestor/shared for cross-module dependency
# Build context must be the parent directory
RUN go build -o datasource .

FROM alpine:latest
CMD ["./datasource"]
```

### docker-compose.yml

Connects to the shared `prod_org-network` for inter-service communication:

```bash
docker-compose up --build
```

Requires the orchestration network from `infra/prod/docker-compose.yml` to be running.

## Contributors

| Author | Contribution |
|--------|-------------|
| **bionicop** (ujjwal) | Core mappers, IP resolver, HTTP client, test suite, Dockerfile |
| **Aishwarya Gilhotra** | SNMP trap generator, syslog simulator, metadata publisher |
| **jamal10101** | Device simulator framework, config loader, README documentation |
| **Myshaa1295** | Simulator runtime validation and startup checks |
| **k-rite** | Kafka producer client, Docker Compose infrastructure |

## Related Repositories

| Repository | Description |
|------------|-------------|
| [ingestor](https://github.com/ibm-live-project-interns/ingestor) | Backend services (API Gateway, Ingestor Core, Event Router) |
| [ai-core](https://github.com/ibm-live-project-interns/ai-core) | IBM Watson AI analysis engine |
| [ui](https://github.com/ibm-live-project-interns/ui) | React frontend dashboard |
| [infra](https://github.com/ibm-live-project-interns/infra) | Infrastructure orchestration and deployment |
