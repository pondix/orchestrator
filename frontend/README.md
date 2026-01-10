# Frontend

Conductor's web UI lives here (React 19 + TypeScript + shadcn/ui).

## Planned screens
- **Dashboard**: cluster cards, status, and recent events.
- **Topology**: real-time tree/canvas view with instance health.
- **Instance detail**: active operations + history + logs.
- **Failover timeline**: step-by-step view with durations.
- **Config**: tab-based config editor with JSON fallback.

## Real-time updates
The UI subscribes to the backend WebSocket event stream for:
- topology updates
- operation progress + output
- audit and alert events

## Initial layout (target)
- `src/app/`: routes and top-level layout.
- `src/components/`: shadcn/ui components and shared UI.
- `src/features/`: feature modules (topology, operations, config).
- `src/lib/`: API clients, WebSocket helpers, and utilities.
