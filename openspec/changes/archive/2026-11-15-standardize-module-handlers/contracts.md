# Contratos HTTP atuais

Inventário conferido em `internal/{auth,consultations,documents,catalog}/internal/server` e nos respectivos `module.go`, antes da refatoração.

## Health dos módulos

| Módulo | Rota | Resultado |
| --- | --- | --- |
| Auth | `GET /auth/health` | `200`, `Content-Type: application/json`, `{"status":"ok"}`; se `Ready` falhar, `503` e `{"status":"unavailable"}` |
| Consultations | `GET /consultations/health` | Igual a Auth |
| Documents | `GET /documents/health` | Igual a Auth |

Os corpos acima não têm quebra de linha final.

## Catalog

| Recurso | Método e rota | Sucesso |
| --- | --- | --- |
| Requirement | `POST /api/v1/services/{serviceId}/requirements` | `201`, `requirementResponse` |
| Requirement | `GET /api/v1/services/{serviceId}/requirements` | `200`, lista JSON (`[]` quando vazia) |
| Requirement | `GET /api/v1/services/{serviceId}/requirements/{id}` | `200`, `requirementResponse` |
| Requirement | `PUT /api/v1/services/{serviceId}/requirements/{id}` | `200`, `requirementResponse` |
| Requirement | `DELETE /api/v1/services/{serviceId}/requirements/{id}` | `204`, sem corpo |
| OfficialSource | `POST /api/v1/services/{serviceId}/sources` | `201`, `sourceResponse` |
| OfficialSource | `GET /api/v1/services/{serviceId}/sources` | `200`, lista JSON (`[]` quando vazia) |
| OfficialSource | `GET /api/v1/services/{serviceId}/sources/{id}` | `200`, `sourceResponse` |
| OfficialSource | `PUT /api/v1/services/{serviceId}/sources/{id}` | `200`, `sourceResponse` |
| OfficialSource | `DELETE /api/v1/services/{serviceId}/sources/{id}` | `204`, sem corpo |

Erros de parâmetros UUID (inclusive UUID nulo): `400`, `{"error":"invalid serviceId"}` ou `{"error":"invalid id"}`. Corpo inválido em POST/PUT: `400`, `{"error":"invalid request body"}`. Erros de serviço por `errors.Is`: `ErrInvalidRequirement` → `400`, `invalid service requirement`; `ErrInvalidOfficialSource` → `400`, `invalid official source`; `ErrPublicServiceNotFound` → `404`, `public service not found`; `ErrRequirementNotFound` → `404`, `service requirement not found`; `ErrOfficialSourceNotFound` → `404`, `official source not found`; `ErrPublicServiceInactive` → `409`, `public service inactive`. Demais erros → `500`, `internal server error`. Todos os erros usam objeto JSON com chave `error`.

`writeJSON` define `Content-Type: application/json` e acrescenta uma quebra de linha ao corpo. O mapeamento de erro é compartilhado pelos dois recursos. Os decodificadores rejeitam campos desconhecidos e conteúdo após o primeiro objeto JSON; `decodeRequirement` exige `description` e `type` (e `display_order` em PUT); `decodeSource` exige `title`, `url` e `agency`. A validação do domínio permanece no serviço.
