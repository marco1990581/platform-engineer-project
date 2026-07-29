# Platform Engineer Lab

A local, no-cloud-cost platform engineering project. It combines a Go runtime
API, an embedded dashboard, Terraform, Ansible, Docker, and Kind into a
reproducible learning environment.

## What it demonstrates

- A containerized Go API for Linux runtime information.
- A dependency-free, read-only browser dashboard.
- Infrastructure provisioning with Terraform and the Docker provider.
- Declarative environment bootstrap with Ansible.
- A local Kubernetes cluster with Kind, `kubectl`, and Helm.
- GitHub Actions validation for code, infrastructure, and the Docker image.

## Architecture

```text
Local Linux host
  ├── Ansible bootstrap
  │     └── Docker, Terraform, Go, kubectl, Kind, Helm
  ├── Terraform
  │     └── Docker network and Nginx example
  ├── Kind cluster: platform-lab
  │     └── single control-plane node
  └── platform-api
        ├── JSON API on :8081
        └── read-only dashboard at /
```

## Requirements

- Debian or Ubuntu Linux with `sudo`.
- At least 2 CPU cores.
- 4 GiB RAM and 10 GiB free disk are recommended before creating a Kind
  cluster.

The bootstrap installs tools only; it does not start a Kubernetes cluster.

## Quick start

### 1. Bootstrap the local tools

```bash
bash platform-lab/scripts/install.sh
```

The shell entry point installs `ansible-core`, then Ansible installs Git,
Docker, Terraform, `kubectl`, Kind, Helm, and Go. Tool URLs, pinned versions,
SHA-256 checksums, and the Docker signing-key fingerprint are defined in
`platform-lab/ansible/group_vars/all.yml`.

If the bootstrap adds your user to the `docker` group, sign out and back in
before continuing.

### 2. Run the API

```bash
cd platform-api
go run ./cmd/api
```

Open <http://localhost:8081/> for the dashboard, or call an endpoint:

```bash
curl http://localhost:8081/health
curl http://localhost:8081/filesystem
```

To run the API as a container:

```bash
docker build -t platform-api .
docker run --rm -p 8081:8081 platform-api
```

### 3. Create the local Kubernetes cluster

```bash
bash platform-lab/scripts/create-kind-cluster.sh
kubectl get nodes
```

The cluster has one control-plane node to reduce local resource use. Remove it
when you are finished:

```bash
kind delete cluster --name platform-lab
```

### 4. Try the Terraform example

```bash
terraform -chdir=platform-lab/terraform init
terraform -chdir=platform-lab/terraform apply
```

This creates a local Docker network and Nginx container on port `8080`. Remove
it with:

```bash
terraform -chdir=platform-lab/terraform destroy
```

## API

| Endpoint | Description |
|---|---|
| `GET /` | Read-only runtime dashboard. |
| `GET /health` | API health status. |
| `GET /hostname` | Current hostname. |
| `GET /memory` | Total, free, and available memory. |
| `GET /uptime` | Runtime uptime. |
| `GET /network` | Network interfaces and addresses. |
| `GET /filesystem` | Mounted filesystems and capacity statistics. |
| `GET /system` | Aggregate hostname, memory, uptime, and network data. |

`GET /filesystem` returns capacity values in bytes:

```json
{
  "filesystems": [
    {
      "filesystem_type": "overlay",
      "mount_point": "/",
      "total_bytes": 42949672960,
      "used_bytes": 8589934592,
      "available_bytes": 34359738368,
      "usage_percent": 20
    }
  ]
}
```

## Scope and safety

This project is intentionally local-first: Terraform uses Docker, and
Kubernetes runs in Kind rather than a paid cloud provider. The API is
Linux-oriented and reports only information visible to its own process or
container. It does not inspect host services.

Do not expose the API or dashboard publicly without authentication and an
appropriate network policy, because they reveal runtime metadata.

## Continuous integration

Run the full local validation suite with:

```bash
bash scripts/ci.sh
```

GitHub Actions runs the same command for every push and pull request. It checks
Go formatting and builds, Bash and Ansible syntax, Terraform formatting and
validation, and the API Docker image build. Deployment is deferred until the
API has Kubernetes manifests and a target environment.

## Repository layout

| Path | Purpose |
|---|---|
| `platform-api/` | Go API, embedded dashboard, and Dockerfile. |
| `platform-lab/ansible/` | Local Ansible bootstrap playbook and pinned tool metadata. |
| `platform-lab/kubernetes/kind/` | Single-node Kind cluster definition. |
| `platform-lab/terraform/` | Local Docker network and Nginx example. |
| `scripts/ci.sh` | Local validation command used by GitHub Actions. |

## Next steps

- Add handler and provider tests.
- Deploy `platform-api` to the local Kind cluster.
- Add a dashboard view for local tool and cluster readiness.
