GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
NAME=qemu

all: clean build
build:
	$(GOBUILD) -ldflags="-s -w" -o $(NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(NAME)
deps:
	# deps
	$(GOGET) -u github.com/libvirt/libvirt-go-xml

	# linters
	curl -sfL https://raw.githubusercontent.com/alecthomas/gometalinter/master/scripts/install.sh | sh -s -- -b $$GOPATH/bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $$GOPATH/bin v1.12.2
	$(GOGET) -u github.com/go-critic/go-critic/...
	$(GOGET) -u github.com/securego/gosec/cmd/gosec/...
lint:
	gometalinter
	golangci-lint run