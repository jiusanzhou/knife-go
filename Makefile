PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

default: help

getdeps:
	@echo "Installing golint" && go get -u golang.org/x/lint/golint
	@echo "Installing gocyclo" && go get -u github.com/fzipp/gocyclo
	@echo "Installing deadcode" && go get -u github.com/remyoudompheng/go-misc/deadcode
	@echo "Installing misspell" && go get -u github.com/client9/misspell/cmd/misspell
	@echo "Installing ineffassign" && go get -u github.com/gordonklaus/ineffassign
	@echo "Installing vendorcheck" && go get -u github.com/FiloSottile/vendorcheck
	@echo "Installing gometalinter" && go get -u github.com/alecthomas/gometalinter

verifiers: lint cyclo deadcode spelling metalinter

vet:
	@echo "Running $@"
	@go tool vet -atomic -bool -copylocks -nilfunc -printf -shadow -rangeloops -unreachable -unsafeptr -unusedresult cmd
	@go tool vet -atomic -bool -copylocks -nilfunc -printf -shadow -rangeloops -unreachable -unsafeptr -unusedresult pkg

fmt:
	@echo "Running $@"

lint:
	@echo "Running $@"
	@${GOPATH}/bin/golint -set_exit_status ./...

metalinter:
	@${GOPATH}/bin/gometalinter --install
	@${GOPATH}/bin/gometalinter --disable-all -E vet -E gofmt -E misspell -E ineffassign -E goimports -E deadcode --tests --vendor ./...

ineffassign:
	@echo "Running $@"
	@${GOPATH}/bin/ineffassign .

cyclo:
	@echo "Running $@"
	@${GOPATH}/bin/gocyclo -over 200 .

deadcode:
	@echo "Running $@"
	@${GOPATH}/bin/deadcode -test $(shell go list ./...) || true

spelling:
	@${GOPATH}/bin/misspell -locale US -error `find .`

check: verifiers test

test:
	@echo "Running unit tests"
	@go test -tags kqueue ./...

coverage:
	@echo "Running all coverage for knife-go"
	@(env bash $(PWD)/go-coverage.sh)

pkg-add:
	@echo "Adding new package $(PKG)"
	@${GOPATH}/bin/govendor add $(PKG)

pkg-update:
	@echo "Updating new package $(PKG)"
	@${GOPATH}/bin/govendor update $(PKG)

pkg-remove:
	@echo "Remove new package $(PKG)"
	@${GOPATH}/bin/govendor remove $(PKG)

pkg-list:
	@$(GOPATH)/bin/govendor list

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm -rvf build
	@rm -rvf release

help:
	@echo "nothing to do!"