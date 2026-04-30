# goreview

Pequena API em Go usando Chi como router.

Pré-requisitos
- Go 1.20+

Como rodar

Instale dependências e rode o servidor:

```bash
cd /Users/dinorah/ufrn/web2/goReview
make tidy
make run
```

O servidor ficará disponível em `http://localhost:8080`.

Targets úteis
- `make run` — roda a aplicação
- `make test` — executa `go test ./...`
- `make tidy` — executa `go mod tidy`
- `make build` — compila binário

Endpoints

- `GET /health` — retorna status
- `GET /reviews` — lista todas as reviews
- `GET /reviews/{id}` — obtém review por id
- `POST /reviews` — cria review (JSON: `title`, `content`)
- `PUT /reviews/{id}` — atualiza review (JSON: `title?`, `content?`)

Exemplos com curl

```bash
# health
curl -sS http://localhost:8080/health

# listar
curl -sS http://localhost:8080/reviews

# criar
curl -sS -X POST -H "Content-Type: application/json" -d '{"title":"Teste","content":"Conteúdo"}' http://localhost:8080/reviews

# obter
curl -sS http://localhost:8080/reviews/1

# atualizar
curl -sS -X PUT -H "Content-Type: application/json" -d '{"content":"Atualizado"}' http://localhost:8080/reviews/1
```

Testes

Rode:

```bash
make test
```

Saída esperada: todos os testes passam (ok).
# goReview