# Requisito não muda de serviço e ausência é erro de domínio

`UpdateRequirement` grava apenas `description`, `type`, `condition_description` e `display_order`. `id`, `public_service_id` e `created_at` são imutáveis: mover um requisito para outro serviço público não existe no domínio (o vínculo nasce junto com o requisito) e reescrever `created_at` apagaria quando o requisito entrou no catálogo. A alternativa era um `Save` da entidade inteira; recusamos porque o `Save` do GORM aceita `public_service_id` novo e, com PK preenchida, pode recriar uma linha já apagada em vez de falhar.

`UpdateRequirement` e `DeleteRequirement` traduzem `RowsAffected == 0` para `exceptions.ErrRequirementNotFound`, como `FindRequirementByID` já faz com `gorm.ErrRecordNotFound`. Isso difere do ADR 0001, que deixa violações de `CHECK`/`FOREIGN KEY` subirem cruas: lá o erro é de invariante do banco, aqui "não existe" é um estado de domínio que o service precisa distinguir de falha de infraestrutura para responder 404. `DeleteRequirement` não é idempotente por escolha — apagar duas vezes é um erro, não um sucesso silencioso.

`FindByPublicServiceID` devolve `[]models.ServiceRequirement` (não ponteiros) ordenado por `display_order`, `created_at` e `id`, e serviço sem requisitos devolve slice vazio com `nil`, não `ErrRequirementNotFound`: lista vazia é uma resposta válida para "quais requisitos este serviço tem?".

## Consequences

- O update usa `Model(...).Select("description", "type", "condition_description", "display_order")`, então o GORM acrescenta `updated_at` sozinho e atualiza o campo no struct devolvido; `created_at` não é tocado.
- `uuid.Nil` é rejeitado antes de qualquer query com `ErrRequirementNotFound`. Sem essa guarda o GORM emite `UPDATE`/`DELETE` sem `WHERE` e atinge a tabela inteira, porque o `WHERE` do update vem da primary key.
- Quem quiser mover um requisito de serviço precisa apagar e criar de novo, o que expõe a operação na auditoria.
- Um futuro soft delete (`deleted_at`) muda a semântica de `DeleteRequirement` e de `FindByPublicServiceID`; esta ADR terá de ser revista junto.
