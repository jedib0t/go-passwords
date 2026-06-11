.PHONY: test

default: test

bench:
	go test -bench=. -benchmem $(shell go list ./...)

cyclo:
	gocyclo -over 15 ./*/*.go

fmt:
	go fmt $(shell go list ./...)

test: fmt vet cyclo
	go test -race -cover -coverprofile=.coverprofile $(shell go list ./...)

tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@v0.5.1
	go install golang.org/x/vuln/cmd/govulncheck@latest

vulncheck:
	govulncheck ./...

vet:
	go vet $(shell go list ./...)

