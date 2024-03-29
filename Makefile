GOBIN=$(shell pwd)/bin
BIN=$(GOBIN)/samedup
CGO=0
VERSION="0.2.6"
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GITCOMMIT=$(shell git rev-parse --short HEAD)
SRC="./cmd/samedup"
INSTALL_FLAG=-v -ldflags "-s -w -X 'samedup/internal/commands.version=$(VERSION)-$(GITCOMMIT)'"

all: build

build:
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(INSTALL_FLAG) -o $(BIN) $(SRC)

build-linux:
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=linux GOARCH=arm go build $(INSTALL_FLAG) -o $(BIN)-linux-arm $(SRC)
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=linux GOARCH=arm64 go build $(INSTALL_FLAG) -o $(BIN)-linux-arm64 $(SRC)
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=linux GOARCH=386 go build $(INSTALL_FLAG) -o $(BIN)-linux-386 $(SRC)
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=linux GOARCH=amd64 go build $(INSTALL_FLAG) -o $(BIN)-linux-amd64 $(SRC)

build-windows:
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=windows GOARCH=arm go build $(INSTALL_FLAG) -o $(BIN)-windows-arm $(SRC)
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=windows GOARCH=386 go build $(INSTALL_FLAG) -o $(BIN)-windows-386 $(SRC)
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=windows GOARCH=amd64 go build $(INSTALL_FLAG) -o $(BIN)-windows-amd64 $(SRC)

build-darwin:
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=darwin GOARCH=amd64 go build $(INSTALL_FLAG) -o $(BIN)-darwin-amd64 $(SRC)
	CGO_ENABLED=$(CGO) GO111MODULE=on GOOS=darwin GOARCH=arm64 go build $(INSTALL_FLAG) -o $(BIN)-darwin-arm64 $(SRC)

build-all: build-linux build-windows build-darwin

clean:
	@rm -rf $(GOBIN)

.PHONY: build build-linux build-windows build-darwin clean all
