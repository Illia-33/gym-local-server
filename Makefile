.PHONY: all
all: configurator server

.PHONY: configurator
configurator:
	go build -o configurator ./cmd/configurator

.PHONY: server
server:
	go build -o server ./cmd/localserver

.PHONY: config
config: configurator
	./configurator

.PHONY: run
run: configurator server
	./configurator
	./server
