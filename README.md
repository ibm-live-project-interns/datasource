# Datasource Service

Simulates or connects to network devices sending SNMP traps, syslogs, and metadata for the NOC Dashboard.

## Overview

This service generates sample network events, validates them using the shared Event model, and forwards them to the Ingestor Core via HTTP with retry logic.

## Features

- SNMP trap simulation
- Syslog message generation
- Metadata enrichment
- **HTTP-based event forwarding** (with retry logic and exponential backoff)
- **Shared Event model validation** (uses `ingestor/shared` package)
- Health check before sending events

## Quick Start

### Prerequisites

- Go 1.23+
- Ingestor Core running at `http://localhost:8001`

### Run Locally

```bash
# Set environment variables
export INGESTOR_CORE_URL=http://localhost:8001

# Run the service
go run main.go
```

### With Docker

Use docker-compose from the [ui repository](https://github.com/ibm-live-project-interns/ui):

```bash
cd ui
docker compose up -d --build
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `INGESTOR_CORE_URL` | **Yes** | Ingestor Core URL (e.g., `http://localhost:8001`) |

## Data Flow

```
Datasource → Ingestor Core → Event Router → API Gateway → UI
  (HTTP)       :8001          :8082          :8080        :3000
```

## Event Model

Uses unified `Event` model from `ingestor/shared/models`:

| Field | Type | Description |
|-------|------|-------------|
| `event_type` | string | syslog, snmp, metadata |
| `source_host` | string | Device hostname |
| `source_ip` | string | Device IP address |
| `severity` | string | critical, high, medium, low, info |
| `category` | string | Event category |
| `message` | string | Event message |
| `raw_payload` | string | Original data |
| `event_timestamp` | time | When event occurred |

## Mappers

| Mapper | Severity Normalization |
|--------|------------------------|
| `mapper.MapSyslog()` | ERROR→critical, WARN→high |
| `mapper.MapSNMP()` | CRITICAL→critical |
| `mapper.MapMetadata()` | Defaults to info |

## HTTP Client Features

- **Retry logic**: 3 attempts with exponential backoff
- **Health check**: Verify Ingestor Core before sending
- **Validation**: Events validated before sending

## Related Repositories

| Repository | Description |
|------------|-------------|
| [ingestor](https://github.com/ibm-live-project-interns/ingestor) | Backend services |
| [ui](https://github.com/ibm-live-project-interns/ui) | Frontend dashboard |
| [docs](https://github.com/ibm-live-project-interns/docs) | Documentation |
