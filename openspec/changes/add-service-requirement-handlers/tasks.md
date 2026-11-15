# Tasks

## 1. Operações do Catalog por serviço público

- [x] 1.1 Adicionar ao service a listagem de requisitos com checagem de existência do `PublicService`; verificar por revisão do fluxo que a ordenação existente, a lista vazia e o serviço inexistente são tratados.
- [x] 1.2 Adicionar consulta de requisito com verificação de `PublicServiceID` e operações de atualização/exclusão delimitadas pelo `serviceId`, preservando os métodos atuais; verificar por revisão que IDs de outro serviço e IDs ausentes retornam não encontrado sem mutação.

## 2. Contrato HTTP e handlers

- [x] 2.1 Criar DTOs de entrada e saída e helpers de decodificação/resposta JSON, com validação de UUID, campos obrigatórios, campos desconhecidos e corpo adicional; verificar por revisão dos caminhos de entrada e saída.
- [x] 2.2 Implementar POST, GET lista e GET por ID com tradução de erros de domínio para `400`, `404`, `409` e `500`; verificar por revisão dos status, corpos, ordenação, serviço vazio/inexistente e ausência de detalhes internos.
- [x] 2.3 Implementar PUT e DELETE com escopo por `serviceId`, resposta `200`/`204` e erros JSON; verificar por revisão a substituição completa, os campos imutáveis, o requisito de outro serviço e o ID ausente.

## 3. Registro das rotas

- [x] 3.1 Montar o handler no `catalog.Module` e registrar as cinco rotas em `main.go`; verificar por revisão das rotas e pela compilação do módulo.
- [x] 3.2 Executar `make fmt` e `make build` em `backend/benefit-api/` e verificar que a formatação e a compilação Go passam, sem executar testes.
