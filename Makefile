modules := serial plm

all: build lint test

build:
	@echo "### Building\n"
	go build -o bin/ion ./ion
	@echo

lint:
	@echo "### Linting\n"
	golint ./...
	go vet ./...
	@echo

test:
	@echo "### Testing\n"
	@for module in $(modules); do \
		go test --coverprofile $$module.coverage ./$$module; \
	done
	@echo
	@echo "### Coverage\n"
	@for module in $(modules); do \
		go tool cover -func=$$module.coverage; \
	done
	@echo

coverage: test
	@for module in $(modules); do \
		go tool cover -html=$$module.coverage; \
	done
