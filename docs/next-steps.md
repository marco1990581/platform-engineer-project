# Next Step: Deploy Platform API to Kind

## Starting point

The local platform configuration work is complete in commit `f6830d7`
(`fix(platform): make runtime configuration deployable`).

The commit has not been pushed because GitHub SSH authentication is not
available in this environment.

## Objective

Deploy `platform-api` to the local Kind cluster using the same runtime
configuration contract as the Docker container.

## Scope

1. Add a Kubernetes `Secret` for the local `users.json` user store.
2. Add a `Deployment` that mounts the Secret read-only at
   `/app/config/users.json` and sets `PLATFORM_API_USERS_FILE`.
3. Add liveness and readiness probes using `GET /health`.
4. Add resource requests and limits.
5. Add a `Service` for local access.
6. Document how to build the image, load it into Kind, deploy it, and run
   authenticated smoke tests.

## Done criteria

- The API runs in the `platform-lab` Kind cluster.
- `GET /health` succeeds through the Kubernetes Service.
- A protected endpoint succeeds with valid HTTP Basic credentials.
- The user store is not committed and is mounted only as runtime
  configuration.

## Out of scope

Do not add JWT, RBAC, Prometheus, Grafana, or cloud deployment in this step.
After the Kind deployment works, add focused tests for authentication,
middleware, and system providers.
