# magic-link-auth

Passwordless authentication service built in Go 1.22 that issues magic links via email and exchanges them for signed JWTs.

![Go version](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go) ![License](https://img.shields.io/badge/license-MIT-green)

---

## What it does

- Accepts an email address, generates a cryptographically secure token (32 bytes, hex-encoded), and persists it with a 15-minute expiry
- Sends a magic link to the user — logged to stdout in local mode, delivered via SMTP or Amazon SES in production
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
(in-memory)  (log/smtp)     (in-memory)  (JWT HS256)
```

All BOs and processors depend exclusively on interfaces. Concrete implementations are injected in `src/containers/server.go`, keeping business rules free of infrastructure concerns.

### Project structure

```
.
├── src/
│   ├── containers/
│   │   ├── server.go                    # Entry point: DI wiring and HTTP server
│   │   └── env.go                       # Config struct loaded from environment variables
│   └── layers/main/
│       ├── enums/
│       │   └── token_status.go          # TokenStatus: PENDING | USED | EXPIRED
│       ├── models/
│       │   └── magic_link.go            # MagicLink domain struct
│       ├── interfaces/                  # Abstract contracts (ISP)
│       │   ├── controller.go
│       │   ├── magic_link_dao.go        # Save, FindByToken, MarkAsUsed
│       │   ├── email_service.go         # SendMagicLink
│       │   ├── token_service.go         # Generate
│       │   ├── auth_token_service.go    # GenerateJWT
│       │   └── secrets_service.go       # GetSecret
│       ├── bo/                          # Pure business rules, framework-independent
│       │   ├── create_magic_link_bo.go
│       │   └── validate_magic_link_bo.go
│       ├── processor/                   # Input validation + BO orchestration
│       │   ├── create_magic_link_processor.go
│       │   └── validate_magic_link_processor.go
│       ├── controller/                  # HTTP parsing and response formatting
│       │   ├── create_magic_link_controller.go
│       │   ├── validate_magic_link_controller.go
│       │   └── response.go
│       └── implementations/
│           ├── memory/                  # Local stubs (active by default)
│           │   ├── magic_link_dao.go    # sync.RWMutex in-memory store
│           │   ├── token_service.go     # crypto/rand → 64-char hex
│           │   ├── auth_token_service.go# JWT HS256 via golang-jwt
│           │   ├── email_service.go     # Logs magic link to stdout
│           │   └── secrets_service.go   # Map-backed local secrets
│           ├── smtp/
│           │   └── email_service.go     # Sends via SMTP (local MailPit or real relay)
│           └── aws/                     # Production implementations
│               ├── magic_link_dao.go    # DynamoDB
│               ├── email_service.go     # Amazon SES v2
│               └── secrets_service.go   # AWS Secrets Manager
├── tests/
│   ├── testutil/
│   │   └── mocks.go                     # Shared mock structs for all interfaces
│   └── unit/
│       ├── bo/
│       │   ├── create_magic_link_bo_test.go
│       │   └── validate_magic_link_bo_test.go
│       └── processor/
│           ├── create_magic_link_processor_test.go
│           └── validate_magic_link_processor_test.go
└── ui/
    └── index.html                       # Minimal frontend for end-to-end validation
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
| `GET` | `/` | Frontend UI (served from `ui/`) |

**`POST /auth/magic-link`**

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
- (Optional) [MailPit](https://github.com/axllent/mailpit) for end-to-end email validation

### Run

```bash
go run ./src/containers
```

The magic link token is printed to stdout by `LogEmailService`:

```
[EMAIL] To: user@example.com | Magic Link: http://localhost:8080?token=<hex>
```

Open `http://localhost:8080` in the browser to use the UI.

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `JWT_SECRET` | `dev-secret-change-in-production` | HMAC key used to sign JWTs |
| `BASE_URL` | `http://localhost:8080` | Base URL embedded in the magic link sent by email |
| `PORT` | `8080` | HTTP server port |
| `SMTP_HOST` | _(empty)_ | SMTP host — if set, enables real email delivery |
| `SMTP_PORT` | `1025` | SMTP port |
| `SMTP_FROM` | `noreply@localhost` | Sender address used in outgoing emails |

When `SMTP_HOST` is empty the service falls back to `LogEmailService`, which only prints the link to stdout.

---

## Testing

### Unit tests

Tests live in `tests/unit/` and are completely separated from production code. Shared mock structs for all interfaces are defined in `tests/testutil/mocks.go`.

**What is tested:**

| Package | Test cases |
|---|---|
| `bo` | Token generation error, DAO save error, email send error, success |
| `bo` | Token not found, already used, status expired, time expired, mark-as-used error, JWT error, success |
| `processor` | Empty email, invalid email format, BO error propagation, success |
| `processor` | Empty token, token not found propagation, success + Bearer type |

**What is not tested:** concrete implementations (`memory`, `smtp`, `aws`) — these are infrastructure and tested through end-to-end or integration flows.

**Run all unit tests:**

```bash
go test ./tests/...
```

**Run with verbose output:**

```bash
go test ./tests/... -v
```

**Run a specific package:**

```bash
go test ./tests/unit/bo/...
go test ./tests/unit/processor/...
```

**Run a single test:**

```bash
go test ./tests/unit/bo/... -run TestCreateMagicLinkBO_Success
```

### End-to-end with MailPit

MailPit provides a local SMTP server with a web inbox, allowing you to receive the magic link email and complete the full authentication flow in the browser.

**1. Start MailPit:**

```bash
docker run -d -p 1025:1025 -p 8025:8025 axllent/mailpit
```

**2. Start the server with SMTP enabled:**

```bash
SMTP_HOST=localhost go run ./src/containers
```

**3. Open the UI and request a magic link:**

```
http://localhost:8080
```

**4. Check the inbox for the email:**

```
http://localhost:8025
```

**5. Click the link in the email.** The browser opens `http://localhost:8080?token=<hex>`, the frontend calls `GET /auth/validate`, and displays the JWT on screen.

### Manual curl flow

If you prefer the terminal without MailPit:

```bash
# 1. Request a magic link (token is printed in server stdout)
curl -X POST http://localhost:8080/auth/magic-link \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# 2. Validate the token (copy from server log)
curl "http://localhost:8080/auth/validate?token=<token-from-log>"
```

---

## Docker

```bash
# Build
docker build -t magic-link-auth .

# Run (log-only mode)
docker run -p 8080:8080 \
  -e JWT_SECRET=change-me \
  -e BASE_URL=https://your-domain.com \
  magic-link-auth

# Run with SMTP
docker run -p 8080:8080 \
  -e JWT_SECRET=change-me \
  -e BASE_URL=https://your-domain.com \
  -e SMTP_HOST=your-smtp-host \
  -e SMTP_PORT=587 \
  -e SMTP_FROM=noreply@your-domain.com \
  magic-link-auth
```

---

## Deploy (AWS ECS)

Swap the `memory` implementations for `aws` in `src/containers/server.go`:

| Local | Production | Notes |
|---|---|---|
| `memory.NewInMemoryMagicLinkDAO()` | `aws.NewDynamoDBMagicLinkDAO(...)` | Requires table name and DynamoDB client |
| `memory.NewLogEmailService()` | `aws.NewSESEmailService(...)` | Requires sender address and SES client |
| `memory.NewJWTAuthTokenService(...)` | `memory.NewJWTAuthTokenService(...)` | No change needed |

No business logic changes — only the wiring in `server.go` is updated.

---

## Security

- Tokens are generated with `crypto/rand` (32 bytes encoded as 64-character hex strings)
- Tokens expire after 15 minutes (Unix timestamp checked at validation time)
- Tokens are single-use — marked `USED` immediately after the first successful validation
- JWTs are signed with HS256 and expire after 24 hours (`sub`, `iat`, `exp` claims)
- Email addresses are validated with `net/mail` before any processing
- The token value is never logged — only the complete magic link URL is printed in local mode
