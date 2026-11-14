# Convenções Go - Benefit Guide Monolito Modular

Este documento descreve as convenções de código para o projeto Go **Benefit Guide**. O objetivo é padronizar a estrutura para que exemplos de código possam ser inseridos e seguidos consistentemente por uma IA (ou desenvolvedores).

## 1. Controller

Responsável por receber requisições HTTP, validar entrada e delegar à camada de serviço.

- Utiliza **Chi** (`github.com/go-chi/chi/v5`) para roteamento.
- Deve conter apenas o "boilerplate" de mapear parâmetros e retornar respostas JSON.
- Nunca contém lógica de negócio — tudo deve ser delegado ao **Service**.

**Estrutura esperada:**

```go
// internal/auth/internal/server/handler.go
package auth

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/carlosEA28/benefitGuide/internal/auth/service"
)

// Handler struct holds the service dependencies.
type Handler struct {
    Service service.Service
}

// NewHandler cria um novo handler.
func NewHandler(s service.Service) *Handler {
    return &Handler{Service: s}
}

// RegisterRoutes injeta as rotas no router Chi.
func (h *Handler) RegisterRoutes(r chi.Router) {
    r.Method(http.MethodPost, "/login", h.Login)
    r.Method(http.MethodPost, "/refresh-token", h.RefreshToken)
    // ... mais rotas
}

// Login lida com POST /auth/login (credenciais de login)
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    // TODO: inserir código de parseio de LoginDTO e chamada h.Service.Login
    writeJSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
}

// RefreshToken lida com POST /auth/refresh-token (renovação de JWT)
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    // TODO: inserir código de refresh token e resposta
    writeJSON(w, http.StatusOK, map[string]interface{}{"status": "refreshed"})
}
```

---

## 2. DTO (Data Transfer Object)

Objetos para transferir dados entre as camadas (entrada/saída da API). Evita expor modelos de banco diretamente.

- **Request DTO**: usado para validar e receber dados da requisição HTTP.
- **Response DTO**: usado para moldar a resposta JSON devolvida ao cliente.
- Use `json` tags nos campos.
- Pode usar bibliotecas como `go-playground/validator` para validação, mas a convenção aqui é manter simples.

**Exemplo de estrutura:**

```go
// internal/shared/dto/request.go
package dto

// LoginDTO recebido na requisição de autenticação
type LoginDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// RefreshTokenDTO recebido para renovação de token
type RefreshTokenDTO struct {
    Token string `json:"token" validate:"required"`
}

// TokenDTO resposta contendo o token gerado
type TokenDTO struct {
    Token string `json:"token"`
}

// UserResponseDTO resposta de usuário (não senha!)
type UserResponseDTO struct {
    ID    int64  `json:"id"`
    Email string `json:"email"`
}
```

---

## 3. Interfaces

Definições de contratos que permitem injeção de dependência e facilita testes unitários.

- **Repository Interface**: método de CRUD mínimo (`Create`, `GetByID`, `List`, `Update`, `Delete`).
- **Service Interface**: método de lógica de negócio exposto pelo handler.
- As interfaces devem ser definidas no pacote **próprio** ou em um pacote `internal`, mas **nunca** no `main`.

**Exemplo:**

```go
// internal/auth/repository/repository.go
package repository

type RepositoryInterface interface {
    Create(*auth.Auth) (*auth.Auth, error)
    GetByID(id int64) (*auth.Auth, error)
    GetByEmail(email string) (*auth.Auth, error)
    // Atualiza um registro
    Update(*auth.Auth) error
    // Remove um registro por ID
    Delete(id int64) error
}
```

---

## 4. Models (Entidades / Domínio)

Structs que representam a estrutura de dados do banco de dados (usando GORM tags).

- Definir na raiz do módulo ou em `internal/<module>/model.go`.
- Usar tags `gorm:"..."` para mapeamento.
- Evitar lógica de negócio aqui — apenas estrutura de dados e relações.

**Exemplo:**

```go
// internal/auth/model.go
package auth

import "time"

type Auth struct {
    ID        int64     `gorm:"primaryKey;autoIncrement"`
    Email     string    `gorm:"size:255;uniqueIndex"`
    Password  string    `gorm:"size:255"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

---

## 5. Repository

Implementação concreta de acesso ao banco de dados usando **GORM**.

- Métodos devem seguir a interface definida na seção de **Interfaces**.
- Receber `*gorm.DB` via construtor.
- Manter transações explícitas quando necessário.

**Exemplo:**

```go
// internal/auth/repository/repository_impl.go
package repository

import (
    "github.com/carlosEA28/benefitGuide/internal/auth/model"
    "gorm.io/gorm"
)

type repository struct {
    db *gorm.DB
}

func New(db *gorm.DB) *repository {
    return &repository{db: db}
}

// Create insere um novo registro no banco
func (r *repository) Create(entity *model.Auth) (*model.Auth, error) {
    if err := r.db.Create(entity).Error; err != nil {
        return nil, err
    }
    return entity, nil
}

// GetByID busca um registro pelo ID principal
func (r *repository) GetByID(id int64) (*model.Auth, error) {
    var entity model.Auth
    if err := r.db.First(&entity, id).Error; err != nil {
        return nil, err
    }
    return &entity, nil
}

// GetByEmail busca um usuário pelo e-mail (único)
func (r *repository) GetByEmail(email string) (*model.Auth, error) {
    var entity model.Auth
    if err := r.db.Where("email = ?", email).First(&entity).Error; err != nil {
        return nil, err
    }
    return &entity, nil
}
```

