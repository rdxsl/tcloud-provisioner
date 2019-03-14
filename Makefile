
BINARY := tcloud-provisioner

UNIX_EXECUTABLES := \
	darwin/amd64/$(BINARY) \
	linux/amd64/$(BINARY) \

ALL_EXECUTABLES := $(UNIX_EXECUTABLES:%=bin/%)

GOFILES := $(shell git ls-files | egrep \.go$ | egrep -v ^vendor/ | egrep -v _test.go$)

all: clean testwithrace build

build: clean $(ALL_EXECUTABLES)

clean:
	rm -rf bin/ pkg/ $(BINARY)

# Run unittests
test:
	$(TEST_ENV_VARS) go test $(TEST_FLAGS) $(ALL_PACKAGES)

# Run unittests with race condition detector on (takes longer)
testwithrace:
	$(TEST_ENV_VARS) go test $(TEST_FLAGS) -race $(ALL_PACKAGES)


bin/darwin/amd64/$(BINARY): $(GOFILES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo -ldflags="-s $(VERSION_UPDATE_FLAG)" -o "$@" main.go

bin/linux/amd64/$(BINARY): $(GOFILES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags="-s $(VERSION_UPDATE_FLAG)" -o "$@" main.go
