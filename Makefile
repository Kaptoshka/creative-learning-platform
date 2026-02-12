.PHONY: generate
generate-from-zero:
	buf generate libs/protos
	cd libs/gen/go && go mod init github.com/Kaptoshka/creative-learning-platform/libs/gen/go && go mod tidy
generate:
	buf generate libs/protos
