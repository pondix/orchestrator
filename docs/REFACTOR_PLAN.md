# Conductor Refactor Plan

## Scope
This plan preserves Orchestrator’s core HA engine (discovery/topology/failover logic) while replacing surface layers (HTTP/UI/config/logging) and adding Conductor features.

## Package Map (current Orchestrator)
**Core engine to preserve**
- `go/discovery`: topology discovery, health checks, instance probing.
- `go/inst`: instance model, replication state, GTID/Pseudo-GTID helpers.
- `go/logic`: failover orchestration, candidate selection, refactoring logic.
- `go/raft`: Raft leader election and consensus primitives.
- `go/db`, `go/kv`, `go/metrics`: storage and metrics support used by engine.

**Surface layers to replace/modernize**
- `go/http`: legacy HTTP API + web UI endpoints.
- `resources/public` + legacy static UI assets.
- `go/config`: JSON config loader (to be wrapped with new tabbed config model + import).
- `go/cmd/orchestrator`: CLI entrypoint (to evolve into Conductor service).

## Boundary Decision
- Introduce a new **Conductor service layer** that depends on the preserved engine packages above.
- Keep existing engine packages in place initially; refactor by adding a new API/UX on top.

## Migration Checklist (Phases)
1. **Phase 0**: Publish this plan and map packages (done).
2. **Phase 1**: Rename product references and introduce new top-level layout: `backend/`, `frontend/`, `docs/`, `docker-compose/`, `kubernetes/`, `scripts/`.
3. **Phase 2**: Create Fiber API shell + WebSocket event stream; wire read-only endpoints to engine state.
4. **Phase 3**: Add operation tracking schema + executor wrapper with streaming events and audit logs.
5. **Phase 4**: Build React 19 UI (shadcn/ui) with real-time updates; replace legacy UI.
6. **Phase 5**: Implement tabbed config model + JSON fallback + import path.
7. **Phase 6**: Add hooks and integrations (Slack, PagerDuty, webhooks, email).
8. **Phase 7**: Add provisioning engine via SSH with streamed progress.
9. **Phase 8**: Add docker-compose quickstart demo.
10. **Phase 9**: Add Kubernetes Helm chart + app RBAC + OIDC.

## Risks & Mitigations
- **Risk**: Accidental behavior changes in failover logic.
  - **Mitigation**: Keep `go/discovery`, `go/inst`, `go/logic` intact; wrap calls with new service layer.
- **Risk**: Breaking existing configs.
  - **Mitigation**: Build an import/translation path; keep JSON fallback.
- **Risk**: Increased operational surface area.
  - **Mitigation**: Add structured logging, audit trails, and operation tracking early.

## Notes
- Avoid ORMs; use explicit SQL or sqlc for new Conductor DB schemas.
- Ensure all new external calls (SSH/hooks/webhooks) have timeouts and retries.
