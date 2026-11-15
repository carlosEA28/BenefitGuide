# Validação de requisitos no service do Catalog

## Contexto

O `RequirementRepository` já criava, consultava, atualizava e excluía requisitos, mas o `Service` do Catalog não aplicava as regras de negócio. O PRD exige que todo requisito pertença a um serviço público existente (RN03) e que os requisitos condicionais expliquem quando se aplicam. As migrations restringem `type` a `MANDATORY` ou `CONDITIONAL` e exigem texto não vazio em `condition_description` para `CONDITIONAL`.

O método `FindByPublicServiceID` do repository retorna requisitos, inclusive uma lista vazia quando o serviço existe sem requisitos; portanto, não serve para confirmar a existência de `PublicService`. A chave estrangeira protege o banco, mas, sozinha, não permite que o service identifique a ausência do serviço público como erro de domínio.

## Decisão

- Consultar `PublicService` por ID no repository do próprio Catalog antes de criar um requisito. A criação só é permitida quando o registro existe e está `ACTIVE`, conforme decisão posterior do responsável do projeto. Serviço inexistente e serviço inativo geram erros de domínio distintos.
- Rejeitar `description` vazia ou composta apenas de espaços. Aceitar somente os tipos `MANDATORY` e `CONDITIONAL`. Para `CONDITIONAL`, exigir `condition_description` não nula e não vazia após remover espaços das bordas. Aplicar as mesmas regras aos dados finais de criação e atualização.
- Ao mudar o tipo de `CONDITIONAL` para `MANDATORY`, limpar `condition_description`.
- Na atualização, preservar `id` e `public_service_id`, conforme ADR 0002. A atualização e a exclusão de IDs inexistentes devem devolver `ErrRequirementNotFound`, já emitido pelo repository; erros de banco distintos devem ser propagados.
- Manter as constraints do PostgreSQL como proteção contra gravações por outros caminhos. Esta decisão complementa a ADR 0001: a regra de `CONDITIONAL` passa também a ser validada no service para produzir um erro de domínio antes da escrita.
- Devolver `ErrInvalidRequirement` para dados inválidos, preservando `ErrRequirementNotFound` para mutações sobre requisitos ausentes. Não definir HTTP nesta etapa.

## Consequências

- O service depende de uma consulta de serviços públicos separada da consulta de requisitos. A verificação de `ACTIVE` ocorre antes da criação; atualizações de requisitos já existentes não alteram seu vínculo.
- A validação no service fornece mensagens previsíveis, enquanto as constraints continuam protegendo contra concorrência e escritas externas. Uma exclusão concorrente do serviço público ainda pode resultar em erro de chave estrangeira; eventual tradução desse erro exige uma decisão separada.
- Nenhuma rota HTTP, migration ou mudança em outros módulos é necessária para esta etapa.
