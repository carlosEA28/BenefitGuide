# Glossário do domínio

## Catalog

Módulo responsável pelo catálogo de serviços públicos, seus requisitos de acesso e fontes oficiais.

## PublicService (serviço público)

Registro de um serviço ou benefício apresentado pelo catálogo. Possui um identificador interno e um status (`ACTIVE` ou `INACTIVE`). A existência do registro é distinta da sua disponibilidade para apresentação ao cidadão. Novos requisitos só podem ser associados a serviços `ACTIVE`.

## ServiceRequirement (requisito de serviço)

Condição ou obrigação descrita para acessar um `PublicService`. Pertence a um único serviço público. Seu vínculo, representado por `public_service_id`, é definido na criação e não pode ser alterado por atualização (ADR 0002).

## OfficialSource (fonte oficial)

Referência identificável a uma página, documento ou publicação governamental que fundamenta as informações de um `PublicService`. Registra título, URL, órgão responsável e, quando conhecida, a data da última verificação. Pertence a um único serviço público existente; o vínculo não muda após o cadastro (ADR 0004). Pode ser cadastrada e mantida mesmo quando o serviço está `INACTIVE`. O service valida os dados e o vínculo conforme a ADR 0005.

## Title (título da fonte oficial)

Nome que identifica a `OfficialSource`. É obrigatório e não pode ser vazio ou composto apenas de espaços. Espaços nas bordas são removidos antes de validar e cadastrar ou atualizar.

## URL da fonte oficial

Endereço da referência. É obrigatório e deve ter esquema `http` ou `https` e host sintaticamente válido. Espaços nas bordas são removidos antes da validação e da persistência (ADR 0005).

## Agency (órgão da fonte oficial)

Órgão responsável indicado na referência. É obrigatório e não pode ser vazio ou composto apenas de espaços (ADR 0005).

## LastVerifiedAt (última verificação da fonte)

Data e hora em que a `OfficialSource` foi verificada pela última vez. É opcional; valor ausente permanece `NULL` e uma atualização pode limpar o valor anterior (ADR 0005).

## Fonte oficial inexistente

Fonte cujo ID não corresponde a um registro no catálogo. Consulta por ID, atualização e exclusão retornam `ErrOfficialSourceNotFound` (ADR 0004). Para operações vinculadas a um `PublicService`, uma fonte que pertence a outro serviço recebe o mesmo erro (ADR 0005).

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
