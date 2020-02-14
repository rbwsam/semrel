CMD ?= ./cmd/semrel
BIN_NAME ?= semrel
LD_FLAGS ?= -s -w

build:
	@go build -ldflags="$(LD_FLAGS)" -o $(BIN_NAME) $(CMD) && echo "./$(BIN_NAME)"

install:
	@go install -ldflags="$(LD_FLAGS)" $(CMD) && echo "$(GOPATH)/bin/$(BIN_NAME)"

uninstall:
	@rm -f $(GOPATH)/bin/$(BIN_NAME)

clean:
	@rm -rf $(BIN_NAME)

.PHONY: build install uninstall clean