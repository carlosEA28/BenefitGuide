# Validação de fontes oficiais no service do Catalog

**Estado:** implementada no service do Catalog.

## Contexto

O PRD exige fontes identificáveis para fundamentar as informações do catálogo (RN04) e o registro opcional da última verificação. O `OfficialSourceRepository` já oferece criação, busca, listagem, atualização e exclusão. A chave estrangeira garante que `public_service_id` aponte para um registro existente, mas não fornece ao service um erro de domínio previsível antes da criação. Também não valida a forma de `title`, `url` e `agency`.

A ADR 0004 estabelece que a fonte pode pertencer a um serviço `INACTIVE`, que seu vínculo com `PublicService` é imutável e que o repository atualiza apenas `title`, `url`, `agency` e `last_verified_at`. O service de requisitos do mesmo módulo já consulta `PublicService` para distinguir serviço inexistente de lista vazia.

## Decisão

- Na criação, validar a fonte e consultar `PublicService` pelo ID antes de persistir. `ErrPublicServiceNotFound` indica serviço ausente; tanto serviços `ACTIVE` quanto `INACTIVE` podem receber fontes.
- Exigir `Title` e `URL` não vazios após remover espaços das bordas. Persistir esses dois valores sem os espaços das bordas. Exigir também `Agency` não vazia nem composta apenas de espaços.
- Considerar válida uma URL somente quando seu esquema for `http` ou `https` e houver host sintaticamente válido. A validação é de formato, sem consulta externa ao endereço.
- `LastVerifiedAt` é opcional. O service não preenche a data automaticamente; `nil` é aceito na criação e, em uma atualização, limpa o valor anterior.
- A atualização recebe o conjunto completo dos campos editáveis: `Title`, `URL`, `Agency` e `LastVerifiedAt`. Validar os dados finais antes da gravação. Preservar `ID`, `PublicServiceID` e `CreatedAt` do registro atual, conforme ADR 0004.
- Antes de atualizar ou excluir, buscar a fonte por ID e confirmar que seu `PublicServiceID` corresponde ao serviço informado. Se a fonte não existir, o ID for nulo ou o vínculo não corresponder, retornar `ErrOfficialSourceNotFound`. A consulta individual vinculada ao serviço segue a mesma regra.
- Na listagem por serviço, consultar `PublicService` para diferenciar serviço inexistente de serviço existente sem fontes. Para o segundo caso, devolver lista vazia.
- Dados inválidos deverão produzir um erro de domínio identificável, `ErrInvalidOfficialSource`. Erros de banco não relacionados à ausência de registros continuam propagados. Nenhuma semântica HTTP é definida nesta etapa.

## Consequências

- O service de fontes dependerá de `OfficialSourceRepository` e `PublicServiceRepository`. A validação da associação ocorre no service, enquanto a chave estrangeira continua protegendo a persistência.
- O service deve copiar o registro atual antes da atualização para impedir que campos de identidade enviados pelo chamador substituam os dados persistidos.
- Esta decisão não requer alteração de migration. Rotas HTTP permanecem fora desta etapa.
