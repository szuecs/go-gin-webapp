.PHONY: clean clean check build.docker scm-source test

APP           ?= appname
REPO_USER     ?= repo/user
SED           ?= sed

BINARY_BASE   ?= go-gin-webapp
TEAM          ?= teapot
REGISTRY      ?= registry.opensource.zalan.do
IMAGE_NAME    ?= $(BINARY_BASE)
VERSION       ?= $(shell git describe --tags --always --dirty)
GIT_NAME      ?= $(shell git config --global --get user.name)
GIT_EMAIL     ?= $(shell git config --global --get user.email)
IMAGE         ?= $(REGISTRY)/$(TEAM)/$(IMAGE_NAME)
TAG           ?= $(VERSION)
TARGET_GOOS   ?= linux
TARGET_GOARCH ?= amd64
DOCKERFILE    ?= Dockerfile
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -X main.Version=$(VERSION) -X main.Buildstamp=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p') -X main.Githash=$(shell git rev-parse HEAD)
GITHEAD   = $(shell git rev-parse --short HEAD)
GITURL    = $(shell git config --get remote.origin.url)
GITSTATUS = $(shell git status --porcelain || echo "no changes")
SOURCES   = $(shell find . -name '*.go')
DST       = $(GOPATH)/src/$(REPO_USER)/$(APP)

default: build.local

clean:
	rm -rf build
	rm -f test/bench.*
	rm -f test/prof.*
	find . -name '*.test' -delete

config:
	@test -d ~/.config/$(BINARY_BASE) || mkdir -p ~/.config/$(BINARY_BASE)
	@test -e ~/.config/$(BINARY_BASE)/config.yaml || cp config.yaml.sample ~/.config/$(BINARY_BASE)/config.yaml
	@echo "modify ~/.config/$(BINARY_BASE)/config.yaml as you need"

test:
	go test $(shell glide novendor)

check:
	golint $(shell glide novendor)
	go vet -v $(shell glide novendor)

build.local: build/$(BINARY_BASE)
build.linux: build/linux/$(BINARY_BASE)
build.osx: build/osx/$(BINARY_BASE)

build/$(BINARY_BASE): $(SOURCES)
	go build -o build/"$(BINARY_BASE)" "$(BUILD_FLAGS)" -ldflags "$(LDFLAGS)" -tags zalandoValidation ./cmd/$(BINARY_BASE)

build/linux/$(BINARY_BASE): $(SOURCES)
	GOOS=linux GOARCH=$(TARGET_GOARCH) CGO_ENABLED=0 go build "$(BUILD_FLAGS)" -o build/linux/"$(BINARY_BASE)" -ldflags "$(LDFLAGS)" -tags zalandoValidation ./cmd/$(BINARY_BASE)

build/osx/$(BINARY_BASE): $(SOURCES)
	GOOS=darwin GOARCH=$(TARGET_GOARCH) CGO_ENABLED=0 go build "$(BUILD_FLAGS)" -o build/osx/"$(BINARY_BASE)" -ldflags "$(LDFLAGS)" -tags zalandoValidation ./cmd/$(BINARY_BASE)

$(DOCKERFILE).upstream: $(DOCKERFILE)
	$(SED) "s@UPSTREAM@$(shell $(shell head -1 $(DOCKERFILE) | $(SED) -E 's@FROM (.*)/(.*)/(.*):.*@pierone tags \2 \3 --url \1@') | awk '{print $$3}' | tail -1)@" $(DOCKERFILE) > $(DOCKERFILE).upstream

build.docker: $(DOCKERFILE).upstream scm-source.json build.linux
	docker build --rm -t "$(IMAGE):$(TAG)" -f $(DOCKERFILE).upstream .

build.push: build.docker
	docker push "$(IMAGE):$(TAG)"

scm-source.json: .git
	scm-source

build.rkt: scm-source.json build.linux
	acbuild begin
	acbuild set-name $(TEAM)/$(BINARY_BASE)
	acbuild copy build/linux/$(BINARY_BASE) /$(BINARY_BASE)
	acbuild copy config.yaml.sample /root/.config/$(BINARY_BASE)/config.yaml
	acbuild copy scm-source.json /scm-source.json
	acbuild set-exec -- /$(BINARY_BASE) --logtostderr -debug -v=2
	acbuild port add 8080 tcp 8080
	acbuild label add version $(VERSION)
	acbuild label add arch $(TARGET_GOARCH)
	acbuild label add os $(TARGET_GOOS)
	acbuild annotation add authors "$(GIT_NAME) <$(GIT_EMAIL)>"
	acbuild write $(BINARY_BASE)-$(VERSION).$(TARGET_GOOS)-$(TARGET_GOARCH).aci
	acbuild end

create.app: create.$(APP)

create.$(APP):
	mkdir -p $(DST)
	rsync -a --exclude=.git $(GOPATH)/src/github.com/szuecs/go-gin-webapp/ $(DST)
	cd $(DST)
	grep -rl github.com/szuecs/go-gin-webapp * | xargs $(SED) -i "s@github.com/szuecs/go-gin-webapp@$(REPO_USER)/$(APP)@"
	grep -rl go-gin-webapp * | xargs $(SED) -i "s@go-gin-webapp@$(APP)@g"
	mv cmd/go-gin-webapp cmd/$(APP)
	mv cmd/go-gin-webapp-cli cmd/$(APP)-cli
	echo "# $(APP)" > README.md
	git add .
	git commit -m "init $(APP)"
	go get -u github.com/Masterminds/glide/...
	glide i

