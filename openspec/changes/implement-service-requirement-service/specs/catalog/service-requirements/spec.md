# Spec Delta

## Purpose

Gerenciar os requisitos associados a serviços públicos do catálogo, garantindo dados válidos e respostas claras para vínculos ou requisitos indisponíveis.

## ADDED Requirements

### Requirement: Associação a serviço ativo
O sistema SHALL criar um requisito somente quando o serviço público associado existir e estiver `ACTIVE`.

#### Scenario: Serviço ativo
- **WHEN** uma criação válida referencia um serviço público `ACTIVE`
- **THEN** o sistema cria o requisito associado a esse serviço

#### Scenario: Serviço inexistente
- **WHEN** uma criação referencia um serviço público inexistente
- **THEN** o sistema rejeita a criação com erro de serviço público não encontrado e não grava o requisito

#### Scenario: Serviço inativo
- **WHEN** uma criação referencia um serviço público `INACTIVE`
- **THEN** o sistema rejeita a criação com erro de serviço público inativo e não grava o requisito

### Requirement: Descrição do requisito
O sistema SHALL rejeitar a criação ou atualização de requisito cuja descrição seja vazia ou contenha apenas espaços.

#### Scenario: Descrição sem conteúdo
- **WHEN** uma criação ou atualização informa descrição vazia ou composta apenas por espaços
- **THEN** o sistema devolve erro de validação e não persiste a alteração

### Requirement: Tipo do requisito
O sistema SHALL aceitar somente `MANDATORY` e `CONDITIONAL` como tipos de requisito.

#### Scenario: Tipo permitido
- **WHEN** uma criação ou atualização informa `MANDATORY` ou `CONDITIONAL`
- **THEN** o tipo passa pela validação de valores permitidos

#### Scenario: Tipo não permitido
- **WHEN** uma criação ou atualização informa qualquer outro tipo
- **THEN** o sistema devolve erro de validação e não persiste a alteração

### Requirement: Condição de requisito condicional
O sistema SHALL exigir uma descrição de condição não vazia para requisitos `CONDITIONAL` e SHALL limpar a condição quando o tipo mudar para `MANDATORY`.

#### Scenario: Condição ausente
- **WHEN** uma criação ou atualização resulta em tipo `CONDITIONAL` com condição nula, vazia ou composta apenas por espaços
- **THEN** o sistema devolve erro de validação e não persiste a alteração

#### Scenario: Mudança para obrigatório
- **WHEN** um requisito `CONDITIONAL` é atualizado para `MANDATORY`
- **THEN** o requisito salvo não contém descrição de condição

### Requirement: Identidade e vínculo imutáveis
O sistema SHALL preservar o ID e o serviço público associado ao atualizar um requisito.

#### Scenario: Tentativa de trocar o serviço associado
- **WHEN** uma atualização informa outro serviço público para um requisito existente
- **THEN** o requisito permanece associado ao serviço original

### Requirement: Ausência de requisito em mutações
O sistema SHALL informar que o requisito não foi encontrado ao tentar atualizar ou excluir um ID inexistente.

#### Scenario: Atualização de requisito inexistente
- **WHEN** uma atualização usa um ID sem requisito correspondente
- **THEN** o sistema devolve erro de requisito não encontrado e não cria um novo requisito

#### Scenario: Exclusão de requisito inexistente
- **WHEN** uma exclusão usa um ID sem requisito correspondente, inclusive um ID já excluído
- **THEN** o sistema devolve erro de requisito não encontrado
