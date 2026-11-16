# Tasks

## 1. Contrato JSON e erros

- [x] 1.1 Definir DTOs de entrada e saída para `OfficialSource`, com `title`, `url` e `agency` obrigatórios, `last_verified_at` opcional e campos de identidade apenas na resposta; verificar por revisão das tags JSON e compilação.
- [x] 1.2 Decodificar exatamente um objeto JSON, rejeitar campos desconhecidos e formatos inválidos de `last_verified_at`; verificar por revisão dos caminhos de erro e compilação.
- [x] 1.3 Mapear `ErrInvalidOfficialSource` para `400` e `ErrOfficialSourceNotFound` para `404` no tratamento de erros HTTP existente; verificar por revisão do mapeamento e preservar as respostas atuais de requisitos.

## 2. Criação e listagem

- [x] 2.1 Registrar `POST` e `GET` em `/api/v1/services/{serviceId}/sources` no handler do Catalog, usando o service já composto; verificar por revisão do registro de rotas e compilação.
- [x] 2.2 Implementar `POST` usando `serviceId` da rota como vínculo, DTO validado e resposta `201` com fonte criada; verificar por revisão de ID inválido, corpo inválido, serviço ausente e sucesso.
- [x] 2.3 Implementar `GET` de coleção com resposta `200`, ordem preservada do service e array `[]` para serviço existente sem fontes; verificar por revisão de lista vazia e serviço ausente.

## 3. Consulta, atualização e exclusão

- [x] 3.1 Registrar `GET`, `PUT` e `DELETE` em `/api/v1/services/{serviceId}/sources/{id}`; verificar por revisão das cinco rotas registradas e ausência de conflito com rotas de requisitos.
- [x] 3.2 Implementar consulta individual usando `serviceId` e `id` validados, devolvendo `200` ou `404` conforme o service; verificar por revisão do fluxo e compilação.
- [x] 3.3 Implementar `PUT` como substituição completa dos campos editáveis, tratando `last_verified_at` omitido ou `null` como limpeza, e devolver `200` com a fonte atualizada; verificar por revisão de campos obrigatórios, vínculo e resposta.
- [x] 3.4 Implementar `DELETE` com `204` sem corpo para fonte do serviço e `404` para fonte ausente ou de outro serviço; verificar por revisão do fluxo e compilação.

## 4. Verificação integrada

- [x] 4.1 Conferir que proposta, spec, design e handlers descrevem o mesmo contrato HTTP; verificar por revisão e `git diff --check`.
- [x] 4.2 Rodar `make fmt` e `go build ./...` em `backend/benefit-api/` e verificar que ambos terminam sem falhas, sem criar nem executar testes conforme a preferência anterior do usuário.
