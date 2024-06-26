GOCMD=go
GORUN=${GOCMD} run
GOBUILD=$(GOCMD) build -ldflags="-s -w"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=bin
WIN_BINARY_NAME=${BINARY_NAME}/win
MAC_BINARY_NAME=${BINARY_NAME}/mac
LINUX_BINARY_NAME=${BINARY_NAME}/linux
NAME=sane

all: test build-mac build-win build-linux

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

run:
	${GORUN} ./main.go

# Cross compilation
# mac
build-mac:
	CGO_ENABLED=0 $(GOBUILD) -o $(MAC_BINARY_NAME)/$(NAME)
# linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(LINUX_BINARY_NAME)/$(NAME)
# windows
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(WIN_BINARY_NAME)/$(NAME).exe