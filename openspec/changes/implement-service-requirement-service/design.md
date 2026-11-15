# Design

## Context

Ver `proposal.md` e a ADR 0003. O `catalog/internal/service/service.go` só injeta `RequirementRepository`. O repository consulta requisitos por ID e por serviço, mas não consulta a tabela `public_services`. `FindByPublicServiceID` retorna lista vazia para um serviço sem requisitos, então não prova a existência dele. O modelo já define os dois tipos, e as migrations protegem tipo, chave estrangeira e condição de `CONDITIONAL`. A ADR 0002 fixa `id` e `public_service_id` em atualizações e define `ErrRequirementNotFound` para update/delete sem linha.

## Goals / Non-Goals

**Goals:**
- Fazer o service decidir sobre validade dos dados e elegibilidade do serviço associado antes da escrita.
- Preservar a separação entre regras de negócio no service e consultas SQL no repository.
- Usar erros de domínio verificáveis por `errors.Is` sem traduzir erros de infraestrutura genericamente.

**Non-Goals:**
- Criar handlers, rotas HTTP ou DTOs de transporte.
- Alterar migrations ou mover requisitos entre serviços públicos.
- Traduzir todas as violações de constraint produzidas por escritas concorrentes ou externas.

## Decisions

1. **Consulta específica de `PublicService` no Catalog.** Ampliar o repository com uma busca por ID que devolva o registro e distinga ausência de falha do banco. O service verificará `Status == ACTIVE` antes de criar. Alternativa considerada: inferir existência por `FindByPublicServiceID` ou pela chave estrangeira; ambas deixam ambíguo o caso de serviço existente sem requisitos. A decisão do responsável posterior à ADR 0003 exige `ACTIVE` para novos requisitos.
2. **Service de requisitos usando as dependências existentes do módulo.** Acrescentar operações de criar, atualizar e excluir ao service do Catalog, preservando a construção em `catalog.New`. Uma interface pequena de consulta de serviço público pode ser injetada separadamente da `RequirementRepository` para facilitar testes. Alternativa considerada: regras no repository; recusada porque misturaria domínio e persistência.
3. **Validação antes da persistência.** Usar `strings.TrimSpace` para determinar se descrição e condição têm conteúdo. Não reescrever automaticamente o texto informado, exceto limpar `ConditionDescription` ao tornar o requisito `MANDATORY`. Validar o tipo contra as constantes do modelo. Alternativa considerada: depender das tags `validate` e das constraints; elas não fornecem sozinhas o contrato completo de erros do service.
4. **Atualizar a partir do requisito existente.** Buscar por ID antes da atualização, aplicar somente campos mutáveis e validar o estado resultante. Isso preserva `PublicServiceID` e `CreatedAt`, evita devolver um modelo com vínculo falsificado e prioriza `ErrRequirementNotFound` quando o ID não existe. O repository continua verificando `RowsAffected` para cobrir exclusão concorrente. Alternativa considerada: repassar o modelo de entrada diretamente; seu vínculo pode não ser persistido, mas aparecer incorretamente no retorno.
5. **Erros distinguíveis.** Usar `ErrRequirementNotFound` existente e adicionar sentinelas para `PublicService` ausente, serviço inativo e requisito inválido, com contexto de campo quando útil. Alternativa considerada: um erro genérico para tudo; prejudicaria o futuro handler ao mapear a causa. Erros de banco não identificados serão propagados.

## Risks / Trade-offs

- [Serviço excluído ou desativado entre consulta e insert] → A chave estrangeira ainda protege contra exclusão. A garantia de status no instante da escrita exigiria transação ou bloqueio; esta etapa verifica o estado lido antes da criação e documenta a limitação.
- [Dados antigos inválidos falham ao serem editados] → Validar o estado final da atualização e devolver erro explícito para que o cadastro seja corrigido.
- [ADR 0003 ainda registra status como ponto aberto] → A proposta e a spec registram a decisão posterior de exigir `ACTIVE`; alinhar a ADR quando a implementação for autorizada.

## Migration Plan

Não há migration de banco ou dependência externa. A futura implementação deve atualizar a composição do módulo, executar `make fmt` e `make test` dentro de `backend/benefit-api/` e permitir reversão pelo revert das alterações Go.
