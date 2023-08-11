.SILENT:

MOCKED_REPOS = Property Record RefType Value ChangedData StoredConfig
MOCKED_BROKERS = Property Record RefType Value
GENERATED_MOCKS = $(foreach var,$(MOCKED_REPOS),./test/mocks/$(var)Repository.go) $(foreach var,$(MOCKED_BROKERS),./test/mocks/$(var)Broker.go)
MOCK_SOURCE= changed_data.go property.go record.go ref_type.go stored_configs.go value.go
COVERAGE = coverage.out

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

clean:
	rm -f ./$(COVERAGE)

.PHONY: mocks