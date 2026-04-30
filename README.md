# magic-link-auth

Passwordless authentication service built in Go 1.22 that issues magic links via email and exchanges them for signed JWTs.

![Go version](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go) ![License](https://img.shields.io/badge/license-MIT-green)

---

## What it does

- Accepts an email address, generates a cryptographically secure token (32 bytes, hex-encoded), and persists it with a 15-minute expiry
- Sends a magic link to the user (`BASE_URL/auth/validate?token=<hex>`) — logged to stdout in local mode, delivered via Amazon SES in production
- Validates the token against existence, expiry, and single-use constraints
- Marks the token as `USED` immediately on first successful validation
- Issues a signed JWT (HS256, 24 h) upon successful authentication
- Provides ready-to-swap AWS implementations (`DynamoDBMagicLinkDAO`, `SESEmailService`) that replace in-memory stubs when deploying to ECS

---

## Architecture

### Request flow

```
POST /auth/magic-link          GET /auth/validate?token=
        |                               |
        v                               v
   Controller                     Controller
        |                               |
        v                               v
   Processor                       Processor
   (validates email                (checks token
    with net/mail)                  is not empty)
        |                               |
        v                               v
CreateMagicLinkBO           ValidateMagicLinkBO
  (generates token,           (checks status/expiry,
   saves, sends link)          marks USED, issues JWT)
        |                               |
   +----+----+                   +------+------+
   v         v                   v             v
MagicLinkDAO EmailService   MagicLinkDAO  AuthTokenService
(in-memory)  (log stub)     (in-memory)  (JWT HS256)
```

All BOs and processors depend exclusively on interfaces. Concrete implementations are injected in `src/containers/server.go`, keeping business rules free of infrastructure concerns.

### Package structure

```
src/
├── containers/
│   ├── server.go           # Entry point: manual DI wiring and HTTP server
│   └── env.go              # Config struct loaded from environment variables
└── layers/main/
    ├── enums/
    │   └── token_status.go # TokenStatus: PENDING | USED | EXPIRED
    ├── models/
    │   └── magic_link.go   # MagicLink domain struct
    ├── interfaces/         # Abstract contracts (ISP)
    │   ├── controller.go         # Controller: Handle(w, r)
    │   ├── magic_link_dao.go     # MagicLinkDAO: Save, FindByToken, MarkAsUsed
    │   ├── email_service.go      # EmailService: SendMagicLink
    │   ├── token_service.go      # TokenService: Generate
    │   └── auth_token_service.go # AuthTokenService: GenerateJWT
    ├── implementations/
    │   ├── memory/               # Active local implementations
    │   │   ├── magic_link_dao.go    # In-memory store with sync.RWMutex
    │   │   ├── token_service.go     # crypto/rand, 32 bytes -> hex
    │   │   ├── auth_token_service.go# JWT HS256 via golang-jwt
    │   │   └── email_service.go     # Logs magic link to stdout
    │   └── aws/                  # Production stubs (DynamoDB, SES)
    │       ├── magic_link_dao.go    # DynamoDBMagicLinkDAO (not yet implemented)
    │       └── email_service.go     # SESEmailService (not yet implemented)
    ├── bo/                       # Pure business rules
    │   ├── create_magic_link_bo.go  # Generates token, saves, dispatches email
    │   └── validate_magic_link_bo.go# Validates token, marks USED, issues JWT
    ├── processor/                # Input validation and orchestration
    │   ├── create_magic_link_processor.go
    │   └── validate_magic_link_processor.go
    └── controller/               # HTTP parsing and response formatting
        ├── create_magic_link_controller.go
        ├── validate_magic_link_controller.go
        └── response.go              # writeJSON helper
```

### SOLID principles applied

| Principle | Where |
|---|---|
| Dependency Inversion | BOs and processors depend only on interfaces; implementations are injected in `server.go` |
| Open/Closed | Swapping `memory` for `aws` requires changing only the wiring in `server.go` — no business logic changes |
| Single Responsibility | Each layer has exactly one reason to change |
| Interface Segregation | Five narrow interfaces instead of one broad contract |

---

## Endpoints

| Method | Route | Description |
|---|---|---|
| `POST` | `/auth/magic-link` | Request a magic link for the given email |
| `GET` | `/auth/validate` | Validate a token and receive a JWT |

**`POST /auth/magic-link`**

Request body:
```json
{"email": "user@example.com"}
```

Response `200`:
```json
{"message": "magic link sent to user@example.com"}
```

**`GET /auth/validate?token=<hex>`**

Response `200`:
```json
{"access_token": "eyJhbGci...", "type": "Bearer"}
```

Error responses:

| HTTP | Condition |
|---|---|
| `400` | Missing or empty token, invalid request body, invalid email |
| `404` | Token not found |
| `422` | Token expired or already used |

---

## Local development

### Prerequisites

- Go 1.22+

### Setup

```bash
git clone <repo-url>
cd magic-link-auth
go mod download
```

### Run

```bash
go run ./src/containers
```

With custom environment variables:

```bash
JWT_SECRET=my-secret BASE_URL=http://localhost:8080 PORT=8080 go run ./src/containers
```

The server starts on `http://localhost:8080`. The magic link token is printed to stdout by `LogEmailService`:

```
[EMAIL] To: user@example.com | Magic Link: http://localhost:8080/auth/validate?token=<hex>
```

### Tests

The project does not have automated tests at this time. To validate the flow manually:

```bash
# 1. Request a magic link
curl -X POST http://localhost:8080/auth/magic-link \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# 2. Copy the token printed in the server log, then validate it
curl "http://localhost:8080/auth/validate?token=<token-from-log>"
```

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `JWT_SECRET` | `dev-secret-change-in-production` | HMAC key used to sign JWTs |
| `BASE_URL` | `http://localhost:8080` | Base URL used to build the magic link |
| `PORT` | `8080` | HTTP server port |

---

## Deploy

### Docker

```bash
# Build
docker build -t magic-link-auth .

# Run
docker run -p 8080:8080 \
  -e JWT_SECRET=change-me \
  -e BASE_URL=https://your-domain.com \
  magic-link-auth
```

The Dockerfile uses a multi-stage build (`golang:1.22-alpine` builder -> `alpine:3.19` runtime) and runs the process as a non-root user.

### AWS ECS

To deploy to ECS, swap the `memory` implementations for `aws` ones in `src/containers/server.go`:

| Local | Production | Notes |
|---|---|---|
| `memory.NewInMemoryMagicLinkDAO()` | `aws.DynamoDBMagicLinkDAO` | Requires table name and DynamoDB client |
| `memory.NewLogEmailService()` | `aws.SESEmailService` | Requires sender address and SES client |
| `memory.NewJWTAuthTokenService()` | `memory.NewJWTAuthTokenService()` | Infrastructure-agnostic — no change needed |

No business logic needs to change. Only the wiring in `server.go` is updated.

---

## Security

- Tokens are generated with `crypto/rand` (32 bytes encoded as 64-character hex strings)
- Tokens expire after 15 minutes (Unix timestamp checked at validation time)
- Tokens are single-use — marked `USED` immediately after the first successful validation
- JWTs are signed with HS256 and expire after 24 hours (`sub`, `iat`, `exp` claims)
- Email addresses are validated with `net/mail` before any processing
- The full token value is never logged — `LogEmailService` prints only the complete magic link URL

---

## Technical reference

| Path | Description |
|---|---|
| `src/containers/server.go` | Entry point and manual dependency injection wiring |
| `src/containers/env.go` | Environment variable loading (`JWT_SECRET`, `BASE_URL`, `PORT`) |
| `src/layers/main/interfaces/` | Abstract contracts defining layer boundaries |
| `src/layers/main/bo/` | Pure business rules, framework-independent |
| `src/layers/main/implementations/aws/` | Production stubs ready to implement for ECS |
