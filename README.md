# Platform Engineer Lab

A local-first platform engineering lab built with Go, Docker, Terraform,
Ansible, and Kind. It exposes Linux runtime information through a small API
and dashboard, then provides the tooling to run the surrounding environment
without a cloud account.

## Included

- Go API with HTTP Basic authentication and a JSON user store
- Read-only dashboard for runtime information
- Docker image and local Terraform example
- Ansible bootstrap for Docker, Go, Terraform, Kind, `kubectl`, and Helm
- Kind cluster configuration and GitHub Actions validation

## Requirements

- Debian or Ubuntu with `sudo`
- Docker and Go, or use the bootstrap script below
- 2 CPU cores; 4 GiB RAM and 10 GiB disk are recommended before creating a
  Kind cluster

## Run the API

```bash
cd platform-api
go run ./cmd/create-user
go run ./cmd/api
```

Open <http://localhost:8081/> or call the API directly:

```bash
curl http://localhost:8081/health
curl --user <username>:<password> http://localhost:8081/filesystem
```

The user store defaults to `configs/users.json` and is ignored by Git. Set
`PLATFORM_API_USERS_FILE` to an absolute path to use a different store.

### Run with Docker

```bash
cd platform-api
docker build -t platform-api .
docker run --rm -p 8081:8081 \
  --mount type=bind,src="$(pwd)/configs/users.json",dst=/app/config/users.json,readonly \
  platform-api
```

## Local tooling

Bootstrap the local tools:

```bash
bash platform-lab/scripts/install.sh
```

The script installs `ansible-core`; Ansible then installs the remaining
tools. Download URLs, versions, checksums, and the Docker signing-key
fingerprint are maintained in `platform-lab/ansible/group_vars/all.yml`.
Sign out and back in if the script adds your user to the `docker` group.

Create the optional local Kind cluster:

```bash
bash platform-lab/scripts/create-kind-cluster.sh
kubectl get nodes
```

## API

| Endpoint | Description |
|---|---|
| `GET /` | Dashboard |
| `GET /health` | Liveness response |
| `GET /hostname`, `/memory`, `/uptime` | Basic runtime information |
| `GET /network`, `/filesystem`, `/system` | Network, storage, and combined runtime information |

The dashboard is public, but its runtime-data requests require HTTP Basic
authentication. `/health` is also public.

## Validation

```bash
bash scripts/ci.sh
```

GitHub Actions runs the same validation on pushes and pull requests.

## Layout

| Path | Purpose |
|---|---|
| `platform-api/` | Go API, dashboard, and Dockerfile |
| `platform-lab/ansible/` | Local tool bootstrap |
| `platform-lab/kubernetes/kind/` | Kind cluster configuration |
| `platform-lab/terraform/` | Docker network and Nginx example |
| `docs/` | Architecture and authentication notes |

## Safety

The API is Linux-specific and reports information visible to its process or
container. Keep it on a trusted local network; do not expose it publicly
without TLS, authentication, and network controls.
