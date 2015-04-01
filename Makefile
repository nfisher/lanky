# ex : shiftwidth=2 tabstop=2 softtabstop=2 :
SHELL := /bin/sh
GOPROCS := 4
SRC := $(wildcard *.go)
EXE := lanky

.PHONY: all
all: get-deps test vet $(EXE)

$(EXE): $(SRC)
	go build

.PHONY: get-deps
get-deps:
	go get -d -v ./...

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: format
format:
	go fmt ./...

coverage.out: $(SRC)
	go test -coverprofile=coverage.out

.PHONY: cov
cov: coverage.out
	go tool cover -func=coverage.out

.PHONY: htmlcov
htmlcov: coverage.out
	go tool cover -html=coverage.out

.PHONY: test
test:
	go test

.PHONY: run
run: $(EXE)
	./$(EXE)

.PHONY: vet
vet:
	go vet -x

.PHONY: install
install: cov
	go install

.PHONY: watch
watch:
	fswatch --one-per-batch *.go | xargs -n1 -I{} make test
