# Proposal

## Why

Os handlers dos módulos usam estruturas diferentes: o Catalog já depende de uma interface local, enquanto Auth, Consultations e Documents dependem diretamente de services concretos. O exemplo fornecido estabelece um padrão único para separar a borda HTTP da lógica de serviço e tornar o tratamento de erros explícito.

## What Changes

- Refatorar os handlers de `auth`, `catalog`, `consultations` e `documents` para depender de interfaces pequenas definidas junto a cada handler, receber essas dependências por construtor e manter métodos HTTP focados em entrada, delegação e resposta.
- Aplicar `errors.Is` ao mapeamento de erros de domínio onde houver tais erros, preservando os códigos, mensagens, corpos e rotas HTTP atuais.
- Atualizar `docs/conventions/golang.md` para mostrar o padrão com Chi, DTOs de borda e helpers de JSON equivalentes aos do exemplo, sem usar Gin.
- Preservar o health check global em `main.go` fora desta refatoração.

## Capabilities

### New Capabilities

Nenhuma. A mudança reorganiza handlers existentes sem criar comportamento externo.

### Modified Capabilities

Nenhuma. Os contratos HTTP permanecem os mesmos; `.openspec.yaml` declara `skip_specs: true`.

## Impact

- Código Go em `backend/benefit-api/internal/{auth,catalog,consultations,documents}/internal/server` e, se necessário, a composição dos respectivos `module.go`.
- Convenções de Go em `docs/conventions/golang.md`.
- Sem migração de Chi para Gin, nova dependência, migration de banco ou mudança no health check global.
- A verificação seguirá a preferência anterior de não criar nem executar testes: revisão dos contratos, formatação e compilação.
