.PHONY: deps test lint lint-check-deps ci-check run-migrations

deps: 
	@echo "[dep] fetching package dependencies"
#	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

test: 
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic ./...

lint: lint-check-deps
	@echo "[golangci-lint] linting sources"
	@golangci-lint run \
		-E misspell \
		-E golint \
		-E gofmt \
		-E unconvert \
		--exclude-use-default=false \
		./...

lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint";\
		go get -u github.com/golangci/golangci-lint/cmd/golangci-lint;\
	fi

ci-check: deps lint run-cdb-migrations test









