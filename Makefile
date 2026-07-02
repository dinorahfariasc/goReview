.PHONY: run test test-integration tidy build sqlc db-up db-down atlas-status atlas-apply atlas-diff atlas-hash

DATABASE_URL ?= postgres://goreview:goreview@127.0.0.1:5435/goreview?sslmode=disable
TEST_DATABASE_URL ?= $(DATABASE_URL)
ATLAS_DATABASE_URL ?= postgres://goreview:goreview@host.docker.internal:5435/goreview?sslmode=disable

run:
	DATABASE_URL=$(DATABASE_URL) go run main.go

test:
	go test ./...

test-integration:
	TEST_DATABASE_URL=$(TEST_DATABASE_URL) go test ./internal/adapter/postgres -run TestRepositoryIntegration -count=1

tidy:
	go mod tidy

build:
	go build -o goreview

sqlc:
	~/go/bin/sqlc generate

db-up:
	docker compose up -d

db-down:
	docker compose down

atlas-status:
	docker run --rm \
		-v $(PWD):/workspace -w /workspace \
		-e DATABASE_URL=$(ATLAS_DATABASE_URL) \
		arigaio/atlas:latest migrate status --env local

atlas-apply:
	docker run --rm \
		--add-host=host.docker.internal:host-gateway \
		-v $(PWD):/workspace -w /workspace \
		-e DATABASE_URL=$(ATLAS_DATABASE_URL) \
		arigaio/atlas:latest migrate apply --env local

atlas-diff:
	docker run --rm \
		--add-host=host.docker.internal:host-gateway \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):/workspace -w /workspace \
		-e DATABASE_URL=$(ATLAS_DATABASE_URL) \
		arigaio/atlas:latest migrate diff $(name) --env local

atlas-hash:
	docker run --rm \
		-v $(PWD):/workspace -w /workspace \
		arigaio/atlas:latest migrate hash --dir file://db/migrations
