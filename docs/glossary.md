# Glossário do domínio

## Catalog

Módulo responsável pelo catálogo de serviços públicos, seus requisitos de acesso e fontes oficiais.

## PublicService (serviço público)

Registro de um serviço ou benefício apresentado pelo catálogo. Possui um identificador interno e um status (`ACTIVE` ou `INACTIVE`). A existência do registro é distinta da sua disponibilidade para apresentação ao cidadão. Novos requisitos só podem ser associados a serviços `ACTIVE`.

## ServiceRequirement (requisito de serviço)

Condição ou obrigação descrita para acessar um `PublicService`. Pertence a um único serviço público. Seu vínculo, representado por `public_service_id`, é definido na criação e não pode ser alterado por atualização (ADR 0002).

## MANDATORY (obrigatório)

Tipo de requisito aplicável a todos os casos descritos para o serviço. Não exige `condition_description`; ao mudar um requisito para esse tipo, o texto da condição é limpo.

## CONDITIONAL (condicional)

Tipo de requisito aplicável somente quando uma condição é satisfeita. Exige `condition_description` com texto significativo.

## Description (descrição do requisito)

Texto que informa qual requisito deve ser atendido. Uma string vazia ou composta apenas de espaços não constitui uma descrição válida.

## ConditionDescription (descrição da condição)

Texto que explica quando um requisito `CONDITIONAL` se aplica. Ao mudar para `MANDATORY`, o campo é limpo.

## Requisito inexistente

Requisito cujo ID não corresponde a um registro no catálogo. Tentativas de atualização ou exclusão geram `ErrRequirementNotFound` (ADR 0002).
