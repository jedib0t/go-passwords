.PHONY: test

default: test

bench:
	go test -bench=. -benchmem $(shell go list ./...)

cyclo:
	gocyclo -over 15 ./*/*.go

fmt:
	go fmt $(shell go list ./...)

test: fmt vet cyclo
	go test -cover -coverprofile=.coverprofile $(shell go list ./...)

tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@v0.5.1

vet:
	go vet $(shell go list ./...)

