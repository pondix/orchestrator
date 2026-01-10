# Backend

Conductor's Go service layer lives here. This service will wrap the preserved Orchestrator engine packages while exposing a new API, event stream, and integrations.

## Planned responsibilities
- Fiber API (`/api/v1`) for clusters, instances, and failover operations.
- WebSocket event stream for topology and operation updates.
- Structured JSON logging (Zap) and Prometheus metrics.
- Operation tracking and audit logging backed by MySQL.
- Hooks and integrations (Slack, PagerDuty, webhooks, email).

## Initial layout (target)
- `cmd/conductor/`: service entrypoint.
- `internal/api/`: Fiber handlers, routing, middleware.
- `internal/events/`: WebSocket broadcast + event serialization.
- `internal/operations/`: operation executor + progress tracking.
- `internal/integrations/`: alerting and webhook clients.
- `internal/storage/`: SQL access (explicit SQL or sqlc).

## Engine dependency
Core topology and failover logic remains in the existing packages:
- `go/discovery`, `go/inst`, `go/logic`, `go/raft`.
