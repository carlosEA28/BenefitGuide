# Design

## Context

O módulo `catalog` já compõe o repository e o service, mas não tem handler nem registro de rotas. O repository lê requisitos por ID ou serviço e consulta serviços públicos; o service implementa criação, atualização e exclusão, porém ainda não oferece métodos de leitura. `main.go` instancia o módulo e não chama `RegisterRoutes` nele. As ADRs 0002 e 0003 fixam o vínculo imutável e a validação de criação em serviço `ACTIVE`.

## Goals / Non-Goals

**Goals:**

- Manter as regras de domínio no service e a tradução HTTP no handler.
- Garantir que a URL aninhada delimite leitura e mutações, inclusive quando o ID do requisito é válido mas pertence a outro serviço.
- Definir um contrato JSON simples e testável, sem expor erros de banco.

**Non-Goals:**

- Criar autenticação, autorização ou endpoints para `PublicService`.
- Alterar schema, migrations ou regras de negócio já decididas nas ADRs.

## Decisions

### Encapsular operações por serviço no Catalog

Adicionar ao service métodos de listagem e consulta que verificam a existência do `PublicService` e o vínculo do requisito. Para PUT e DELETE, criar operações no service que consultam o requisito e conferem `PublicServiceID` antes de usar a lógica de atualização/exclusão existente. Preservar os métodos atuais para não quebrar consumidores internos. A alternativa de fazer a conferência apenas no handler colocaria uma regra de escopo de dados na camada HTTP.

### Usar DTOs HTTP próprios

O handler recebe apenas `description`, `type`, `condition_description` e `display_order` em `snake_case`; o `serviceId` e o `id` vêm da URL. Em POST, `description` e `type` são obrigatórios; `display_order` pode ser omitido e assume zero. Em PUT, `description`, `type` e `display_order` são obrigatórios para substituição completa; `condition_description` pode ser nulo ou omitido para `MANDATORY`, mas precisa conter texto para `CONDITIONAL`. Decodificar um único objeto JSON, rejeitar campos desconhecidos e conteúdo adicional. Respostas têm campos explícitos de requisito (`id`, `public_service_id`, campos editáveis e timestamps). Usar DTOs evita que o cliente forneça `id`, `public_service_id` ou timestamps, o que a desserialização direta no modelo permitiria.

### Registrar rotas Chi no módulo Catalog

Criar `internal/catalog/internal/server` com um handler que depende de interface pequena do service e registrá-lo em `catalog.Module`. `main.go` chama `catalogModule.RegisterRoutes(router)`, agrupando as cinco rotas sob `/api/v1/services/{serviceId}/requirements`. Reutilizar Chi e a injeção manual já presentes; não adicionar dependência.

### Mapear erros de domínio para HTTP

UUID inválido, JSON malformado, corpo com formato ou campo inesperado e `ErrInvalidRequirement` retornam `400`. `ErrPublicServiceNotFound` e `ErrRequirementNotFound` retornam `404`; requisito de outro serviço é tratado como não encontrado. `ErrPublicServiceInactive` retorna `409`. Outros erros retornam `500` com mensagem genérica. Todas as falhas usam `{"error":"..."}`. Para leitura da lista, serviço existente e vazio gera `[]`, enquanto serviço inexistente gera `404`, exigindo consulta explícita ao serviço público antes da lista. Essa distinção preserva o significado da URL e a semântica já definida na ADR 0002 para lista vazia.

## Risks / Trade-offs

- [Consulta de vínculo seguida de mutação ocorre em duas operações] → O vínculo é imutável pela ADR 0002; uma exclusão concorrente é detectada pelo repository e devolve `404`. Conferir o tratamento de ID de outro serviço na revisão do código.
- [Serviço público pode ser excluído entre consulta e criação/listagem] → A chave estrangeira protege a criação; falha inesperada é tratada sem expor detalhes internos. Não ampliar o escopo para transações ou tradução de códigos PostgreSQL neste change.

## Migration Plan

Publicar os novos handlers junto ao binário da API. Não há migração de dados. O rollback consiste em voltar ao binário anterior; as tabelas e regras de domínio permanecem compatíveis.
