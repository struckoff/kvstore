# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -v
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -race -v
GOGET=$(GOCMD) get
BINARY_FOLDER=bin
STORE_BINARY_NAME=kvstore
ROUTER_BINARY_NAME=kvrouter

.PHONY: all
all: store router

.PHONY: store
store:
	@mkdir -p $(BINARY_FOLDER)
	$(GOBUILD) -o $(BINARY_FOLDER)/$(STORE_BINARY_NAME) ./store/cmd/.

.PHONY: store-test
store-test:
	$(GOTEST) -v ./store/...

.PHONY: store-clean
store-clean:
	$(GOCLEAN)
	rm -f $(BINARY_FOLDER)/$(STORE_BINARY_NAME)

.PHONY: store-docker
store-docker:
	docker build -f store/build.Dockerfile

.PHONY: store-docker-debug
store-docker-debug:
	docker build -f store/debug.Dockerfile

.PHONY: router
router:
	@mkdir -p $(BINARY_FOLDER)
	$(GOBUILD) -o $(BINARY_FOLDER)/$(ROUTER_BINARY_NAME) ./router/cmd/.

.PHONY: router-test
router-test:
	$(GOTEST) ./router/...

.PHONY: router-clean
router-clean:
	$(GOCLEAN)
	rm -f $(BINARY_FOLDER)/$(ROUTER_BINARY_NAME)
