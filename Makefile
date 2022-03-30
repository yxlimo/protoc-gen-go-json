.PHONY: proto
proto:
	mkdir -p build
	go build -o build/protoc-gen-go-json .
	export PATH=$(CURDIR)/build/:$$PATH && buf generate
