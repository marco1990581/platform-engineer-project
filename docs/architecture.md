# Architecture

The Platform API follows a layered architecture.

```
HTTP

↓

Middleware

↓

Authentication

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

### handlers

HTTP handlers.

---

## Future Roadmap

Authentication

- JSON repository
- SQLite
- JWT
- RBAC

Observability

- Audit logging
- Prometheus metrics
- Grafana

Security

- TLS
- Rate limiting
- API Keys
