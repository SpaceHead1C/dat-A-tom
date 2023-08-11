.SILENT:

BINARY_NAME = datatom

MOCKED_REPOS = Property Record RefType Value ChangedData StoredConfig
MOCKED_BROKERS = Property Record RefType Value
GENERATED_MOCKS = $(foreach var,$(MOCKED_REPOS),./test/mocks/$(var)Repository.go) $(foreach var,$(MOCKED_BROKERS),./test/mocks/$(var)Broker.go)
MOCK_SOURCE = changed_data.go property.go record.go ref_type.go stored_configs.go value.go
COVERAGE = coverage.out

default: build

all: tests coverage-total build

mocks:
	go generate ./test/mocks/...

$(GENERATED_MOCKS): $(foreach var,$(MOCK_SOURCE),./internal/domain/$(var))
	$(MAKE) mocks

tests: $(GENERATED_MOCKS)
	go test -short -count=1 -covermode=atomic -coverpkg=./internal/api/...,./internal/handlers/... -coverprofile=$(COVERAGE) ./...

coverage: $(COVERAGE)
	go tool cover -html=$(COVERAGE)

coverage-total: $(COVERAGE)
	go tool cover -func $(COVERAGE) | grep total | awk '{print $3}'

$(COVERAGE):
	$(MAKE) tests

proto:
	go generate ./third_party/dataway

clean:
	go clean
	rm -f ./$(COVERAGE)
	rm -f $(wildcard ./internal/pb/*.pb.go)
	rm -f $(foreach var,$(GENERATED_MOCKS),$(var))
	rm -f $(wildcard ./.build/${BINARY_NAME}-*)


build: proto
	GOARCH=amd64 GOOS=darwin go build -o '.build/${BINARY_NAME}-darwin-amd' ./cmd/app
	GOARCH=arm64 GOOS=darwin go build -o '.build/${BINARY_NAME}-darwin' ./cmd/app
	GOARCH=amd64 GOOS=linux go build -o '.build/${BINARY_NAME}-linux' ./cmd/app
	GOARCH=amd64 GOOS=windows go build -o '.build/${BINARY_NAME}-windows.exe' ./cmd/app

run: build
	'.build/${BINARY_NAME}-windows' -c ./conf/config.toml

.PHONY: mocks proto clean