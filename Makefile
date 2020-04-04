build:
	go build

.PHONY: fmt
fmt:
	goimports -local github.com/vllry/gameapi -w .

.PHONY: test
test:
	go test -cover -race ./...
