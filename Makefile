help:
	@egrep -h '\s#@\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?#@ "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'

docs: #@ Generate docs
	swag init -g cmd/api/main.go
.PHONY:docs
test: fmt #@ Run tests
	go test -v ./...
fmt: #@ Format the code
	go fmt ./...
vet: fmt #@ VET the code
	go vet ./...
lint: fmt #@ Run the linter
	golint ./...
run: test docs vet #@ Start locally
	go run cmd/api/main.go
sync: test vet #@ Sync local data with API
	go run cmd/sync/main.go
update: #@ Update dependencies
	go mod tidy
clear-build: #@ Clear build folder
	rm -rf build && mkdir build && mkdir build/api && mkdir build/api/pkg && mkdir build/api/pkg/i18n
copy-translations: #@ Copy translations
	cp -r pkg/i18n/*.toml build/api/pkg/i18n
build: test docs vet clear-build copy-translations #@ Build the binary
	go build -o build/api/main cmd/api/main.go
.PHONY:build
image: docs vet #@ Build docker image
	docker build -t electricity-prices . --load
