# Architecture and operations

The service exposes distinct live and ready checks plus non-sensitive release metadata. CI rejects unformatted/racy Go code and invalid Terraform. The container contains only the compiled binary and runs without root. Kubernetes adds resource bounds, probe semantics, a read-only filesystem, and zero ambient Linux capabilities.

## Release model

CI should build once, attach provenance and an SBOM, scan the image, sign its digest, and promote that same digest through environments. Terraform creates shared cloud resources; a deployment controller applies Kubernetes changes. Secrets come from a managed secret store, not manifests or image layers.

Before production, replace `latest` with an immutable digest, add NetworkPolicy, PodDisruptionBudget, autoscaling, workload identity, ingress policy, TLS, and environment-specific overlays. Those omissions are documented to keep the reference focused.

## Twelve-factor mapping

| Factor | Implementation |
|---|---|
| Codebase | One repository maps to one deployable service |
| Dependencies | Go modules, provider constraints, and container stages are explicit |
| Config | `PORT`, `APP_ENV`, and `RELEASE_SHA` are environment values |
| Backing services | Registry/cloud services are replaceable attached resources |
| Build, release, run | CI, immutable image, and Kubernetes runtime are distinct |
| Processes | Stateless service process; no local durable writes |
| Port binding | Binary exports HTTP on `$PORT` |
| Concurrency | Replica count and autoscaling are platform concerns |
| Disposability | Fast startup plus correct readiness/liveness behavior |
| Dev/prod parity | The same container digest is promoted across environments |
| Logs | Structured logs go to stdout/stderr |
| Admin processes | Migrations and diagnostics run as one-off jobs from the image |
