# Proposal

## Why

O Catalog já persiste requisitos, mas seu service ainda não aplica as regras descritas na ADR 0003. Isso deixa a validação dependente das constraints do banco e não permite distinguir, na camada de negócio, um serviço público ausente ou inativo de uma falha de persistência.

## What Changes

- Acrescentar ao Catalog uma consulta por ID que confirme a existência e o status do `PublicService`. `FindByPublicServiceID` consulta requisitos e não atende a esse uso.
- Permitir a criação de requisitos somente para `PublicService` existente e `ACTIVE`, conforme decisão do responsável do projeto posterior à ADR 0003.
- Validar descrição, tipo e condição de requisitos antes de criar ou atualizar; limpar a condição quando o tipo passar a `MANDATORY`.
- Garantir que atualização e exclusão de requisitos inexistentes retornem `ErrRequirementNotFound` e preservem o vínculo imutável com o serviço.
- Definir erros de domínio distinguíveis para entrada inválida e serviço público ausente ou inativo; manter a proteção do banco.

## Capabilities

### New Capabilities

- `catalog/service-requirements`: Gerenciamento de requisitos de serviços públicos com validação de dados, vínculo com serviço ativo e semântica de ausência.

### Modified Capabilities

Nenhuma. O projeto ainda não possui specs OpenSpec publicadas.

## Impact

- Código Go no módulo `backend/benefit-api/internal/catalog` e erros compartilhados em `internal/shared/exceptions`.
- Contratos internos do service e repository; sem rotas HTTP ou mudanças de schema/migration nesta etapa.
- Testes do service para regras de validação e estados de existência; `make fmt` e `make test` na futura implementação.
