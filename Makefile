VERSION := $(shell grep -E '\bCliVersion\s*=\s*"' cmd/helpers/build/build.go | head -1 | sed -E 's/.*"([^"]+)".*/\1/')
ifeq ($(VERSION),)
$(error Could not extract CliVersion from cmd/helpers/build/build.go)
endif

.PHONY: all
all: clean get build
	@

.PHONY: build
build:
	@echo "Building version $(VERSION)"
	GORELEASER_CURRENT_TAG=$(VERSION) goreleaser release --clean --skip=publish --snapshot

.PHONY: get
get:
	go mod download

.PHONY: clean
clean:
	rm -rf bin
	mkdir bin

.PHONY: docs
docs:
	go run ./main.go docs

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: _linter
_linter:
	go install github.com/mgechev/revive@v1.10.0

.PHONY: lint
lint:
	golangci-lint run ./... && \
	revive -config .revive.toml ./...

.PHONY: _mockery-install
_mockery-install:
	go install github.com/vektra/mockery/v3@v3.3.4

.PHONY: mockery
mockery:
	mockery --config .mockery.yml

.PHONY: test_unit
test_unit:
	go test ./cmd/... ./services/... ./output/...

# add binary and config to tests/cmd/bin/ before run test integration
.PHONY: test_integration
test_integration:
	go test ./tests/cmd/files/push/...
	go test ./tests/cmd/files/pull/...
	go test ./tests/cmd/jobs/progress/...
	go test ./tests/cmd/files/list/...
	go test ./tests/cmd/files/status/...
	go test ./tests/cmd/files/rename/...
	go test ./tests/cmd/files/delete/...
	go test ./tests/cmd/projects/...
	go test ./tests/cmd/init/...
	go test ./tests/cmd/mt/detect/...
	go test ./tests/cmd/mt/translate/...
	go test ./tests/cmd/jobs/locales/...
