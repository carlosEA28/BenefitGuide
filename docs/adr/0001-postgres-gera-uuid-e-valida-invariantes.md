# Postgres gera o UUID e valida as invariantes do requisito

As tabelas do catálogo geram o próprio `id` com `DEFAULT gen_random_uuid()` (PostgreSQL 16), e a regra de que requisito `CONDITIONAL` exige `condition_description` não nula e não vazia é uma `CHECK` constraint, não validação em Go. A alternativa era gerar o UUID no repository (`uuid.New()`) e validar no service; recusamos porque qualquer caminho de escrita — repository, migration de dados ou `psql` manual — conseguiria gravar um requisito em estado inválido, e porque o repository voltaria a carregar regra de domínio. A consequência aceita é que erros do PostgreSQL (`23514`, `23503`) sobem crus, sem tradução para erros de domínio nesta etapa.

## Consequences

- O `Create` não precisa preencher `id`; o GORM omite o campo (tag `default`) e o driver lê o valor de volta via `RETURNING`. Se a tag `default` e o `DEFAULT` do banco divergirem, o insert falha com `NOT NULL`.
- Erros de violação de constraint são indistinguíveis entre si na camada de service até alguém traduzi-los. Se isso virar problema, a tradução acontece no repository, e não na constraint.
