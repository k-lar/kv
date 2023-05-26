PREFIX ?= /usr
BINARY_NAME=kv

.PHONY: all
all: build

.PHONY: build
build:
	go build -o ${BINARY_NAME} main.go


.PHONY: run
run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

.PHONY: install
install:
	@install -Dm755 kv $(DESTDIR)$(PREFIX)/bin/kv

.PHONY: uninstall
uninstall:
	@rm -f $(DESTDIR)$(PREFIX)/bin/kv

.PHONY: clean
clean:
	go clean
	rm ${BINARY_NAME}

