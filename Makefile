.PHONY: gen-proto lint-protos

gen-proto:
	buf generate libs/protos --template libs/protos/buf.gen.yaml

lint-protos:
	buf lint libs/protos
