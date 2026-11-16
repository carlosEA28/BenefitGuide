# Tasks

## 1. Contratos atuais

- [x] 1.1 Inventariar rotas, códigos, mensagens e corpos de resposta dos handlers de Auth, Catalog (requirements e sources), Consultations e Documents; verificar o inventário contra os handlers e `module.go` atuais.

## 2. Handlers de Auth, Consultations e Documents

- [x] 2.1 Refatorar o handler de Auth para receber uma interface local mínima pelo construtor, mantendo o comportamento de `/health`; verificar assinatura, chamada a `Ready` e respostas por revisão do diff.
- [x] 2.2 Refatorar o handler de Consultations pelo mesmo padrão, mantendo o comportamento de `/health`; verificar assinatura, chamada a `Ready` e respostas por revisão do diff.
- [x] 2.3 Refatorar o handler de Documents pelo mesmo padrão, mantendo o comportamento de `/health`; verificar assinatura, chamada a `Ready` e respostas por revisão do diff.

## 3. Handlers de Catalog

- [x] 3.1 Ajustar interfaces locais, construtores e métodos HTTP dos handlers de requirements para o padrão de entrada, serviço, erro e resposta; verificar que rotas e respostas coincidem com o inventário.
- [x] 3.2 Ajustar os handlers de OfficialSource, preservando as alterações já presentes e o vínculo com serviço nas rotas; verificar que rotas, validação de entrada e respostas coincidem com o inventário.

## 4. Convenção Go

- [x] 4.1 Atualizar `docs/conventions/golang.md` com exemplo Chi de interface local, injeção por construtor, método HTTP e `errors.Is`; verificar coerência dos exemplos e orientações com os handlers refatorados e as demais seções do documento.

## 5. Integração

- [x] 5.1 Conferir as composições nos `module.go`, formatar o Go alterado e executar `go build ./...` no backend; verificar compilação e ausência de mudanças no health check global ou nos contratos inventariados, sem criar nem executar testes.
