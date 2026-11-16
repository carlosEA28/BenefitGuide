# Design

## Context

Ver `proposal.md` para a motivação. A API usa Chi e `net/http`. O Catalog já define uma interface local para o serviço; Auth, Consultations e Documents recebem `*service.Service` diretamente e expõem apenas uma rota de health por módulo. Os handlers de OfficialSource no Catalog estão presentes no estado atual da árvore de trabalho e fazem parte do escopo. O health check global em `main.go` fica fora dele.

## Goals / Non-Goals

**Goals:**

- Deixar explícito, em cada pacote `server`, o conjunto mínimo de métodos de serviço usado pelo handler.
- Tornar uniforme a sequência de cada método HTTP: ler e validar entrada, chamar o serviço, mapear erros conhecidos e escrever a resposta.
- Fazer a convenção documentada refletir a arquitetura em uso e os contratos existentes.

**Non-Goals:**

- Alterar serviços, repositórios, DTOs públicos, rotas ou semântica HTTP.
- Migrar de Chi para Gin ou introduzir dependências.
- Refatorar o health check global.

## Decisions

### Interfaces locais e construtores

Cada pacote `server` define a menor interface necessária ao seu handler. Nos três módulos com rota de health, ela contém `Ready(context.Context) error`; no Catalog, abrange somente as operações de serviço chamadas por seus handlers, incluindo requirements e sources. Os construtores aceitam a interface e armazenam a dependência no handler. Manter os nomes exportados necessários à composição em `module.go`; o exemplo `userHandler/newUserHandler` orienta a estrutura, não exige nomes privados que impediriam o uso entre pacotes.

Alternativa considerada: interface compartilhada em `service` ou tipos concretos. Uma interface local expressa melhor o contrato de consumo e evita acoplar a borda HTTP à implementação.

### Adaptação do exemplo para Chi

Os métodos continuam com assinatura `func(http.ResponseWriter, *http.Request)`. Parâmetros de rota vêm de `chi.URLParam`, o corpo continua sendo decodificado pelos helpers atuais de `encoding/json`, e as respostas continuam usando os escritores HTTP existentes. O tratamento de erros de domínio usa `errors.Is` com os erros já definidos; falhas sem erro de domínio conhecido seguem o mapeamento atual. As rotas e suas respostas permanecem iguais, inclusive os health checks dos módulos.

Alternativa considerada: copiar `gin.Context`, `ShouldBindJSON` e `c.JSON` literalmente. Isso exigiria migração de roteador e mudaria a superfície HTTP sem necessidade para o padrão solicitado.

### Convenção e verificação

Atualizar os exemplos e orientações pertinentes de `docs/conventions/golang.md` para demonstrar interface local, construtor, método HTTP com Chi e mapeamento explícito de erros, sem contradizer seções vizinhas. A implementação será verificada por revisão dos contratos antes e depois, formatação e compilação. A preferência anterior de não criar nem executar testes permanece aplicada.

Alternativa considerada: acrescentar testes de caracterização nesta mudança. O usuário pediu implementação sem testes nas etapas anteriores; a verificação escolhida respeita essa restrição.

## Risks / Trade-offs

- [Assinaturas públicas dos construtores podem mudar] → Conferir todas as chamadas nos `module.go` e compilar todos os pacotes Go.
- [Refatoração pode alterar respostas por acidente] → Registrar os códigos, mensagens, corpos e rotas atuais antes da edição e revisar o diff contra esse inventário.
- [Mudanças pré-existentes nos handlers de OfficialSource] → Trabalhar sobre o estado atual, preservando o comportamento e as alterações já presentes.

## Migration Plan

Não há migração de dados nem de infraestrutura. Refatorar os handlers dos quatro módulos e a convenção no mesmo change; formatar e compilar o backend. Uma reversão consiste em reverter apenas os arquivos editados por esta mudança, sem afetar os handlers de OfficialSource pré-existentes.
