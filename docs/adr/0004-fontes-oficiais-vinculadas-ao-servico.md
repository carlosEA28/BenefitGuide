# Fontes oficiais vinculadas ao serviço público

## Contexto

O PRD exige que as informações de um serviço público sejam fundamentadas em referências identificáveis e que a última verificação da fonte possa ser registrada. A tabela `official_sources` já contém `public_service_id`, `title`, `url`, `agency`, `last_verified_at` e `created_at`. Sua chave estrangeira impede apagar um `PublicService` que ainda possui fontes.

O `OfficialSourceRepository` é a primeira etapa da implementação. Não há nesta etapa service ou rotas HTTP para fontes oficiais.

## Decisão

- Cada `OfficialSource` pertence a um único `PublicService`. O vínculo é definido na criação e não pode ser alterado por `Update`; para mover uma fonte, deve-se excluir e cadastrar outra.
- Fontes podem ser cadastradas e mantidas para serviços `INACTIVE`. O repository não valida o status do serviço; a chave estrangeira garante a existência do registro no banco.
- `Create` deixa o PostgreSQL gerar o UUID e retorna a fonte com o ID preenchido. Violações de constraints são propagadas como erros do banco, seguindo a ADR 0001.
- `FindByID`, `Update` e `Delete` retornam `ErrOfficialSourceNotFound` quando o ID não existe ou é `uuid.Nil`. Uma segunda exclusão também retorna esse erro.
- `FindByPublicServiceID` lista por `created_at ASC, id ASC` e devolve slice vazio quando não há fontes. A lista vazia não comprova que o serviço existe; o chamador deve consultar `PublicService` se precisar dessa distinção.
- `Update` altera somente `title`, `url`, `agency` e `last_verified_at`, inclusive quando a data passa a `NULL`. `id`, `public_service_id` e `created_at` permanecem imutáveis.

## Consequências

- O repository de fontes é uma implementação separada do repository já usado para requisitos, pois ambos têm um método `FindByPublicServiceID` com tipos de retorno diferentes.
- Regras de validação de título, URL e órgão, além da resposta HTTP para os erros, ficam para a futura camada de service e handler. As constraints existentes continuam protegendo a persistência.
