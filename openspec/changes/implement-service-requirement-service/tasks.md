# Tasks

## 1. Consulta de serviço público no Catalog

- [x] 1.1 Adicionar ao repository do Catalog uma busca de `PublicService` por ID que distinga registro ausente de erro de banco; verificar com teste da consulta que serviço sem requisitos também é encontrado.
- [x] 1.2 Expor a consulta ao service pela composição de `catalog.New`, sem usar `FindByPublicServiceID` como prova de existência; verificar com `go test ./internal/catalog/...` dentro de `backend/benefit-api/`.

## 2. Regras do ServiceRequirementService

- [x] 2.1 Definir erros de domínio distinguíveis para serviço público ausente, serviço inativo e requisito inválido; verificar com testes que `errors.Is` reconhece cada causa.
- [x] 2.2 Implementar criação com validação de descrição, tipo, condição e `PublicService ACTIVE` antes da escrita; verificar com testes de serviço ativo, inativo, inexistente e campos inválidos que o repository não grava nas rejeições.
- [x] 2.3 Implementar atualização a partir do requisito existente, preservando ID e vínculo, limpando a condição na mudança para `MANDATORY` e validando o estado final; verificar com testes de mudança de tipo, tentativa de troca de serviço e ID inexistente.
- [x] 2.4 Implementar exclusão com propagação de `ErrRequirementNotFound` para ID ausente, inclusive após segunda exclusão; verificar com testes de sucesso e ausência.
- [x] 2.5 Atualizar a ADR 0003 para registrar a decisão de exigir `ACTIVE` na criação e a consulta por ID escolhida; verificar que a ADR, o glossário e a spec não se contradizem.

## 3. Verificação integrada

- [x] 3.1 Rodar `make fmt` e `make test` em `backend/benefit-api/` e verificar que ambos terminam sem falhas.
