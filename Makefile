LOCAL_DB_DSN:=$(shell grep -A1 'database' values/local.yaml | tail -n1 | sed "s/.*dsn: //g" | sed "s/\"//g")

run:
	go run cmd/main.go

t:
	go test -v -count=1 -p=1 -cover -coverprofile=cover.out.tmp -covermode=atomic -coverpkg ./... ./...

migrate-new:
	migrate create -ext sql -dir migrations/ -seq init_schema

migrate-up:
	migrate -path migrations -database "$(LOCAL_DB_DSN)" up

migrate-down:
	migrate -path migrations -database "$(LOCAL_DB_DSN)" down

jet:bin-deps
	@PATH=$(LOCAL_BIN):$(PATH) jet -dsn $(LOCAL_DB_DSN) -path=./internal/generated/ -schema=public

bin-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/go-jet/jet/v2/cmd/jet@latest
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up-ci:bin-deps
	migrate -path migrations -database "$(LOCAL_DB_DSN)" up
