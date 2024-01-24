## 3rd-party packages
* sqlc CLI - [source](https://github.com/sqlc-dev/sqlc) | [Info](https://sqlc.dev/)
* Goose CLI - [source](https://github.com/pressly/goose)

### Generate Queries
Use [sqlc](https://github.com/sqlc-dev/sqlc) CLI:
```shell
sqlc generate
```

### Migration
Use [Goose](https://github.com/pressly/goose) CLI:
```shell
goose -dir=sql/schemas postgres <db/connection/string> up
```
Create new migration file:
```shell
goose -dir=sql/schemas create <brief_description> sql
```