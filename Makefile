.PHONY: generate run

run:
	go run main.go sources
run-check-sources-errors:
	go run main.go checksources

test:
	go test ./...

sync-db:
	pg_dump postgres://postgres:iv6najnlzz@localhost:5433/sources > /tmp/bk.sql
	dbmate -u postgres://postgres:postgres@localhost:5432/rc_sources?sslmode=disable drop
	dbmate -u postgres://postgres:postgres@localhost:5432/rc_sources?sslmode=disable create
	psql postgres://postgres:postgres@localhost:5432/rc_sources < /tmp/bk.sql

pgcli:
	pgcli postgres://postgres:postgres@localhost:5432/rc_sources

GOBIN=${PWD}/.bin

download:
	echo Download go.mod dependencies
	go mod download

install-tools: download
	echo Installing tools from tools.go
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % env GOBIN=${GOBIN} go install %

generate-proto:
	PATH=${GOBIN} buf generate

generate: generate-proto
