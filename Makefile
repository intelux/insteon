modules := serial plm

ifeq ($(OS),Windows_NT)
EXT:=.exe
else
EXT:=
endif

all: build lint test

build:
	@echo "### Building\n"
	go build -o bin/ion${EXT} ./ion
	@echo

lint:
	@echo "### Linting\n"
	@for module in $(modules); do \
		golint ./$$module; \
		go vet ./$$module; \
	done
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
