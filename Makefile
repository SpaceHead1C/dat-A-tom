.SILENT:

RELEASE_VERSION = 0.3.0
BINARY_NAME = datatom

PB_NAMES = dataway dataway_grpc
GENERATED_PB = $(foreach var,$(PB_NAMES),./internal/pb/$(var).pb.go)

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

$(GENERATED_PB): ./third_party/dataway/dataway.proto
	$(MAKE) proto

tests: $(GENERATED_MOCKS) $(GENERATED_PB)
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
	rm -f $(wildcard ./.build/$(BINARY_NAME)-*)


build: $(GENERATED_PB)
	go build -v -ldflags \
		'-X main.Version=$(RELEASE_VERSION)' \
		-o '.build/$(BINARY_NAME)' ./cmd/app

run: build
	'.build/$(BINARY_NAME)' -c ./conf/config.toml

.PHONY: mocks proto clean