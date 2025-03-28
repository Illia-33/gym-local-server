.PHONY: all
all: configurator server

.PHONY: configurator
configurator:
	go build -o ./cfgr ./configurator

.PHONY: server
server:
	go build -o ./srvr ./server

.PHONY: config
config: configurator
	./cfgr

.PHONY: run
run: configurator server
	./cfgr
	./srvr
