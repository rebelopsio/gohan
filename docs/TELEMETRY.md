# OpenTelemetry Observability

Gohan includes built-in OpenTelemetry support for distributed tracing, metrics, and observability in both CLI and server modes.

## Configuration

Telemetry can be configured via environment variables:

```bash
# Enable telemetry
export TELEMETRY_ENABLED=true

# Service identification
export TELEMETRY_SERVICE_NAME=gohan
export TELEMETRY_SERVICE_VERSION=0.1.0
export TELEMETRY_ENVIRONMENT=production

# OTLP endpoint (defaults to localhost:4318 for HTTP)
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
```

## Running with Jaeger (Local Development)

### Using Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "4318:4318"   # OTLP HTTP receiver
      - "16686:16686" # Jaeger UI
    environment:
      - COLLECTOR_OTLP_ENABLED=true

  gohan:
    build: .
    ports:
      - "8080:8080"
    environment:
      - TELEMETRY_ENABLED=true
      - TELEMETRY_SERVICE_NAME=gohan
      - TELEMETRY_ENVIRONMENT=development
      - OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4318
    depends_on:
      - jaeger
```

Start services:

```bash
docker-compose up -d
```

Access Jaeger UI at: http://localhost:16686

## CLI vs Server Telemetry

Gohan supports telemetry for both CLI commands and the long-running API server.

### Server Mode

When running as a server (`gohan server`), telemetry:
- ✅ Collects continuous metrics every 10 seconds
- ✅ Traces all HTTP requests automatically
- ✅ Keeps connections open to the OTLP endpoint
- ✅ Suitable for long-running observability

### CLI Mode

When running CLI commands (`gohan init`, `gohan install`, etc.), telemetry:
- ✅ Creates single-command spans
- ✅ Records command execution metrics
- ✅ Flushes immediately before exit (ensures data is sent)
- ✅ Optimized for short-lived executions
- ✅ Minimal performance impact

### Enabling Telemetry

#### For Server

```bash
# Enable telemetry for the server
export TELEMETRY_ENABLED=true
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318

gohan server
```

#### For CLI

```bash
# Enable telemetry for CLI commands
export TELEMETRY_ENABLED=true
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318

# All CLI commands will now send telemetry
gohan init
gohan install --config myconfig.yaml
```

### Using Jaeger Directly

```bash
# Start Jaeger all-in-one
docker run -d --name jaeger \
  -p 4318:4318 \
  -p 16686:16686 \
  jaegertracing/all-in-one:latest

# Start Gohan with telemetry enabled
TELEMETRY_ENABLED=true \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318 \
make server
```

## Running with Grafana Stack

### Docker Compose with Grafana, Tempo, and Loki

```yaml
version: '3.8'

services:
  tempo:
    image: grafana/tempo:latest
    command: ["-config.file=/etc/tempo.yaml"]
    volumes:
      - ./tempo.yaml:/etc/tempo.yaml
    ports:
      - "4318:4318"  # OTLP HTTP
      - "3200:3200"  # Tempo HTTP

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin
    volumes:
      - ./grafana-datasources.yaml:/etc/grafana/provisioning/datasources/datasources.yaml

  gohan:
    build: .
    ports:
      - "8080:8080"
    environment:
      - TELEMETRY_ENABLED=true
      - TELEMETRY_SERVICE_NAME=gohan
      - OTEL_EXPORTER_OTLP_ENDPOINT=tempo:4318
    depends_on:
      - tempo
```

Create `tempo.yaml`:

```yaml
server:
  http_listen_port: 3200

distributor:
  receivers:
    otlp:
      protocols:
        http:
        grpc:

storage:
  trace:
    backend: local
    local:
      path: /tmp/tempo/traces
```

Create `grafana-datasources.yaml`:

```yaml
apiVersion: 1

datasources:
  - name: Tempo
    type: tempo
    access: proxy
    url: http://tempo:3200
    isDefault: true
