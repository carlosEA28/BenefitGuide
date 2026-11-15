# Proposal

## Why

O Catalog já persiste e valida requisitos de serviços públicos, mas não oferece rotas HTTP para gerenciá-los. Expor esses recursos permite cadastrar e consultar requisitos pela API mantendo as regras de domínio existentes.

## What Changes

- Adicionar os cinco endpoints de requisitos sob `/api/v1/services/{serviceId}/requirements`, para criar, listar, consultar, atualizar e excluir.
- Validar parâmetros e corpos JSON na borda HTTP, devolver respostas e erros consistentes e impedir acesso a um requisito por meio do ID de outro serviço público.
- Ligar os handlers ao `Service` do Catalog e registrar as rotas no roteador principal.
- Conferir o contrato HTTP por revisão do código e compilação, sem criar ou executar testes nesta aplicação.

## Capabilities

### New Capabilities

- `catalog/service-requirements`: gerenciamento HTTP dos requisitos associados a um serviço público.

### Modified Capabilities

Nenhuma. Ainda não há especificações principais publicadas neste projeto; o delta anterior dessa capacidade descreve as regras internas do service.

## Impact

- API Go: módulo Catalog, handlers HTTP, registro de rotas em `main.go` e consultas do service necessárias para leitura.
- API pública: cinco novas operações em `/api/v1/services/{serviceId}/requirements`.
- Sem nova dependência, migration ou alteração em outros módulos.
