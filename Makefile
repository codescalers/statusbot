GOPATH := $(shell go env GOPATH)

getverifiers:
	@echo "Installing staticcheck" && go get -u honnef.co/go/tools/cmd/staticcheck && go install honnef.co/go/tools/cmd/staticcheck
	@echo "Installing gocyclo" && go get -u github.com/fzipp/gocyclo/cmd/gocyclo && go install github.com/fzipp/gocyclo/cmd/gocyclo
	@echo "Installing deadcode" && go get -u github.com/remyoudompheng/go-misc/deadcode && go install github.com/remyoudompheng/go-misc/deadcode
	@echo "Installing misspell" && go get -u github.com/client9/misspell/cmd/misspell && go install github.com/client9/misspell/cmd/misspell
	@echo "Installing golangci-lint" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	go mod tidy

verifiers: fmt lint cyclo deadcode spelling staticcheck

fmt:
	@echo "Running $@"
	@gofmt -d .

lint:
	@echo "Running $@"
	@${GOPATH}/bin/golangci-lint run

cyclo:
	@echo "Running $@"
	@${GOPATH}/bin/gocyclo -over 100 .

deadcode:
	@echo "Running $@"
	@${GOPATH}/bin/deadcode -test $(shell go list ./...) || true

spelling:
	@echo "Running $@"
	@${GOPATH}/bin/misspell -i monitord -error `find .`

staticcheck:
	@echo "Running $@"
	@${GOPATH}/bin/staticcheck -- ./...

build:
	@echo "Running $@"
	@go build -o statusbot main.go
