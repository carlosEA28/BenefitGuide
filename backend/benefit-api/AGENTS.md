# benefit-api (Go)

API principal. Monolito modular em Go com PostgreSQL.

## Leia antes de escrever código

- Estrutura de pastas: `../../docs/golang_architecure.md`
- Padrões de código: `../../docs/conventions/golang.md`

## Comandos (rodar nesta pasta)

- `make run`: sobe a API
- `make build`: gera `bin/benefit-api`
- `make test`: `go test ./...`
- `make fmt`: `go fmt ./...`
- `make tidy`: `go mod tidy`
- `make migration-up`: aplica migrations
- `make migration-rollback`: reverte só a última
- `make migration-down`: reverte TODAS (destrutivo; não rodar sem pedido explícito)

## Migrations

- Ficam em `migrations/` e usam Goose (SQL).
- Nunca edite uma migration já aplicada; crie uma nova.

## Antes de concluir uma tarefa

Rode `make fmt` e `make test`.
