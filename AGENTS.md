# Benefit

Monorepo com frontend, backend principal em Go e (futuro) serviço de IA.

## Estrutura

- `frontend/`: Vite + pnpm.
- `backend/benefit-api/`: API principal em Go, monolito modular (PostgreSQL, migrations com Goose).
- `backend/ai-service/`: serviço de IA em FastAPI. **Ainda não implementado**; não criar código ali sem pedido explícito.
- `docs/`: documentação de arquitetura e convenções.

## Antes de editar código

Leia o `AGENTS.md` da área em que for trabalhar:

- Go: `backend/benefit-api/AGENTS.md`
- Frontend: `frontend/AGENTS.md` (se existir)

Se o arquivo da área não existir, siga os padrões do código vizinho.

## Comandos

Rode sempre dentro da pasta do projeto.

Frontend (`frontend/`):
- `pnpm install`
- `pnpm dev`
- `pnpm build`

Backend (`backend/benefit-api/`): veja o `AGENTS.md` de lá.

## Regras gerais

- Use `pnpm` no frontend. Não use npm nem yarn.
- Nunca leia, edite ou exiba o conteúdo de arquivos `.env`. Para novas variáveis, atualize apenas `.env.example`.
- Não adicione dependências sem justificar.
- Mudanças pequenas e focadas; não refatore código fora do escopo do pedido.
- Use a `PRD.md` para um contexto geral do projeto.
- Se a duvida puder ser respondida olhando para o codigo, use ele.
