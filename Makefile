bootstrap-test:
	ginkgo bootstrap

run-tests:
	go test -p 1 -coverprofile coverage.out ./...

run-single-test:
	ginkgo --focus "test 1" testDir

coverage:
	@go tool cover -html=coverage.out -o coverage.html
	@if grep -qi microsoft /proc/version; then \
		explorer.exe coverage.html; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		open coverage.html; \
	else \
		xdg-open coverage.html; \
	fi

mocks:
	mockery

build:
	go build -o subcordant ./cmd/main.go
