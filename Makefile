APP = nuki-cli
BUILD_DIR = build/bin
GIT_VER=$(shell git rev-parse HEAD)

LDFLAGS=-ldflags "-X main.version=${GIT_VER}"

.PHONY: test clean nuki-cli

# Build the project
all: clean test nuki-cli

nuki-cli:
	@mkdir -p ${BUILD_DIR}
	GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP}-linux-amd64 -v cmd/nuki-cli/main.go

test:
	go test -v

clean:
	-rm -f ${BUILD_DIR}/${BINARY}-*

distclean:
	rm -rf ./build

mrproper: distclean
	git ls-files --others | xargs rm -rf
