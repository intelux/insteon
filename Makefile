modules := serial

all: build lint test

build:
	@echo "### Building\n"
	go build ./serial
	@echo

lint:
	@echo "### Linting\n"
	golint ./...
	go vet ./...
	@echo

test:
	@echo "### Testing\n"
	@for module in $(modules); do \
		go test --coverprofile $$module.coverage ./$$module --trace=trace; \
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
