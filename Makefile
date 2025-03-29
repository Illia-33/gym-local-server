.PHONY: all
all: configurator server

.PHONY: configurator
configurator:
	go build -o cfgr ./cmd/configurator

.PHONY: server
server:
	go build -o srvr ./cmd/localserver

.PHONY: config
config: configurator
	./cfgr

.PHONY: run
run: configurator server
	./cfgr
	./srvr
