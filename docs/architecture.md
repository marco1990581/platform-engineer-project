# Architecture

The Platform API follows a layered architecture.

```
HTTP

↓

Middleware

↓

Authentication

↓

Authorization

↓

Repository

↓

System

↓

Linux
```

---

## Goals

- Separation of concerns
- Small reusable packages
- Easy future migration
- Simple testing

---

## Directory Layout

```
cmd/

internal/
    config/

configs/

logs/

docs/
```

---

## Internal Packages

### auth

Authentication subsystem.

### config

Runtime configuration resolution.

### middleware

HTTP middleware.

### models

Domain entities.

### system

Linux information providers.

### authz

Role-based authorization. A static role-to-permission policy and an
`Enforce(role, resource, action)` check, kept separate from `auth`
(authentication answers "who are you"; `authz` answers "what can you do").

### handlers

HTTP handlers.

---

## Future Roadmap

Authentication

- JSON repository
- SQLite
- JWT

Authorization

- ✅ Static role-based policy (`admin`, `viewer`)
- Persisted/manageable policy store
- Fine-grained per-user permissions

Observability

- Audit logging
- Prometheus metrics
- Grafana

Security

- TLS
- Rate limiting
- API Keys