---

## 6. Server (Configuração e Middleware)

Configuração do servidor HTTP, roteador Chi e middlewares compartilhados.

- O `main.go` inicia o servidor aqui.
- Middlewares comuns: RequestID, RealIP, Logger, Recoverer, CORS.
- Evitar colocar rotas específicas aqui; elas pertencem aos **Controllers**.

**Exemplo (em main.go ou separate file):**

```go
// internal/shared/server/router.go
package server

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/carlosEA28/benefitGuide/internal/auth/controller"
    "github.com/carlosEA28/benefitGuide/internal/catalog/controller"
    // ... outros controllers
)

func NewRouter(authCtrl *auth.Controller, catalogCtrl *catalog.Controller) *chi.Mux {
    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(corsMiddleware) // definido abaixo

    r.Route("/auth", func(r chi.Router) {
        authCtrl.RegisterRoutes(r)
    })
    r.Route("/catalog", func(r chi.Router) {
        catalogCtrl.RegisterRoutes(r)
    })
    // ... mais rotas

    return r
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

## 7. Service

Lógica de negócio central. Coordenador entre repositórios e outros serviços.

- Receber repositórios (e outros serviços) via construtor ou parâmetros de função.
- Retornar DTOs (ou modelos transformados) para o controller.
- Lancar erros de domínio (ex.: usando pacote `errors` ou erros customizados do `shared/exceptions`).

**Exemplo:**

```go
// internal/auth/service/service_impl.go
package service

import (
    "errors"
    "time"
    "github.com/carlosEA28/benefitGuide/internal/auth/model"
    "github.com/carlosEA28/benefitGuide/internal/auth/repository"
    "github.com/carlosEA28/benefitGuide/shared/exceptions"
)

type service struct {
    repo repository.RepositoryInterface
}

func New(repo repository.RepositoryInterface) *service {
    return &service{repo: repo}
}

// Login realiza autenticação de usuário com e-mail/senha
func (s *service) Login(email, password string) (*service.LoginResult, error) {
    user, err := s.repo.GetByEmail(email)
    if err != nil {
        return nil, exceptions.ErrInvalidCredentials
    }
    // validar senha (comparar hash)
    if !checkPasswordHash(password, user.Password) {
        return nil, exceptions.ErrInvalidCredentials
    }
    // gerar token JWT com expiração em 24h
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })
    signedToken, err := token.SignedString([]byte("sua-chave-secreta"))
    if err != nil {
        return nil, exceptions.ErrTokenGeneration
    }
    return &service.LoginResult{Token: signedToken}, nil
}

// RefreshToken renova um token de acesso usando um refresh token válido
func (s *service) RefreshToken(refreshToken string) (*service.RefreshResult, error) {
    // validar refresh token, gerar novo access token
    // ...
    return &service.RefreshResult{NewToken: "novo-token"}, nil
}

// Logout invalida um token ou refresh token (dependendo da implementação)
func (s *service) Logout() error {
    // implementação depende da estratégia de revogação
    return nil
}

// Helper privado para comparar hash de senha
func (s *service) checkPasswordHash(password, hash string) bool {
    // usar bcrypt.CompareHashAndPassword ou equivalente
    return true // placeholder
}

// Estruturas de resultado
type LoginResult struct {
    Token string `json:"token"`
}

type RefreshResult struct {
    NewToken string `json:"new_token"`
}
```

---

## 8. Utils

Funções auxiliares e helpers usados em toda a aplicação.

- Funções de escrita JSON (`writeJSON`).
- Funções de hash de senha.
- Formatação de data/hora.
- anything que seja reutilizado across modules mas que não justifique um pacote próprio.

**Exemplo:**

```go
// internal/shared/utils.go
package utils

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"
)

// writeJSON escreve uma resposta JSON com status code.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(v)
}

// SanitizeString remove espaços extras e normaliza.
func SanitizeString(s string) string {
    return strings.TrimSpace(strings.Trim(s, " "))
}

// IsValidEmail verifica formato básico de e-mail.
func IsValidEmail(e string) bool {
    return strings.Contains(e, "@")
}

// FormatDateISO formata um time.Time para string ISO 8601
func FormatDateISO(t time.Time) string {
    return t.Format(time.RFC3339)
}
```

---

## Como usar este documento

1. Copie os esboços acima para os respectivos arquivos da sua estrutura `internal/<module>/`.
2. Preencha os `TODO` e implemente a lógica conforme o domínio do Benefit Guide.
3. Mantenha os nomes de pacote e arquivos consistentes com as convenções aqui definidas.
4. Quando for inserir exemplos para a IA, basta adicionar o código real nas respectivas seções — a IA seguirá o padrão definido.

---
*Última atualização: $(date +%Y-%m-%d)*