# goreview

API REST em Go com PostgreSQL, `sqlc`, migrations com Atlas e separacao em camadas inspirada em Clean Architecture.

## Arquitetura

Estrutura principal:

- `internal/domain` — entidades e DTOs de entrada
- `internal/usecase` — regras de negocio e contratos
- `internal/adapter/http` — handlers e router HTTP
- `internal/adapter/postgres` — implementacao de repositorio usando `sqlc`
- `internal/app` — bootstrap da aplicacao
- `internal/config` — leitura de configuracao
- `db/schema` — schema desejado usado por `sqlc` e Atlas
- `db/migrations` — migrations versionadas do Atlas
- `db/sqlc` — codigo gerado

## Stack

- Go
- Chi
- PostgreSQL
- `sqlc`
- Atlas
- Docker Compose

## Como rodar

Suba o banco:

```bash
make db-up
```

O PostgreSQL sobe vazio. A criacao de tabelas e controlada pelo Atlas.

Rode a API:

```bash
make run
```

Padrao local:

```bash
DATABASE_URL=postgres://goreview:goreview@localhost:5432/goreview?sslmode=disable
```

## Makefile

- `make run` — sobe a API
- `make test` — executa todos os testes unitarios
- `make test-integration` — executa o teste de integracao do repositorio Postgres
- `make sqlc` — regenera o codigo do `sqlc`
- `make db-up` — sobe o PostgreSQL
- `make db-down` — derruba o ambiente local
- `make atlas-status` — mostra o estado das migrations
- `make atlas-apply` — aplica migrations no banco alvo
- `make atlas-diff name=nome_da_migration` — gera nova migration a partir de `db/schema`

## Migrations com Atlas

Configuracao:

- [atlas.hcl](/Users/dinorah/ufrn/web2/goReview/atlas.hcl)
- [db/migrations/202605220001_init.sql](/Users/dinorah/ufrn/web2/goReview/db/migrations/202605220001_init.sql)
- [db/schema/schema.sql](/Users/dinorah/ufrn/web2/goReview/db/schema/schema.sql)

Fluxo:

```bash
make atlas-hash
make atlas-status
make atlas-apply
make atlas-diff name=add_new_table
```

Os comandos usam a imagem oficial `arigaio/atlas`, entao nao dependem de uma instalacao local do Atlas.

## Endpoints

- `GET /health`
- `GET /movies`
- `GET /movies/{id}`
- `POST /movies`
- `PUT /movies/{id}`
- `DELETE /movies/{id}`
- `GET /movies/{id}/details`
- `GET /movies/{id}/reviews`
- `POST /movies/{id}/reviews`
- `GET /reviews/{id}`
- `PUT /reviews/{id}`
- `DELETE /reviews/{id}`

## Exemplo rapido

```bash
curl -sS -X POST http://localhost:8080/movies \
  -H "Content-Type: application/json" \
  -d '{"title":"Arrival","synopsis":"Sci-fi first contact","release_year":2016}'

curl -sS -X POST http://localhost:8080/movies/1/reviews \
  -H "Content-Type: application/json" \
  -d '{"reviewer_name":"Dinorah","rating":5,"content":"Excelente atmosfera"}'

curl -sS http://localhost:8080/movies/1/details
```

## Testes

Cobertura atual:

- testes de use case
- testes de handler HTTP
- teste de integracao do repositorio Postgres

Para o teste de integracao:

```bash
TEST_DATABASE_URL=postgres://goreview:goreview@localhost:5432/goreview?sslmode=disable make test-integration
```
