# Go Modular Monolith Architecture

## Visão Geral

Este repositório contém a API **Benefit Guide**, um sistema monolítico modular escrito em Go. O objetivo é organizar o código em módulos independentes e interchangeáveis, cada um com sua própria responsabilidade (auth, catalog, consultations, documents), mas todos compartilhando uma base comum de configuração e banco de dados.

A arquitetura segue princípios de **modular monolith** (monolito modular), onde:

- Cada módulo encapsula o seu próprio domínio, com repository, service e HTTP handler.
- As dependências entre módulos são explicitas via parâmetros de função (injeção de serviço).
- O ponto de entrada (`main.go`) orquestra a criação dos módulos e o registro das rotas no router Chi.
- A configuração e a conexão com o banco de dados são compartilhadas via pacote `shared`.

---

## Estrutura de Pastas

```
benefitGuide/
├── backend/
│   ├── ai-service/               # Serviço de IA em Python (fora do escopo Go)
│   └── benefit-api/              # API Go (modular monolith)
│       ├── go.mod
│       ├── go.sum
│       ├── main.go               # Entry point: configura DB, modules e HTTP server
│       ├── Makefile
│       └── internal/
│           ├── auth/             # Módulo de autenticação
│           │   ├── module.go     # Estrutura Module, New, RegisterRoutes
│           │   ├── internal/
│           │   │   ├── repository/   # CRUD ao banco
│           │   │   ├── service/      # Lógica de negócio
│           │   │   └── server/       # Handlers Chi
│           │   ├── cmd/            # Comandos/flags (se houver)
│           │   └── module.go
│           ├── catalog/            # Módulo de catálogo de benefícios
│           │   ├── module.go
│           │   ├── internal/
│           │   │   ├── repository/
│           │   │   ├── service/
│           │   │   └── server/
│           ├── consultations/      # Módulo de consulta de benefícios
│           │   ├── module.go
│           │   ├── internal/
│           │   │   ├── repository/
│           │   │   ├── service/      # Depende de auth e catalog services
│           │   │   └── server/
│           ├── documents/          # Módulo de gestão de documentos
│           │   ├── module.go
│           │   ├── internal/
│           │   │   ├── repository/
│           │   │   ├── service/      # Depende do service de consultations
│           │   │   └── server/
│           └── shared/             # Componentes compartilhados
│               ├── config/           # Load .env, Config struct, DatabaseConfig
│               │   ├── application.go
│               │   └── config.go
│               ├── db/               # Conexão GORM com PostgreSQL
│               │   └── database.go
│               └── exceptions/       # (pasta vazia atualmente, reservada para erros padronizados)
├── migrations/                      # Migrações de banco (Ex: Airplane, MigraGo, etc.)
└── docs/                            # Documentação gerada
    └── go_architecture.md
```

### Detalhamento dos Módulos Internos

Cada módulo (`auth`, `catalog`, `consultations`, `documents`) segue o mesmo pattern estrutural definido em seu `module.go`:

| Arquivo | Responsabilidade |
|---|---|
| `module.go` | Define a struct `Module`, a função `New(db *gorm.DB)` (ou com serviços de módulos dependentes) e o método `RegisterRoutes(r chi.Router)` |
| `internal/repository/` | Implementação de repositórios que abstraem o acesso ao banco via GORM |
| `internal/service/` | Lógica de negócio específica do domínio |
| `internal/server/` | Handlers HTTP (Chi) que expõem endpoints e conectam o handler ao service |

#### Módulo `auth`
- **Rotas**: `/auth`
- **Dependências**: Apenas o banco de dados (`*gorm.DB`).
- **Funcionalidades**: Registro, login, autenticação de usuários.

#### Módulo `catalog`
- **Rotas**: `/catalog`
- **Dependências**: Apenas o banco de dados.
- **Funcionalidades**: CRUD de catálogos de benefícios.

#### Módulo `consultations`
- **Rotas**: `/consultations`
- **Dependências**: `*authmodule.Service` e `*catalogmodule.Service` (injeção de serviço).
- **Funcionalidades**: Consultas de benefícios que podem exigir informações de autenticação e dados do catálogo.

#### Módulo `documents`
- **Rotas**: `/documents`
- **Dependências**: `*consultationsmodule.Service`.
- **Funcionalidades**: Upload, listagem e gestão de documentos relacionados a consultas.

### Pacote `shared`

- **`shared/config`**: Carrega variáveis de ambiente do arquivo `.env` e fornece o struct `Config` ( `ServerPort`, `Database` ) e funções auxiliares (`Load`, `envOrDefault`, `loadEnvFile`).
- **`shared/db`**: Cria a conexão com PostgreSQL usando GORM, a partir do `DatabaseConfig` lido pela config. Configura pools de conexões, timeouts e logger em modo `info`.

### Ponto de Entrada (`main.go`)

O arquivo `main.go` orquestra tudo:

1. Carrega a configuração (`config.Load`).
2. Abre a conexão com o banco de dados (`db.New`).
3. Cria instâncias dos módulos:
   - `authModule := auth.New(database)`
   - `catalogModule := catalog.New(database)`
   - `consultationsModule := consultations.New(database, authModule.Service, catalogModule.Service)`
   - `documentsModule := documents.New(database, consultationsModule.Service)`
4. Cria o router Chi e utiliza middlewares (RequestID, RealIP, Logger, Recoverer, CORS).
5. Registra as rotas de cada módulo: `authModule.RegisterRoutes(router)`, `catalogModule.RegisterRoutes(router)`, etc.
6. Inicia o servidor HTTP na porta configurada.

---

## Padrões e Decisões de Arquitetura

| Aspecto | Decisão |
|---|---|
| **Tipo de arquitetura** | Modular Monolith: um único processo/binário, mas com frontières claras entre módulos. |
| **Roteamento** | `github.com/go-chi/chi/v5` — leve, rápido, suporte a routes aninhadas. |
| **ORM** | `gorm.io/gorm` com driver `gorm.io/driver/postgres`. Permite migrações simples e queries convenientes. |
| **Injeção de dependência** | Manual (via parâmetros de função) dentro de `main.go`. Isso mantém o código explícito e fácil de testar. |
| **Configuração** | Arquivo `.env` lido no startup; variáveis prefixadas (`SERVER_PORT`, `DB_HOST`, etc.). |
| **Comunicação entre módulos** | Injeção de interfaces (services) — um módulo passa o seu `Service` para outro ao ser instanciado. |
| **Testabilidade** | Cada módulo pode ser testado isoladamente mockando o repositório e o banco; o `main.go` facilita a criação de fixtures com conexão em memória ou mock. |
| **Extensibilidade** | Para adicionar um novo módulo, basta criar a pasta `internal/<nome>`, implementar `module.go`, `repository/`, `service/`, `server/` e registrá-lo em `main.go`. |

---

## Como Executar

```bash
# 1. Certifique‑se de ter o Go 1.26.1+ instalado
go version

# 2. Clone o repositório e entre na pasta backend
cd benefitGuide/backend/benefit-api

# 3. Instale as dependências
go mod tidy

# 4. Renomeie .env.example para .env e configure as variáveis
cp .env.example .env
# edite .env com seus valores de DB, porta, etc.

# 5. Rode a aplicação
go run main.go
# ou use o Makefile incluso (make run, make build, etc.)
```
