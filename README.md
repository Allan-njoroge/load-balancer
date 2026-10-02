# Load Balancer Lab

A small Go-based load balancer laboratory for learning how load balancing works under the hood.

This is intentionally **not** a production load balancer.

The goal is to build the system incrementally, implementing different load-balancing strategies one at a time and observing how their behavior changes under different workloads and server conditions.

---

## Why This Exists

Load balancing is often explained with a simple diagram:

```text
                 ┌──────────────┐
Client ────────► │Load Balancer │
                 └──────┬───────┘
                        │
             ┌──────────┼──────────┐
             ▼          ▼          ▼
          Server A   Server B   Server C
```

That diagram is useful, but it hides most of the interesting engineering.

This project explores what happens inside that box.

Questions this project should eventually answer:

- How does a load balancer choose a server?
- What makes an algorithm static or dynamic?
- What information does a dynamic algorithm need?
- How are active connections tracked?
- What happens when a server dies?
- How does a load balancer detect that failure?
- What happens when a server comes back?
- How do different algorithms behave under the same workload?
- What are the trade-offs between simplicity, fairness, latency, and throughput?
- What concurrency problems appear when multiple requests arrive simultaneously?

---

## Project Goals

### Primary goal

Understand load balancing by implementing it from first principles.

### Secondary goals

- Learn Go through a real systems-oriented project.
- Understand HTTP request forwarding.
- Learn practical Go concurrency.
- Experiment with static and dynamic routing strategies.
- Understand health checks and failure handling.
- Measure algorithm behavior rather than relying only on theory.
- Create a repository that another engineer can inspect and learn from.

---

## Non-Goals

This project is not intended to become:

- Nginx
- HAProxy
- Envoy
- a cloud load-balancing service
- a production-grade distributed infrastructure platform

It does not need:

- Kubernetes
- service meshes
- distributed consensus
- cloud deployment
- complex service discovery
- enterprise configuration systems

Those things may be interesting later, but they are not the point.

# Conceptual Architecture

The system will eventually look roughly like:

```text
                         Client
                           │
                           ▼
                  ┌─────────────────┐
                  │  Load Balancer  │
                  │                 │
                  │  Strategy       │
                  │  Health Checks  │
                  │  Metrics        │
                  └────────┬────────┘
                           │
             ┌─────────────┼─────────────┐
             ▼             ▼             ▼
        ┌─────────┐   ┌─────────┐   ┌─────────┐
        │Server A │   │Server B │   │Server C │
        │ :8081   │   │ :8082   │   │ :8083   │
        └─────────┘   └─────────┘   └─────────┘
```

The strategy should be replaceable:

```text
                Load Balancer
                      │
                      ▼
             LoadBalancingStrategy
                      │
        ┌─────────────┼─────────────┐
        ▼             ▼             ▼
   Round Robin     Random     Least Connections
```

---

# Planned Algorithms

Algorithms will be implemented **one at a time**.

Possible progression:

1. Round Robin
2. Random
3. Weighted Round Robin
4. Least Connections
5. Weighted Least Connections
6. Least Response Time
7. Health-aware routing
8. Adaptive/dynamic strategies
9. Experimental strategies

The exact algorithm list may evolve as the project develops.

# Static vs Dynamic

A major learning objective is understanding the distinction.

## Static

The routing decision primarily depends on predefined information.

Example:

```text
A → weight 5
B → weight 3
C → weight 2
```

The load balancer does not necessarily need to observe the current workload before every decision.

Examples include:

- Round Robin
- Weighted Round Robin
- Random
- Weighted Random

## Dynamic

The routing decision uses information about the current state of the system.

Example:

```text
Server A → 20 connections
Server B →  4 connections
Server C → 11 connections
```

A least-connections strategy may therefore choose B.

Examples include:

- Least Connections
- Weighted Least Connections
- Least Response Time
- health-aware strategies
- resource-aware experimental strategies

# Project Philosophy

The project follows one important rule:

> **Implement less, understand more.**

A small implementation that is fully understood is more valuable than a large implementation assembled from libraries.

Each stage should answer a question.

For example:

### Stage: Round Robin

Question:

> How can a load balancer distribute requests sequentially?

### Stage: Weighted Round Robin

Question:

> How can predefined server capacity influence distribution?

### Stage: Least Connections

Question:

> How can current workload influence routing?

### Stage: Health Checks

Question:

> How can the load balancer know that a server should temporarily receive no traffic?

### Stage: Dynamic Algorithms

Question:

> How should routing decisions react to changing system conditions?

# Experiments

The project should include experiments rather than only unit tests.

For example, send 10,000 requests through:

```text
Round Robin
```

and observe:

```text
Server A: 3334
Server B: 3333
Server C: 3333
```

Then change the configuration:

```text
A weight = 5
B weight = 3
C weight = 2
```

and observe how the distribution changes.

The same workload should be usable across different algorithms whenever practical.

# Failure Experiments

Servers should eventually be independently controllable.

For example:

```text
Server A   🟢
Server B   🟢
Server C   🟢
```

Kill Server B:

```text
Server A   🟢
Server B   🔴
Server C   🟢
```

The load balancer should eventually detect the failure and stop routing traffic to B.

When B returns:

```text
Server A   🟢
Server B   🟢
Server C   🟢
```

the system should demonstrate how B becomes eligible again.

The exact recovery behavior is itself part of the experiment.

# Technology

Initial implementation:

- **Go**
- Go standard library
- HTTP
- local server processes

Third-party dependencies should be kept to a minimum.

The project should deliberately use the standard library wherever it makes sense.

# Running the Project

The exact commands will be established during the initial implementation.

The intended local setup will resemble:

```text
Load Balancer   :8080

Backend A       :8081
Backend B       :8082
Backend C       :8083
```

A request might therefore look like:

```text
Client
  │
  │ GET /hello
  ▼
:8080
  │
  ├────► :8081
  ├────► :8082
  └────► :8083
```

# Testing

Run the complete test suite:

```bash
go test ./...
```

For concurrency-related work:

```bash
go test -race ./...
```

Format the project:

```bash
gofmt -w .
```

Run static analysis when useful:

```bash
go vet ./...
```

# Repository Structure

The structure should evolve with the project rather than being over-designed at the beginning.

A likely eventual structure:

```text
.
├── README.md
├── go.mod
├── cmd/
│   ├── loadbalancer/
│   └── backend/
├── internal/
│   ├── balancer/
│   ├── strategy/
│   ├── health/
│   └── metrics/
├── experiments/
├── docs/
└── tests/
```

This is an eventual direction, not a requirement to create all of these directories immediately.

# What Success Looks Like

The project succeeds if, after completing it, you can explain without hand-waving:

- what happens when an HTTP request reaches the load balancer
- how the backend is selected
- how each algorithm differs
- how concurrent requests affect shared state
- how server health is tracked
- what happens when a backend fails
- how a backend recovers
- why one algorithm behaves differently from another
- how you measured those differences
