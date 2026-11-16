# Spec Delta

## Purpose

Gerenciar as fontes oficiais que fundamentam informações de serviços públicos, garantindo referências válidas e operações limitadas ao serviço ao qual cada fonte pertence.

## ADDED Requirements

### Requirement: Associação a serviço público existente
O sistema SHALL criar uma fonte oficial somente quando o `PublicService` associado existir, independentemente de seu status `ACTIVE` ou `INACTIVE`.

#### Scenario: Serviço existente e ativo
- **WHEN** uma fonte válida é criada para um serviço `ACTIVE`
- **THEN** a fonte é associada ao serviço

#### Scenario: Serviço existente e inativo
- **WHEN** uma fonte válida é criada para um serviço `INACTIVE`
- **THEN** a fonte é associada ao serviço

#### Scenario: Serviço inexistente
- **WHEN** a criação informa um `PublicService` inexistente
- **THEN** o sistema devolve erro de serviço público não encontrado e não grava a fonte

### Requirement: Campos de identificação da fonte
O sistema SHALL exigir `Title`, `URL` e `Agency` com conteúdo na criação e na atualização completa. SHALL remover espaços das bordas de `Title` e `URL` antes de validar e persistir esses campos.

#### Scenario: Título ou URL contém espaços nas bordas
- **WHEN** uma criação ou atualização informa `Title` ou `URL` válidos com espaços nas bordas
- **THEN** o sistema persiste esses valores sem os espaços das bordas

#### Scenario: Campo obrigatório vazio
- **WHEN** `Title`, `URL` ou `Agency` é vazio ou contém apenas espaços
- **THEN** o sistema devolve erro de fonte oficial inválida e não grava a alteração

### Requirement: Formato da URL
O sistema SHALL aceitar apenas URLs com esquema `http` ou `https` e host sintaticamente válido.

#### Scenario: URL válida
- **WHEN** a fonte informa uma URL HTTP ou HTTPS com host válido
- **THEN** a URL passa pela validação de formato

#### Scenario: URL inválida
- **WHEN** a fonte informa URL sem host, com esquema diferente de HTTP ou HTTPS, ou com formato inválido
- **THEN** o sistema devolve erro de fonte oficial inválida e não grava a alteração

### Requirement: Última verificação opcional
O sistema SHALL aceitar `LastVerifiedAt` ausente e SHALL permitir limpar esse campo em uma atualização completa.

#### Scenario: Criação sem verificação
- **WHEN** uma fonte válida é criada sem `LastVerifiedAt`
- **THEN** a fonte é salva sem data de última verificação

#### Scenario: Remoção da data anterior
- **WHEN** uma fonte com `LastVerifiedAt` preenchido é atualizada com valor ausente nesse campo
- **THEN** a fonte passa a não ter data de última verificação

### Requirement: Consulta de fontes por serviço
O sistema SHALL listar as fontes de um `PublicService` existente e SHALL restringir a consulta individual às fontes pertencentes ao serviço informado.

#### Scenario: Serviço existente sem fontes
- **WHEN** uma listagem é solicitada para serviço existente sem fontes
- **THEN** o sistema devolve uma lista vazia

#### Scenario: Serviço inexistente na listagem
- **WHEN** uma listagem é solicitada para serviço inexistente
- **THEN** o sistema devolve erro de serviço público não encontrado

#### Scenario: Consulta individual fora do serviço
- **WHEN** uma consulta individual informa uma fonte pertencente a outro serviço
- **THEN** o sistema devolve erro de fonte oficial não encontrada

### Requirement: Mutações restritas ao serviço proprietário
O sistema SHALL atualizar ou excluir uma fonte somente quando ela pertencer ao `PublicService` informado. A atualização SHALL substituir todos os campos editáveis e preservar `ID`, `PublicServiceID` e `CreatedAt`.

#### Scenario: Atualização da fonte do serviço
- **WHEN** uma atualização válida informa o serviço proprietário e uma fonte existente
- **THEN** o sistema substitui `Title`, `URL`, `Agency` e `LastVerifiedAt` e preserva a identidade e o vínculo da fonte

#### Scenario: Alteração ou exclusão fora do serviço
- **WHEN** uma atualização ou exclusão informa um serviço diferente do proprietário da fonte
- **THEN** o sistema devolve erro de fonte oficial não encontrada e não altera a fonte

#### Scenario: Fonte inexistente ou ID nulo
- **WHEN** uma atualização ou exclusão informa fonte inexistente, já excluída ou ID nulo
- **THEN** o sistema devolve erro de fonte oficial não encontrada
