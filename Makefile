.PHONY: generate-protos generate-protos-fs lint-protos

generate-protos-fs:
	buf generate libs/protos
	cd libs/gen/go && go mod init github.com/Kaptoshka/creative-learning-platform/libs/gen/go && go mod tidy

generate-protos:
	buf generate libs/protos

lint-protos:
	buf lint libs/protos
