cover:
	go test -short -count=1 -covermode=atomic -coverpkg=./internal/api/...,./internal/handlers/... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
