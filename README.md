# magic-link-auth

Passwordless authentication in Go. User enters their email, receives a magic link, clicks it, gets a JWT.

---

## How it works

```
1. POST /auth/magic-link  { "email": "user@example.com" }
        │
        ├── generates a secure random token (32 bytes hex)
        ├── saves it with a 15-minute expiry
        └── sends an email with the link: http://localhost:8080?token=<hex>

2. User clicks the link → browser opens the frontend → frontend calls GET /auth/validate?token=<hex>
        │
        ├── checks token exists, not expired, not already used
        ├── marks token as USED
        └── returns a signed JWT (HS256, 24h)
```

---

## Running locally

**Log-only mode** (no email server needed — token is printed to stdout):

```bash
go run ./src/containers
```

**With real email via MailPit:**

```bash
# Start MailPit (SMTP on :1025, inbox UI on :8025)
docker run -d -p 1025:1025 -p 8025:8025 axllent/mailpit

# Start the server
SMTP_HOST=localhost go run ./src/containers
```

Then:
1. Open `http://localhost:8080` → enter your email
2. Open `http://localhost:8025` → click the link in the inbox
3. Browser redirects back and shows the JWT

---

## Running tests

```bash
go test ./tests/...
```

Tests cover the `bo` and `processor` layers using manual mocks. Concrete implementations (`memory`, `smtp`, `aws`) are not unit tested.

---

## Project structure

```
src/
├── containers/         # Entry point and env config
└── layers/main/
    ├── interfaces/     # Contracts: DAO, EmailService, TokenService, AuthTokenService
    ├── bo/             # Business rules (token generation, validation)
    ├── processor/      # Input validation + BO orchestration
    ├── controller/     # HTTP layer
    └── implementations/
        ├── memory/     # Local stubs (default)
        ├── smtp/       # SMTP email delivery
        └── aws/        # DynamoDB, SES, Secrets Manager (production)
tests/
├── testutil/           # Shared mock structs
└── unit/               # Unit tests for bo/ and processor/
ui/
└── index.html          # Minimal frontend
```

---

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `JWT_SECRET` | `dev-secret-change-in-production` | JWT signing key |
| `BASE_URL` | `http://localhost:8080` | Base URL used in the magic link |
| `PORT` | `8080` | HTTP port |
| `SMTP_HOST` | _(empty)_ | If set, enables real email via SMTP |
| `SMTP_PORT` | `1025` | SMTP port |
| `SMTP_FROM` | `noreply@localhost` | Sender address |

---

## Deploying to AWS

Swap implementations in `src/containers/server.go`:

| Local | Production |
|---|---|
| `memory.NewInMemoryMagicLinkDAO()` | `aws.NewDynamoDBMagicLinkDAO(...)` |
| `memory.NewLogEmailService()` | `aws.NewSESEmailService(...)` |

No business logic changes required.
