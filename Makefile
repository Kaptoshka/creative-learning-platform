.PHONY: generate
generate:
	buf generate libs/protos
	cd libs/gen/go && go mod init github.com/Kaptoshka/creative-learning-platform/libs/gen/go || true && go mod tidy
