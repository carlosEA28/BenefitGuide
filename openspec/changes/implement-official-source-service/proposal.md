# Proposal

## Why

O Catalog já persiste `OfficialSource`, mas ainda não aplica no service as regras da ADR 0005. Sem essa camada, dados inválidos e alterações fora do `PublicService` informado podem chegar ao repository, e a ausência de serviço não recebe um erro de domínio previsível.

## What Changes

- Acrescentar operações de criação, listagem, consulta, atualização completa e exclusão de fontes oficiais ao service do Catalog.
- Validar a existência do `PublicService` na criação e na listagem, aceitando serviços `ACTIVE` e `INACTIVE`.
- Exigir `Title`, `URL` e `Agency` válidos; remover espaços nas bordas de `Title` e `URL`; aceitar `LastVerifiedAt` ausente e permitir limpá-lo na atualização.
- Restringir consulta individual, atualização e exclusão à fonte pertencente ao `PublicService` informado, com `ErrOfficialSourceNotFound` para ausência ou vínculo divergente.
- Adicionar `ErrInvalidOfficialSource`, preservando os erros de banco não relacionados.

## Capabilities

### New Capabilities

- `catalog/official-sources`: Gerenciamento de fontes oficiais vinculadas a serviços públicos, incluindo validação, consulta e proteção do vínculo.

### Modified Capabilities

Nenhuma. O projeto ainda não possui specs OpenSpec publicadas.

## Impact

- Código Go em `backend/benefit-api/internal/catalog` e erro de domínio em `internal/shared/exceptions`.
- Composição interna do módulo Catalog para injetar `OfficialSourceRepository`; sem rotas HTTP, mudanças de schema ou novas dependências nesta etapa.
- Verificação por formatação, compilação e revisão estática, conforme o pedido de implementar sem testes.
