GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=gopherss
MAIN_PATH=cmd/$(BINARY_NAME)/main.go

all: test build
build: 
		$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		make build
		./$(BINARY_NAME)
