bootstrap-test:
	~/go/bin/ginkgo bootstrap

bootstrap-test:
	~/go/bin/ginkgo generate

run-tests:
	go test -p 1 -coverprofile coverage.out ./...

run-single-test:
	~/go/bin/ginkgo --focus "test 1" testDir

coverage:
	@go tool cover -html coverage.out -o coverage.html
	explorer.exe coverage.html
