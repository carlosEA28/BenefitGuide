# Convenções Go - Benefit Guide Monolito Modular

Este documento descreve as convenções de código para o projeto Go **Benefit Guide**. O objetivo é padronizar a estrutura para que exemplos de código possam ser inseridos e seguidos consistentemente por uma IA (ou desenvolvedores).

## 1. Handler HTTP

O handler fica em `internal/<modulo>/internal/server`, recebe uma interface pequena definida no mesmo pacote e usa Chi com `net/http`. Cada método lê parâmetros e corpo, chama o serviço, mapeia erros conhecidos com `errors.Is` e escreve a resposta. Validação de domínio e regras de negócio ficam no serviço.

Use `chi.URLParam` para parâmetros de rota. Decodifique JSON com `encoding/json` e os helpers do pacote, sem aceitar campos desconhecidos nem conteúdo extra quando esse for o contrato da rota. Use DTOs de entrada e saída; não exponha o modelo de persistência diretamente. Preserve os códigos, mensagens e formato JSON de cada endpoint.

Trecho baseado nos handlers de OfficialSource do Catalog (outros métodos da interface, imports e funções auxiliares do mesmo pacote omitidos):

```go
type officialSourceService interface {
    CreateOfficialSource(context.Context, *models.OfficialSource) (*models.OfficialSource, error)
}

type sourceHandler struct {
    service officialSourceService
}

func newSourceHandler(service officialSourceService) *sourceHandler {
    return &sourceHandler{service: service}
}

func (h *sourceHandler) createSource(w http.ResponseWriter, r *http.Request) {
    serviceID, err := parseUUID(r, "serviceId") // usa chi.URLParam
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    source, err := decodeSource(r) // DTO JSON e validação da entrada
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
        return
    }
    source.PublicServiceID = serviceID

    created, err := h.service.CreateOfficialSource(r.Context(), source)
    if err != nil {
        switch {
        case errors.Is(err, exceptions.ErrInvalidOfficialSource):
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid official source"})
        case errors.Is(err, exceptions.ErrPublicServiceNotFound):
            writeJSON(w, http.StatusNotFound, map[string]string{"error": "public service not found"})
        default:
            writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
        }
        return
    }
    writeJSON(w, http.StatusCreated, sourceResponseFor(created))
}
```

O Catalog concentra esse mapeamento no helper `writeServiceError`, compartilhado pelos handlers de requirements e sources. Um helper é adequado quando várias rotas seguem o mesmo contrato de erro. O `Handler` exportado compõe os handlers de recurso e registra as rotas com `r.Post`, `r.Get`, `r.Put` e `r.Delete`. Os nomes exportados necessários ao `module.go` permanecem exportados.

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

Definições de contratos para injeção de dependência.

- **Repository Interface**: método de CRUD mínimo (`Create`, `GetByID`, `List`, `Update`, `Delete`).
- **Service Interface**: declare junto ao handler, no pacote consumidor `server`, apenas os métodos que ele chama. O construtor recebe essa interface.
- Evite interfaces de serviço globais quando os recursos usam conjuntos diferentes de métodos. No Catalog, `requirementService` e `officialSourceService` são separados.

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

- `main.go` configura o roteador, os middlewares compartilhados e o servidor HTTP.
- Cada `module.go` compõe o serviço e o handler com `server.New(svc)` e chama `Handler.RegisterRoutes` no roteador Chi.
- As rotas específicas ficam no `RegisterRoutes` do handler do módulo. O health check global em `main.go` é separado dos handlers dos módulos.

---

## 7. Service

Lógica de negócio central. Coordenador entre repositórios e outros serviços.

- Receber repositórios (e outros serviços) via construtor ou parâmetros de função.
- Retornar dados que o handler converte para os DTOs de resposta da API.
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

1. Consulte os handlers existentes em `internal/<modulo>/internal/server` antes de adicionar uma rota.
2. Defina no pacote consumidor apenas os métodos de serviço necessários e injete a interface pelo construtor.
3. Siga os contratos HTTP do módulo e mantenha regras de negócio no serviço.

---
*Última atualização: $(date +%Y-%m-%d)*
