# Spec Delta

## Purpose

Permitir que clientes da API cadastrem e consultem fontes oficiais de um serviço público e mantenham essas referências atualizadas sem alterar seu vínculo.

## ADDED Requirements

### Requirement: Criação de fonte oficial via HTTP
O sistema SHALL oferecer `POST /api/v1/services/{serviceId}/sources` para criar uma fonte vinculada ao serviço da URL, recebendo JSON com `title`, `url`, `agency` obrigatórios e `last_verified_at` opcional.

#### Scenario: Criação válida
- **WHEN** o cliente envia dados válidos para um `PublicService` existente, `ACTIVE` ou `INACTIVE`
- **THEN** a API devolve `201` com a fonte criada em JSON, incluindo `id` e `public_service_id` igual a `serviceId`

#### Scenario: Serviço inexistente
- **WHEN** a criação referencia um `PublicService` inexistente
- **THEN** a API devolve `404` e não cria a fonte

### Requirement: Listagem de fontes por serviço
O sistema SHALL oferecer `GET /api/v1/services/{serviceId}/sources` e retornar somente as fontes do serviço indicado, ordenadas por `created_at` e `id` em ordem crescente.

#### Scenario: Serviço com fontes
- **WHEN** o cliente lista as fontes de um serviço existente
- **THEN** a API devolve `200` com um array JSON contendo apenas as fontes desse serviço na ordem definida

#### Scenario: Serviço sem fontes
- **WHEN** o cliente lista um serviço existente sem fontes
- **THEN** a API devolve `200` com o array JSON `[]`

#### Scenario: Serviço inexistente
- **WHEN** o cliente lista as fontes de um serviço inexistente
- **THEN** a API devolve `404`

### Requirement: Consulta individual no contexto do serviço
O sistema SHALL oferecer `GET /api/v1/services/{serviceId}/sources/{id}` e retornar uma fonte somente quando ela pertencer ao serviço informado.

#### Scenario: Fonte vinculada ao serviço
- **WHEN** o cliente consulta uma fonte pertencente ao `serviceId`
- **THEN** a API devolve `200` com a fonte em JSON

#### Scenario: Fonte ausente ou de outro serviço
- **WHEN** o ID não existe ou a fonte pertence a outro serviço
- **THEN** a API devolve `404` sem expor dados da fonte de outro serviço

### Requirement: Atualização completa no contexto do serviço
O sistema SHALL oferecer `PUT /api/v1/services/{serviceId}/sources/{id}` para substituir `title`, `url`, `agency` e `last_verified_at` de uma fonte do serviço informado, preservando ID, vínculo e data de criação.

#### Scenario: Atualização válida
- **WHEN** o cliente envia `title`, `url` e `agency` válidos para uma fonte do `serviceId`
- **THEN** a API devolve `200` com a fonte atualizada em JSON e seu vínculo original

#### Scenario: Data de verificação omitida ou nula
- **WHEN** o cliente atualiza uma fonte com `last_verified_at` omitido ou `null`
- **THEN** a data anterior de última verificação é removida

#### Scenario: Fonte ausente ou de outro serviço
- **WHEN** o cliente tenta atualizar uma fonte inexistente ou vinculada a outro serviço
- **THEN** a API devolve `404` e não altera a fonte

### Requirement: Exclusão no contexto do serviço
O sistema SHALL oferecer `DELETE /api/v1/services/{serviceId}/sources/{id}` para remover somente uma fonte pertencente ao serviço informado.

#### Scenario: Exclusão válida
- **WHEN** o cliente exclui uma fonte pertencente ao `serviceId`
- **THEN** a API devolve `204` sem corpo e a fonte deixa de existir

#### Scenario: Fonte ausente ou de outro serviço
- **WHEN** o cliente tenta excluir um ID inexistente, já excluído ou pertencente a outro serviço
- **THEN** a API devolve `404` e não exclui nenhuma outra fonte

### Requirement: Contrato JSON e erros HTTP
O sistema SHALL validar `serviceId` e `id` como UUIDs não nulos, aceitar apenas um objeto JSON com os campos editáveis esperados em `POST` e `PUT`, e devolver falhas como JSON com a chave `error`.

#### Scenario: Corpo ou ID inválido
- **WHEN** o ID da URL é inválido, o corpo não é um único objeto JSON válido, faltam campos obrigatórios ou há campos desconhecidos
- **THEN** a API devolve `400` e não persiste a operação

#### Scenario: Fonte inválida pelo domínio
- **WHEN** `title`, `url` ou `agency` violam as regras de validação do Catalog
- **THEN** a API devolve `400` com erro JSON e não persiste a operação

#### Scenario: Falha inesperada
- **WHEN** uma operação falha por erro interno não classificado como entrada inválida ou registro ausente
- **THEN** a API devolve `500` com erro JSON sem expor detalhes internos
