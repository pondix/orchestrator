# Docker Compose

Quickstart demo assets for running Conductor + a sample MySQL cluster.

## Goals
- `docker-compose up` boots Conductor plus a 3-node MySQL topology.
- Provide a scripted failover simulation to validate visibility and recovery.
- Keep the demo fast and self-contained for local development.

## Planned services
- `mysql-primary`: writable master.
- `mysql-replica-1` / `mysql-replica-2`: replicas.
- `conductor`: backend API + UI.
- `init-replication`: one-off job to configure replication and seed users.

## Usage (target)
```bash
docker-compose up
```

Then simulate failure:
```bash
docker-compose stop mysql-primary
```
