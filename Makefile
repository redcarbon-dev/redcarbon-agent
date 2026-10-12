.PHONY: generate run mocks test e2e

mocks:
	rm -rf mocks/
	mockery --with-expecter --all

run:
	go run cmd/main.go

test:
	go test ./...

e2e:
	docker build -f test/e2e/Dockerfile -t redcarbon-agent-e2e .
	docker run --rm --network none redcarbon-agent-e2e

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