```

Access Grafana at: http://localhost:3000

## What is Collected?

### Traces (Spans)

With telemetry enabled, Gohan automatically traces:

1. **HTTP Requests** (Server Mode)
   - Request method and path
   - Response status code
   - Request duration
   - Request ID

2. **CLI Commands** (CLI Mode)
   - Command name (e.g., "gohan init")
   - Command execution time
   - Success/failure status
   - Error details

3. **Installation Operations**
   - Session creation
   - Package installation steps
   - System snapshot operations
   - Configuration merging

4. **Repository Operations**
   - Session save/load operations
   - List operations

### Metrics

Gohan collects comprehensive metrics across different areas:

#### CLI Metrics

- `cli.command.duration` (histogram) - Duration of CLI command execution in milliseconds
  - Attributes: `command`, `status`
- `cli.command.count` (counter) - Number of CLI commands executed
  - Attributes: `command`, `status`
- `cli.command.errors` (counter) - Number of CLI command errors
  - Attributes: `command`

#### HTTP Server Metrics

- `http.server.request.duration` (histogram) - HTTP request duration in milliseconds
  - Attributes: `http.method`, `http.route`, `http.status_code`
- `http.server.request.count` (counter) - Number of HTTP requests
  - Attributes: `http.method`, `http.route`, `http.status_code`
- `http.server.active_connections` (gauge) - Number of active HTTP connections
- `http.server.request.size` (histogram) - Size of HTTP request bodies in bytes
  - Attributes: `http.method`, `http.route`
- `http.server.response.size` (histogram) - Size of HTTP response bodies in bytes
  - Attributes: `http.method`, `http.route`

#### Installation Metrics

- `installation.duration` (histogram) - Duration of installation operations in milliseconds
  - Attributes: `status`
- `installation.count` (counter) - Number of installations
  - Attributes: `status` (success/failed/cancelled)
- `installation.packages.count` (counter) - Number of packages installed
  - Attributes: `status`
- `installation.size` (histogram) - Total size of installed packages in bytes
  - Attributes: `status`
- `installation.sessions.active` (gauge) - Number of active installation sessions
- `installation.package.duration` (histogram) - Duration of individual package installation
  - Attributes: `package`, `status`
- `installation.configuration.validation` (counter) - Configuration validation results
  - Attributes: `type`, `status` (valid/invalid)

## Example Traces

### Installation Workflow

```
POST /api/installation/start
├─ CreateSession (span)
├─ SaveSession (span)
└─ Return session ID

POST /api/installation/{id}/execute
├─ LoadSession (span)
├─ TakeSnapshot (span)
├─ MergeConfiguration (span)
├─ InstallPackages (span)
│  ├─ AptUpdate (span)
│  ├─ InstallPackage: docker (span)
│  └─ InstallPackage: kubectl (span)
└─ SaveSession (span)
```

## Custom Telemetry (For Developers)

### CLI Command Telemetry

To add telemetry to a CLI command:

```go
import (
    "context"
    "github.com/rebelopsio/gohan/internal/infrastructure/telemetry"
)

func runMyCommand(cmd *cobra.Command, args []string) error {
    // Load telemetry config
    cfg, _ := config.Load()

    // Initialize telemetry for CLI
    provider, metrics, err := telemetry.InitCLITelemetry(telemetry.Config{
        ServiceName:    cfg.Telemetry.ServiceName,
        ServiceVersion: cfg.Telemetry.ServiceVersion,
        Environment:    cfg.Telemetry.Environment,
        OTLPEndpoint:   cfg.Telemetry.OTLPEndpoint,
        Enabled:        cfg.Telemetry.Enabled,
    })
    if err != nil {
        log.Printf("Warning: failed to initialize telemetry: %v", err)
    }

    // Wrap command execution with telemetry
    return telemetry.CLICommand(
        context.Background(),
        provider,
        metrics,
        "gohan mycommand",
        func(ctx context.Context) error {
            // Your command implementation here
            return performWork(ctx)
        },
    )
}
```

### Custom Spans

To add custom spans in your code:

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

func MyFunction(ctx context.Context) error {
    tracer := otel.Tracer("gohan")
    ctx, span := tracer.Start(ctx, "MyFunction")
    defer span.End()

    // Add attributes to the span
    span.SetAttributes(
        attribute.String("operation", "custom"),
        attribute.Int("count", 42),
    )

    // Do work...

    return nil
}
```

### Installation Metrics

To record installation metrics:

```go
import (
    "github.com/rebelopsio/gohan/internal/infrastructure/telemetry"
)

// Create metrics instance
metrics, err := telemetry.NewInstallationMetrics()
if err != nil {
    return err
}

// Track active session
metrics.IncrementActiveSessions(ctx)
defer metrics.DecrementActiveSessions(ctx)

// Record individual package installation
startTime := time.Now()
err := installPackage(ctx, "docker")
metrics.RecordPackageInstallation(
    ctx,
    "docker",
    time.Since(startTime),
    err == nil,
)

// Record complete installation
metrics.RecordInstallation(
    ctx,
    "success",
    totalDuration,
    packageCount,
    totalSize,
)
```

### HTTP Metrics

