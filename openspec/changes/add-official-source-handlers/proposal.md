# Proposal

## Why

O Catalog já persiste e valida `OfficialSource`, mas ainda não expõe fontes oficiais pela API. Os cinco endpoints permitem cadastrar e consultar as referências de cada `PublicService` e manter suas informações atualizadas.

## What Changes

- Adicionar `POST` e `GET` em `/api/v1/services/{serviceId}/sources` para criar e listar fontes.
- Adicionar `GET`, `PUT` e `DELETE` em `/api/v1/services/{serviceId}/sources/{id}` para consultar, substituir os campos editáveis e excluir uma fonte.
- Validar IDs e corpos JSON na borda HTTP, usar o `serviceId` da rota como vínculo e traduzir erros de domínio em respostas HTTP consistentes.
- Registrar as rotas no handler já usado pelo módulo Catalog, preservando os endpoints de requisitos.

## Capabilities

### New Capabilities

- `catalog/official-sources`: gerenciamento de fontes oficiais por serviço público, agora acessível pelos endpoints HTTP de criação, consulta, atualização e exclusão.

### Modified Capabilities

Nenhuma. Ainda não há specs principais publicadas; o delta anterior de `catalog/official-sources` descreve as regras internas do service.

## Impact

- API Go: `backend/benefit-api/internal/catalog/internal/server` e composição do handler no módulo Catalog.
- API pública: cinco novas operações sob `/api/v1/services/{serviceId}/sources`.
- Sem nova dependência ou migration. A verificação de implementação seguirá a preferência anterior de não criar nem executar testes, usando formatação, compilação e revisão do contrato.
