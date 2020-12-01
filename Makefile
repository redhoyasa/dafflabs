.PHONY: all
all: build fmt vet lint test test-cover-html vendor

SRC_PKGS=$(shell go list ./... | grep -v "vendor")
APP_EXECUTABLE=dafflabs

build: compile fmt vet lint

.PHONY: compile
compile:
	mkdir -p bin/
	GO111MODULE=on go build -o bin/$(APP_EXECUTABLE)

.PHONY: clean
clean:
	GO111MODULE=on go clean
	rm -rf bin/

.PHONY: fmt
fmt:
	GO111MODULE=on go fmt $(SRC_PKGS)

.PHONY: lint
lint:
	@for p in $(SRC_PKGS); do \
		echo "==> Linting $$p"; \
		golint $$p | grep -vwE "exported (function|method|type) \S+ should have comment" | true; \
	done

.PHONY: vet
vet:
	GO111MODULE=on go vet $(SRC_PKGS)

.PHONY: test
test:
	mkdir -p reports
	ENVIRONMENT=test GO111MODULE=on go test -coverprofile=reports/coverage.out ./...

.PHONY: test-cover-html
test-cover-html:
	GO111MODULE=on go tool cover -html=reports/coverage.out -o reports/coverage.html

.PHONY: vendor
vendor:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
