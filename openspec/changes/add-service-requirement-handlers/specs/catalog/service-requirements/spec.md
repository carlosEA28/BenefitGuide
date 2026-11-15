# Spec Delta

## Purpose

Permitir o gerenciamento dos requisitos de um serviço público pela API, preservando o vínculo de cada requisito e as validações do Catalog.

## ADDED Requirements

### Requirement: Criação de requisito via HTTP
O sistema SHALL oferecer `POST /api/v1/services/{serviceId}/requirements` para criar um requisito associado ao serviço público indicado na URL, recebendo JSON com `description`, `type`, `condition_description` e `display_order`.

#### Scenario: Criação em serviço ativo
- **WHEN** o cliente envia dados válidos para um serviço público `ACTIVE`
- **THEN** a API devolve `201` com o requisito criado em JSON, incluindo seu ID e `public_service_id` igual a `serviceId`

#### Scenario: Serviço inexistente ou inativo
- **WHEN** o cliente tenta criar em um serviço público inexistente ou `INACTIVE`
- **THEN** a API devolve `404` para o inexistente ou `409` para o inativo e não cria requisito

### Requirement: Listagem de requisitos do serviço
O sistema SHALL oferecer `GET /api/v1/services/{serviceId}/requirements` e retornar somente requisitos vinculados ao serviço indicado, ordenados por `display_order`, `created_at` e `id`, nessa ordem crescente.

#### Scenario: Serviço existente sem requisitos
- **WHEN** o cliente lista requisitos de um serviço público existente sem requisitos
- **THEN** a API devolve `200` com o array JSON `[]`

#### Scenario: Serviço inexistente
- **WHEN** o cliente lista requisitos de um serviço público inexistente
- **THEN** a API devolve `404`

### Requirement: Consulta de requisito no contexto do serviço
O sistema SHALL oferecer `GET /api/v1/services/{serviceId}/requirements/{id}` e retornar o requisito somente se ele pertencer ao serviço indicado.

#### Scenario: Requisito vinculado ao serviço
- **WHEN** o cliente consulta um ID de requisito que pertence ao `serviceId`
- **THEN** a API devolve `200` com o requisito em JSON

#### Scenario: Requisito ausente ou de outro serviço
- **WHEN** o ID do requisito não existe ou pertence a outro serviço público
- **THEN** a API devolve `404` sem revelar os dados do requisito de outro serviço

### Requirement: Atualização de requisito no contexto do serviço
O sistema SHALL oferecer `PUT /api/v1/services/{serviceId}/requirements/{id}` para substituir os campos editáveis `description`, `type`, `condition_description` e `display_order`, preservando ID, serviço associado e data de criação.

#### Scenario: Atualização válida
- **WHEN** o cliente envia os campos editáveis válidos para um requisito vinculado ao `serviceId`
- **THEN** a API devolve `200` com o requisito atualizado em JSON e preserva seu vínculo

#### Scenario: Requisito ausente ou de outro serviço
- **WHEN** o cliente tenta atualizar um ID inexistente ou vinculado a outro serviço público
- **THEN** a API devolve `404` e não altera nenhum requisito

### Requirement: Exclusão de requisito no contexto do serviço
O sistema SHALL oferecer `DELETE /api/v1/services/{serviceId}/requirements/{id}` para remover somente um requisito pertencente ao serviço indicado.

#### Scenario: Exclusão válida
- **WHEN** o cliente exclui um requisito vinculado ao `serviceId`
- **THEN** a API devolve `204` sem corpo e o requisito deixa de existir

#### Scenario: Requisito ausente ou de outro serviço
- **WHEN** o cliente tenta excluir um ID inexistente ou vinculado a outro serviço público
- **THEN** a API devolve `404` e não exclui nenhum requisito

### Requirement: Validação e erros HTTP
O sistema SHALL validar os IDs da rota como UUIDs e os corpos de `POST` e `PUT` como objetos JSON com os campos editáveis esperados. Erros SHALL ser devolvidos em JSON com a chave `error`.

#### Scenario: Entrada inválida
- **WHEN** um ID é inválido, o JSON é malformado, há campos desconhecidos ou os dados do requisito violam as regras de domínio
- **THEN** a API devolve `400` com erro JSON e não persiste a operação

#### Scenario: Falha inesperada
- **WHEN** uma operação falha por erro interno não classificado como validação, ausência ou conflito de estado
- **THEN** a API devolve `500` com erro JSON sem expor detalhes internos
