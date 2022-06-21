GO=go
HOOKS=main
GOFLAGS=-ldflags="-s -w"

build: $(HOOKS)

main:
	go env -w GO111MODULE=on
	$(GO) build $(GOFLAGS) -o $@ ./mainhook

tidy:
	go mod tidy -compat=1.17

clean:
	rm -f $(HOOKS)

