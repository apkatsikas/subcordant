bootstrap-test:
	ginkgo bootstrap

run-tests:
	go test -p 1 -coverprofile coverage.out ./...

run-single-test:
	ginkgo --focus "test 1" testDir

check-formatting:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files need formatting:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

coverage:
	@go tool cover -html=coverage.out -o coverage.html
	@if grep -qi microsoft /proc/version; then \
		explorer.exe coverage.html; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		open coverage.html; \
	else \
		xdg-open coverage.html; \
	fi

install-staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

mocks:
	mockery

build:
	go build -o subcordant ./cmd/main.go
