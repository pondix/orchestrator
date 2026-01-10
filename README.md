# Conductor

Conductor is a MySQL high availability and replication management tool, forked from Orchestrator to preserve proven discovery and failover logic while modernizing the UI, configuration, and integrations. See the refactor roadmap in [docs/REFACTOR_PLAN.md](docs/REFACTOR_PLAN.md).

#### Discovery

Conductor actively crawls through your topologies and maps them. It reads basic MySQL info such as replication status and configuration.

It provides with slick visualization of your topologies, including replication problems, even in the face of failures.

#### Refactoring

Conductor understands replication rules. It knows about binlog file:position, GTID, Pseudo GTID, Binlog Servers.

Refactoring replication topologies can be a matter of drag & drop a replica under another master. Moving replicas around is safe: Conductor will reject an illegal refactoring attempt.

Fine-grained control is achieved by various command line options.

#### Recovery

Conductor uses a holistic approach to detect master and intermediate master failures. Based on information gained from the topology itself, it recognizes a variety of failure scenarios.

Configurable, it may choose to perform automated recovery (or allow the user to choose type of manual recovery). Intermediate master recovery achieved internally to Conductor. Master failover supported by pre/post failure hooks.

Recovery process utilizes Conductor's understanding of the topology and of its ability to perform refactoring. It is based on _state_ as opposed to _configuration_: Conductor picks the best recovery method by investigating/evaluating the topology at the time of
recovery itself.

#### The interface

Conductor supports:

- Command line interface (love your debug messages, take control of automated scripting)
- Web API (HTTP GET access)
- Web interface (modernized UI in progress).

#### Additional perks

- Highly available
- Controlled master takeovers
- Manual failovers
- Failover auditing
- Audited operations
- Pseudo-GTID
- Datacenter/physical location awareness
- MySQL-Pool association
- HTTP security/authentication methods
- Coupled with [orchestrator-agent](https://github.com/github/orchestrator-agent), seed new/corrupt instances
- There is also an [orchestrator-mysql](https://groups.google.com/forum/#!forum/orchestrator-mysql) Google groups forum to discuss topics related to Orchestrator
- More...

Read more in [docs/](docs/)

Authored by [Shlomi Noach](https://github.com/shlomi-noach) at [GitHub](http://github.com). Previously at [Booking.com](http://booking.com) and [Outbrain](http://outbrain.com)

#### License

Conductor is free and open sourced under the [Apache 2.0 license](LICENSE).
