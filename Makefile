.PHONY: run fmt vet check

run:
		go run ./...

fmt:
		go fmt ./...

vet:
		go vet ./...

check: fmt vet