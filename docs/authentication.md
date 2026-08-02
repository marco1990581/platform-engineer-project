# Authentication Subsystem

## Overview

The Platform API authentication subsystem has been designed around
separation of responsibilities.

The objective is to provide a secure authentication mechanism that
can evolve over time without changing the HTTP layer.

Current authentication method:

- HTTP Basic Authentication
- bcrypt password hashing
- Repository abstraction
- Middleware-based authentication

---

## Architecture

```
HTTP Request
      │
      ▼
Authentication Middleware
      │
      ▼
Basic Authenticator
      │
      ▼
User Repository
      │
      ▼
JSON Storage (future SQLite)
```

---

## Components

### User

Represents an authenticated user.

Passwords are never stored in plaintext.

Only bcrypt hashes are persisted.

---

## Authentication diagnostics

The authentication backend logs the outcome of these steps:

- loading and parsing the configured user file
- looking up the requested username
- comparing the submitted password with bcrypt

Logs include the configured file path and username where needed for
troubleshooting. Passwords and bcrypt hashes are never logged.

---

### Repository

Responsible for loading users.

The authenticator never knows where users are stored.

Current implementation:

- Memory Repository
- File Repository

Future implementations:

- SQLite
- LDAP
- Keycloak

---

### Authenticator

Responsible only for validating credentials.

Responsibilities:

- Validate password
- Compare bcrypt hashes
- Return authenticated user

---

### Middleware

Responsible for protecting HTTP endpoints.

If authentication succeeds:

- User is attached to request context.

Otherwise:

- HTTP 401 Unauthorized is returned.

---

## Security

Passwords are never:

- stored in plaintext
- logged
- returned to the client

Only bcrypt hashes are persisted.
