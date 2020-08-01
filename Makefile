DOCKER_HUB_REPO?=tusharraut
DOCKER_HUB_REGISTRY_IMAGE?=internal-service-monitor
DOCKER_HUB_REGISTRY_TAG?=1.0.0

REGISTRY_IMG=$(DOCKER_HUB_REPO)/$(DOCKER_HUB_REGISTRY_IMAGE):$(DOCKER_HUB_REGISTRY_TAG)

export GO111MODULE=on
export GOFLAGS = -mod=vendor
HAS_GOMODULES := $(shell go help mod why 2> /dev/null)
ifndef HAS_GOMODULES
$(error service monitor can only be built with go 1.14+ which supports go modules)
endif

ifndef PKGS
PKGS := $(shell GOFLAGS=-mod=vendor go list ./... 2>&1)
endif

GO_FILES := $(shell find . -name '*.go' | grep -v vendor | \
                                   grep -v '\.pb\.go' | \
                                   grep -v '\.pb\.gw\.go' | \
                                   grep -v 'externalversions' | \
                                   grep -v 'versioned' | \
                                   grep -v 'generated')

BASE_DIR    := $(shell git rev-parse --show-toplevel)
GIT_SHA     := $(shell git rev-parse --short HEAD)
BIN         :=$(BASE_DIR)/bin

LDFLAGS += "-s -w"
BUILD_OPTIONS := -ldflags=$(LDFLAGS)

.DEFAULT_GOAL: all

all: internal-service-monitor pretest test

internal-service-monitor:
	@echo "Bin directory: $(BIN)"
	CGO_ENABLED=0 go build $(BUILD_OPTIONS) -o $(BIN)/internal-service-monitor

container:
	@echo "Building container: docker build --tag $(REGISTRY_IMG) -f Dockerfile ."
	sudo docker build --tag $(REGISTRY_IMG) -f Dockerfile .

deploy:
	sudo docker push $(REGISTRY_IMG)

test:
	echo "" > coverage.txt
	for pkg in $(PKGS);	do \
		go test -v -coverprofile=profile.out -covermode=atomic -coverpkg=$${pkg}/... $${pkg} || exit 1; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done

lint:
	go get -u golang.org/x/lint/golint
	for file in $(GO_FILES); do \
        golint $${file}; \
        if [ -n "$$(golint $${file})" ]; then \
            exit 1; \
        fi; \
    done

vet:
	go vet $(PKGS)

staticcheck:
	go get -u honnef.co/go/tools/cmd/staticcheck
	staticcheck $(PKGS)

errcheck:
	go get -u github.com/kisielk/errcheck
	errcheck -ignoregenerated -verbose -blank $(PKGS)

check-fmt:
	bash -c "diff -u <(echo -n) <(gofmt -l -d -s -e $(GO_FILES))"

do-fmt:
	gofmt -s -w $(GO_FILES)

gocyclo:
	go get -u github.com/fzipp/gocyclo
	gocyclo -over 15 $(GO_FILES)

pretest: lint vet staticcheck

imports:
	goimports -w $(GO_FILES)

clean:
	rm -rf ./bin/*
