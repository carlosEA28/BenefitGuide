# Design

## Context

Ver `proposal.md`, a spec `catalog/official-sources` e as ADRs 0004 e 0005. O `Service` do Catalog já oferece `CreateOfficialSource`, `ListOfficialSources`, `GetOfficialSource`, `UpdateOfficialSource` e `DeleteOfficialSource`, com validação e conferência de vínculo. O handler atual expõe apenas requisitos sob `/api/v1/services/{serviceId}/requirements`; `main.go` já registra as rotas do módulo Catalog.

## Goals / Non-Goals

**Goals:**
- Reutilizar o padrão HTTP do Catalog sem deslocar validação de domínio para o handler.
- Fazer a URL aninhada determinar o serviço da fonte e impedir que campos de identidade no corpo alterem esse vínculo.
- Devolver respostas JSON explícitas e erros estáveis para os cinco endpoints.

**Non-Goals:**
- Criar autenticação, autorização ou endpoints para cadastrar `PublicService`.
- Alterar regras do service, repository, migrations ou rotas de requisitos.
- Consultar a disponibilidade da URL oficial.

## Decisions

### Acrescentar as rotas ao handler existente

Ampliar a interface interna usada pelo `Handler` com as cinco operações de fontes e registrar as rotas `sources` em `RegisterRoutes`. O `Service` já satisfaz esse contrato, e `main.go` já chama `catalogModule.RegisterRoutes(router)`. Criar outro handler ou registrar rotas diretamente no `main.go` duplicaria a composição do módulo.

### Usar DTOs próprios para entrada e saída

O corpo de `POST` e `PUT` aceita apenas `title`, `url`, `agency` e `last_verified_at`; os três primeiros são obrigatórios e `last_verified_at` pode ser omitido ou `null`. Em `PUT`, omissão e `null` limpam a data anterior porque a operação é substituição completa. `serviceId` e `id` vêm somente da rota. O DTO de resposta contém `id`, `public_service_id`, `title`, `url`, `agency`, `last_verified_at` e `created_at`; a data opcional aparece como `null` quando ausente. Decodificar um único objeto JSON, rejeitando campos desconhecidos ou dados adicionais. A desserialização direta no modelo aceitaria campos de identidade enviados pelo cliente.

### Reutilizar helpers HTTP e mapear erros

Usar `parseUUID` para IDs não nulos e `writeJSON` para as respostas. Reutilizar o mapeamento central de erros do handler, acrescentando `ErrInvalidOfficialSource` como `400` e `ErrOfficialSourceNotFound` como `404`. `ErrPublicServiceNotFound` já produz `404`; outros erros seguem como `500` genérico. `POST` responde `201`, leituras e `PUT` respondem `200`, e `DELETE` responde `204` sem corpo. A lista vazia é serializada como `[]`, conforme o contrato do service e do repository.

## Risks / Trade-offs

- [Campo `last_verified_at` omitido em `PUT` apaga a data existente] → Documentar a semântica de substituição completa e mapear omissão para `nil`, conforme ADR 0005.
- [Cliente tenta enviar `public_service_id` ou `id` no JSON] → Rejeitar como campo desconhecido; IDs são obtidos da rota e o service mantém o vínculo original.
- [Erro de banco contém detalhes internos] → Responder `500` com mensagem genérica, seguindo o handler atual.

## Migration Plan

Não há migration de banco nem dependência adicional. A implementação registra as novas rotas no binário existente e pode ser revertida retirando os handlers e DTOs de fontes. Conforme a preferência anterior do usuário, verificar por revisão estática, `make fmt` e `go build ./...`, sem criar ou executar testes nesta mudança.
