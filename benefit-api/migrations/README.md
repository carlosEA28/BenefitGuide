# Database migrations

Install the Goose CLI if it is not already available:

```sh
go install github.com/pressly/goose/v3/cmd/goose@v3.27.1
```

Put Goose SQL migrations in this directory. For example:

```sh
goose -dir migrations create create_users sql
```

Run them from the project root:

```sh
make migration-up       # Apply all pending migrations
make migration-rollback # Revert the latest migration
make migration-down     # Revert all applied migrations
```

The Makefile uses the same `DB_*` variables as the application. Set them in
`.env` or in the environment. The `.env` file must use shell-compatible
`KEY=value` syntax.