HTTP metrics are automatically collected via the tracing middleware. To add custom HTTP metrics:

```go
import (
    "github.com/rebelopsio/gohan/internal/infrastructure/telemetry"
)

// Create metrics instance
metrics, err := telemetry.NewHTTPMetrics()
if err != nil {
    return err
}

// Record request
metrics.RecordRequest(
    ctx,
    "POST",
    "/api/installation/start",
    201,
    duration,
    requestSize,
    responseSize,
)
```

## Production Considerations

1. **Sampling**: The default configuration uses `AlwaysSample()`. In production, consider using probability-based sampling:

   ```go
   // In telemetry.go
   sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)) // 10% sampling
   ```

2. **Security**: Use TLS for OTLP endpoint in production:

   ```go
   // Remove WithInsecure() and configure TLS
   otlptracehttp.New(ctx,
       otlptracehttp.WithEndpoint(cfg.OTLPEndpoint),
       otlptracehttp.WithTLSCredentials(credentials),
   )
   ```

3. **Resource Limits**: Configure batch size and timeout appropriately for your traffic:

   ```go
   sdktrace.WithBatcher(exporter,
       sdktrace.WithBatchTimeout(10*time.Second),
       sdktrace.WithMaxExportBatchSize(1024),
   )
   ```

## Troubleshooting

### Traces Not Appearing

1. Check telemetry is enabled:
   ```bash
   echo $TELEMETRY_ENABLED
   ```

2. Verify OTLP endpoint is reachable:
   ```bash
   curl http://localhost:4318/v1/traces
   ```

3. Check Gohan logs for telemetry initialization:
   ```
   Telemetry enabled - exporting to localhost:4318
   ```

### Performance Impact

- Tracing adds minimal overhead (~1-2ms per request)
- Use sampling in high-traffic environments
- Monitor memory usage if batching large numbers of spans

## Visualizing Metrics

### Prometheus + Grafana

For production monitoring, use Prometheus to scrape metrics and Grafana for visualization.

#### Docker Compose with Prometheus

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin
    volumes:
      - ./grafana-datasources.yaml:/etc/grafana/provisioning/datasources/datasources.yaml

  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    ports:
      - "4318:4318"  # OTLP HTTP
      - "8889:8889"  # Prometheus metrics exporter
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    command: ["--config=/etc/otel-collector-config.yaml"]

  gohan:
    build: .
    ports:
      - "8080:8080"
    environment:
      - TELEMETRY_ENABLED=true
      - OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318
    depends_on:
      - otel-collector
```

#### OpenTelemetry Collector Configuration

Create `otel-collector-config.yaml`:

```yaml
receivers:
  otlp:
    protocols:
      http:
      grpc:

processors:
  batch:

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
  otlp/jaeger:
    endpoint: jaeger:4317
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/jaeger]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [prometheus]
```

#### Prometheus Configuration

Create `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'otel-collector'
    static_configs:
      - targets: ['otel-collector:8889']
```

#### Example Grafana Dashboards

**Installation Performance Dashboard**

Queries:
```promql
# Installation duration (p95)
histogram_quantile(0.95,
  rate(installation_duration_bucket[5m])
)

# Installations per minute
rate(installation_count[1m])

# Success rate
sum(rate(installation_count{status="success"}[5m])) /
sum(rate(installation_count[5m]))

# Active sessions
installation_sessions_active
```

**HTTP Performance Dashboard**

Queries:
```promql
# Request duration (p99)
histogram_quantile(0.99,
  rate(http_server_request_duration_bucket[5m])
)

# Requests per second by route
sum by (http_route) (
  rate(http_server_request_count[1m])
)

# Error rate
sum(rate(http_server_request_count{http_status_code=~"5.."}[5m])) /
sum(rate(http_server_request_count[5m]))

# Active connections
http_server_active_connections
```

**CLI Usage Dashboard**

Queries:
```promql
# Command execution time (p90)
histogram_quantile(0.90,
  rate(cli_command_duration_bucket[5m])
)

# Most used commands
topk(10,
  sum by (command) (
    rate(cli_command_count[1h])
  )
)

# Command failure rate
sum by (command) (
  rate(cli_command_errors[5m])
)
```

### Query Examples

**Find slow installations:**
```promql
installation_duration > 300000  # >5 minutes
```

**Find failing packages:**
```promql
sum by (package) (
  rate(installation_package_duration{status="failure"}[10m])
)
```

**HTTP endpoint latency:**
```promql
histogram_quantile(0.95,
  sum by (http_route, le) (
    rate(http_server_request_duration_bucket[5m])
  )
)
```
