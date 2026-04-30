# magic-link-auth

Sistema de autenticação passwordless via magic link implementado em Go 1.22 com Clean Architecture e princípios SOLID.

![Go version](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go) ![License](https://img.shields.io/badge/license-MIT-green)

---

## O que faz

- Recebe um endereço de email e gera um token criptograficamente seguro com expiração de 15 minutos
- Persiste o token em memória com status `PENDING`, `USED` ou `EXPIRED`
- Simula o envio de um magic link por email (log no console em ambiente local)
- Valida o token quanto à existência, expiração e uso único
- Emite um JWT Bearer (HS256, 24h) ao confirmar autenticação bem-sucedida
- Expõe stubs prontos para substituição por DynamoDB e Amazon SES ao migrar para ECS

---

## Arquitetura

```
POST /auth/magic-link          GET /auth/validate?token=
        │                               │
        ▼                               ▼
   Handler (HTTP)              Handler (HTTP)
        │                               │
        ▼                               ▼
   Processor                      Processor
   (valida input,                (valida token
   net/mail)                      não vazio)
        │                               │
        ▼                               ▼
CreateMagicLinkBO           ValidateMagicLinkBO
   (gera token,              (checa status/expiração,
    salva, envia)             marca USED, emite JWT)
        │                               │
   ┌────┴────┐                   ┌──────┴──────┐
   ▼         ▼                   ▼             ▼
Repository  EmailService     Repository  AuthTokenService
(in-memory) (log stub)       (in-memory) (JWT HS256)
```

### Estrutura de pacotes

```
src/
├── functions/api/
│   ├── main.go                           # Entrypoint: wiring manual + servidor HTTP :8080
│   └── handler/
│       ├── create_magic_link_handler.go  # Parse HTTP → delega ao processor
│       ├── validate_magic_link_handler.go
│       └── response.go                   # Helper writeJSON
└── layers/main/
    ├── enums/token_status.go             # TokenStatus: PENDING | USED | EXPIRED
    ├── models/magic_link.go              # Struct de domínio MagicLink
    ├── interfaces/                       # Contratos abstratos (ISP)
    │   ├── magic_link_repository.go      # Save, FindByToken, MarkAsUsed
    │   ├── email_service.go              # SendMagicLink
    │   ├── token_service.go              # Generate
    │   └── auth_token_service.go         # GenerateJWT
    ├── implementations/
    │   ├── memory/                       # Implementações locais ativas
    │   │   ├── magic_link_repository.go  # In-memory com sync.RWMutex
    │   │   ├── token_service.go          # crypto/rand, 32 bytes → hex
    │   │   ├── auth_token_service.go     # JWT HS256 via golang-jwt
    │   │   └── email_service.go          # Log stub
    │   └── aws/                          # Stubs para DynamoDB e SES
    ├── bo/                               # Regras de negócio puras
    │   ├── create_magic_link_bo.go       # Gera token, salva, dispara email
    │   └── validate_magic_link_bo.go     # Valida, marca USED, emite JWT
    └── processor/                        # Validação de input + orquestração
        ├── create_magic_link_processor.go
        └── validate_magic_link_processor.go
```

### Princípios aplicados

| Princípio | Aplicação |
|---|---|
| **Dependency Inversion** | BOs e processors dependem apenas de interfaces — implementações são injetadas em `main.go` |
| **Open/Closed** | Trocar `memory` por `aws` exige alterar apenas o wiring em `main.go`, sem tocar nas regras de negócio |
| **Single Responsibility** | Cada camada tem exatamente uma razão para mudar |
| **Interface Segregation** | Quatro interfaces narrow em vez de uma interface gorda |
| **Injeção manual de dependência** | Sem container de DI — wiring explícito e legível em `main.go` |

---

## Endpoints

| Método | Rota | Input | Resposta de sucesso |
|--------|------|-------|---------------------|
| `POST` | `/auth/magic-link` | `{"email":"user@example.com"}` | `200 {"message":"magic link sent to user@example.com"}` |
| `GET` | `/auth/validate` | `?token=<hex>` | `200 {"access_token":"eyJ...","type":"Bearer"}` |

**Erros do `GET /auth/validate`:**

| HTTP | Condição |
|------|----------|
| `400` | Token ausente ou input inválido |
| `404` | Token não encontrado |
| `422` | Token expirado ou já utilizado |

---

## Desenvolvimento local

### Pré-requisitos

- Go 1.22+

### Setup

```bash
git clone https://github.com/carlos-sousa/magic-link-auth.git
cd magic-link-auth
go mod download
```

### Executar

```bash
go run ./src/functions/api
```

Com variáveis customizadas:

```bash
JWT_SECRET=minha-chave-secreta BASE_URL=http://localhost:8080 go run ./src/functions/api
```

O servidor sobe em `http://localhost:8080`.

### Testes

O projeto não possui testes automatizados neste momento. Para validar o fluxo manualmente, veja a seção de exemplos abaixo.

### Variáveis de ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `JWT_SECRET` | `dev-secret-change-in-production` | Chave HMAC para assinar JWTs |
| `BASE_URL` | `http://localhost:8080` | Base URL usada para montar o magic link |

---

## Exemplos de uso

```bash
# 1. Solicitar magic link
curl -X POST http://localhost:8080/auth/magic-link \
  -H "Content-Type: application/json" \
  -d '{"email":"voce@example.com"}'

# Resposta:
# {"message":"magic link sent to voce@example.com"}

# 2. O log do servidor exibe o token gerado:
# [EMAIL] To: voce@example.com | Magic Link sent (use /auth/validate to authenticate)
# O link completo fica disponível como: http://localhost:8080/auth/validate?token=<hex>

# 3. Validar o token e obter o JWT
curl "http://localhost:8080/auth/validate?token=<token-do-log>"

# Resposta:
# {"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...","type":"Bearer"}
```

---

## Deploy com Docker

```bash
# Build
docker build -t magic-link-auth .

# Executar
docker run -p 8080:8080 -e JWT_SECRET=minha-chave magic-link-auth
```

O Dockerfile usa build multi-stage (`golang:1.22-alpine` → `alpine:3.19`) e executa o processo como usuário não-root.

---

## Segurança

- Token gerado com `crypto/rand` (32 bytes → 64 chars hex)
- Expiração de 15 minutos (Unix timestamp verificado no momento da validação)
- Token de uso único — marcado como `USED` imediatamente após consumo
- JWT HS256 com expiração de 24 horas (`sub`, `iat`, `exp`)
- Email validado com `net/mail` antes de qualquer processamento
- Token completo nunca logado — o stub de email registra apenas a confirmação de envio

---

## Preparado para AWS ECS

Os stubs em `implementations/aws/` fornecem a estrutura para:

- `DynamoDBMagicLinkRepository` — substituir o repositório in-memory por DynamoDB
- `SESEmailService` — substituir o log stub pelo envio real via Amazon SES

Para migrar, basta implementar as interfaces existentes e substituir as dependências no wiring de `main.go`. Nenhuma regra de negócio precisa ser alterada.

---

## Documentação técnica

| Documento | Descrição |
|---|---|
| `src/layers/main/interfaces/` | Contratos que definem os limites entre camadas |
| `src/layers/main/bo/` | Regras de negócio puras, independentes de framework |
| `src/layers/main/implementations/aws/` | Stubs prontos para produção no ECS |
