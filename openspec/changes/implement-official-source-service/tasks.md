# Tasks

## 1. Dependências e erros do Catalog

- [x] 1.1 Adicionar `ErrInvalidOfficialSource` em `internal/shared/exceptions` e verificar por inspeção que segue o padrão das sentinelas existentes e pode ser identificado por `errors.Is`.
- [x] 1.2 Injetar `OfficialSourceRepository` no `Service` do Catalog e compô-lo em `catalog.New`, preservando chamadas existentes do construtor; verificar por compilação de `backend/benefit-api/`.

## 2. Criação e validação da fonte

- [x] 2.1 Implementar validação de entrada nula, `Title`, `URL` e `Agency`, incluindo normalização de `Title` e `URL`, esquema HTTP/HTTPS e host válido; verificar por revisão das condições de rejeição e compilação.
- [x] 2.2 Implementar criação com consulta ao `PublicService`, aceitando `ACTIVE` e `INACTIVE` e `LastVerifiedAt` ausente; verificar por revisão do fluxo e compilação que só uma fonte válida de serviço existente chega ao repository.

## 3. Consultas e mutações no escopo do serviço

- [x] 3.1 Implementar listagem com consulta prévia ao `PublicService` e consulta individual limitada ao vínculo da fonte; verificar por revisão dos caminhos de serviço inexistente, lista vazia e vínculo divergente.
- [x] 3.2 Implementar atualização completa a partir da fonte persistida, preservando `ID`, `PublicServiceID` e `CreatedAt` e permitindo limpar `LastVerifiedAt`; verificar por revisão dos campos atribuídos e da ordem entre checagem do vínculo, validação e escrita.
- [x] 3.3 Implementar exclusão apenas para fonte do serviço informado; verificar por revisão dos caminhos de ID nulo, fonte ausente, vínculo divergente e propagação de erros do repository.

## 4. Verificação integrada

- [x] 4.1 Conferir que implementação, ADR 0005, glossário e spec descrevem as mesmas regras; verificar por revisão dos artefatos e `git diff --check`.
- [x] 4.2 Rodar `make fmt` e `go build ./...` em `backend/benefit-api/` e verificar que ambos terminam sem falhas, sem criar nem executar testes.
