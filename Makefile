VERSION=0.1.9

.PHONY: all
all: release

.PHONY: release
release: release-deps
	gox \
		-osarch="!darwin/386" \
		-ldflags="-s -w -X main.version=$(VERSION)" \
		-asmflags="-trimpath" \
		-output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		.

.PHONY: release-deps
release-deps:
	go install github.com/mitchellh/gox@latest

.PHONY: clean
clean:
	rm -rf build
