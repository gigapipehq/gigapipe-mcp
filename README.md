<a href="https://gigapipe.com" target="_blank">
  <img src='https://github.com/user-attachments/assets/fc8c7ca9-7a18-403d-b2a6-17899a534d33' style="margin-left:-10px;width:200px;" height=200 >
</a>

# Gigapipe MCP Server

A Machine Control Protocol (MCP) server for Gigapipe to query metrics, logs, and traces

## Features

### Prometheus Integration
- Query metrics using PromQL
- List all available labels
- Get values for specific labels
- Support for time range queries

### Loki Integration
- Query logs using LogQL
- List all available labels
- Get values for specific labels
- Support for time range queries

### Tempo Integration
- Query traces by trace ID
- List all available trace tags
- Get values for a specific trace tag
- Support for JSON trace format

## Installation

```bash
# Clone the repository
git clone https://github.com/lmangani/gigapipe-mcp.git
cd gigapipe-mcp

# Build the server
go build -o gigapipe-mcp
```

## Configuration

The server can be configured using environment variables:

```bash
# Required: Gigapipe server address (default: localhost:3100)
export GIGAPIPE_HOST="your-host:3100"

# Optional: HTTP Basic Authentication
export GIGAPIPE_USERNAME="your-username"
export GIGAPIPE_PASSWORD="your-password"
```

## Usage

### Prometheus Tools

1. Query Metrics
```bash
# Query metrics with PromQL
prometheus_query --query="rate(http_requests_total[5m])" --start="2024-01-01T00:00:00Z" --end="2024-01-01T01:00:00Z" --step="1m"
```

2. List Labels
```bash
# List all available labels
prometheus_labels --start="2024-01-01T00:00:00Z" --end="2024-01-01T01:00:00Z"
```

3. Get Label Values
```bash
# Get values for a specific label
prometheus_label_values --label="instance" --start="2024-01-01T00:00:00Z" --end="2024-01-01T01:00:00Z"
```

### Loki Tools

1. Query Logs
```bash
# Query logs with LogQL
loki_query --query='{job="varlogs"}' --start="2024-01-01T00:00:00Z" --end="2024-01-01T01:00:00Z" --limit="100"
```

2. List Labels
```bash
# List all available labels
loki_labels --start="2024-01-01T00:00:00Z" --end="2024-01-01T01:00:00Z"
```

3. Get Label Values
```bash
# Get values for a specific label
loki_label_values --label="job" --start="2024-01-01T00:00:00Z" --end="2024-01-01T01:00:00Z"
```

### Tempo Tools

1. Query Traces
```bash
# Query a trace by ID
tempo_query --trace_id="1234567890abcdef"
```

2. List Tags
```bash
# List all available trace tags
tempo_tags
```

3. Get Tag Values
```bash
# Get values for a specific trace tag
tempo_tag_values --tag="service.name"
```

## API Endpoints

The server communicates with Gigapipe using the following endpoints:

### Prometheus
- Query Range: `/api/v1/query_range`
- Query: `/api/v1/query`
- Labels: `/api/v1/labels`
- Label Values: `/api/v1/label/:name/values`

### Loki
- Query Range: `/loki/api/v1/query_range`
- Query: `/loki/api/v1/query`
- Labels: `/loki/api/v1/label`
- Label Values: `/loki/api/v1/label/:name/values`

### Tempo
- Query Traces: `/api/traces/:traceId`
- Query Traces (JSON): `/api/traces/:traceId/json`
- Tags: `/api/search/tags`
- Tag Values: `/api/search/tag/{name}/values`

## Development

### Prerequisites
- Go 1.16 or later
- Git

### Building from Source
```bash
# Clone the repository
git clone https://github.com/lmangani/gigapipe-mcp.git
cd gigapipe-mcp

# Install dependencies
go mod download

# Build the server
go build -o gigapipe-mcp

# Run tests
go test ./...
```

## License

MIT License
