ifneq (,$(wildcard .env))
        include .env
        export
endif

.PHONY: update
update:
	go mod tidy
	go mod verify

.PHONY: test
test:
	gotestsum --format=testname -- ./... -tags=units

.PHONY: lint
lint:
	golangci-lint run