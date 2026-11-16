# Design

## Context

Ver `proposal.md`, a spec `catalog/official-sources` e as ADRs 0004 e 0005. `OfficialSourceRepository` já oferece as cinco operações de persistência, mas o `Service` do Catalog injeta apenas `RequirementRepository` e `PublicServiceRepository`. `ErrOfficialSourceNotFound` e `ErrPublicServiceNotFound` já existem. O modelo usa `LastVerifiedAt *time.Time`; a coluna correspondente aceita `NULL`. O repository de fontes atualiza apenas os quatro campos editáveis, sem trocar `PublicServiceID`.

## Goals / Non-Goals

**Goals:**
- Concentrar validação e conferência de vínculo no service, com erros de domínio distinguíveis por `errors.Is`.
- Preservar a composição atual do Catalog e facilitar testes com stubs das interfaces de repository.
- Manter os dados de identidade da fonte lidos do banco durante a atualização.

**Non-Goals:**
- Criar handlers, rotas, DTOs HTTP ou regras de autorização de usuário.
- Alterar schema, migrations, o repository de fontes ou a política de status de `ServiceRequirement`.
- Verificar se a URL responde ou se seu conteúdo é governamental.

## Decisions

1. **Injetar o repository de fontes no `Service` existente.** Ampliar o construtor do service para receber `OfficialSourceRepository` e criar essa dependência em `catalog.New`. Preservar compatibilidade com as chamadas existentes do construtor para não alterar arquivos de teste nesta etapa. Isso segue a composição manual do módulo; um service paralelo duplicaria a dependência de `PublicServiceRepository` e a superfície do Catalog.
2. **Validar dados antes da escrita.** Criar `ErrInvalidOfficialSource` em `internal/shared/exceptions` e uma validação do service para entrada nula, título, URL e órgão. Aplicar `strings.TrimSpace` e persistir os valores normalizados de `Title` e `URL`; para `Agency`, verificar que `TrimSpace` não é vazio. Usar `net/url` para analisar a URL e aceitar somente `http`/`https` com hostname sintaticamente válido. Não depender apenas das tags `validate`, pois elas não executam sozinhas nem padronizam o erro do service.
3. **Consultar o serviço público quando sua existência é parte da resposta.** Na criação, buscar `PublicService` antes de chamar `Create`, sem exigir `ACTIVE` (ADR 0004). Na listagem, buscar o serviço antes de listar para distinguir serviço inexistente de serviço sem fontes. A alternativa de inferir existência de uma lista vazia é ambígua.
4. **Consultar a fonte sob o escopo informado antes de ler ou alterar.** Buscar por ID e comparar `PublicServiceID` com o serviço recebido; divergência retorna `ErrOfficialSourceNotFound`. Atualização e exclusão só chamam o repository após essa checagem. A alternativa de confiar apenas no ID permitiria alterar ou excluir uma fonte de outro serviço.
5. **Atualizar uma cópia da fonte persistida.** Aplicar os quatro campos editáveis da entrada à cópia, validar o estado final e enviá-la ao repository. Assim, `ID`, `PublicServiceID` e `CreatedAt` retornados refletem o registro real. Um `LastVerifiedAt` nulo na entrada limpa a data; a implementação do repository já seleciona esse campo explicitamente. Repassar a entrada inteira ao repository poderia devolver um vínculo informado pelo chamador que não foi persistido.
6. **Preservar semântica de erro.** Fonte ausente, ID nulo ou vínculo divergente produzem `ErrOfficialSourceNotFound`; fonte inválida produz `ErrInvalidOfficialSource`; serviço ausente na criação ou listagem produz `ErrPublicServiceNotFound`. Outros erros do repository são propagados. Nenhum mapeamento HTTP será introduzido.

## Risks / Trade-offs

- [Serviço excluído entre consulta e inserção] → A chave estrangeira continua protegendo a escrita; o erro de concorrência do banco é propagado, como na ADR 0001.
- [Fonte excluída entre checagem do vínculo e mutação] → `Update` e `Delete` do repository verificam `RowsAffected` e devolvem `ErrOfficialSourceNotFound`.
- [Validação de URL sintática aceita um endereço inacessível] → O contrato exige formato válido; verificação de disponibilidade fica fora do service nesta etapa.

## Migration Plan

Não há migration de banco nem dependência externa. Na implementação, compor o novo repository no módulo, executar `make fmt` e compilar com `go build ./...` dentro de `backend/benefit-api/`. Não criar nem executar testes nesta etapa, conforme pedido do usuário. A reversão consiste em retirar as operações e a injeção adicionadas ao service e ao módulo.
