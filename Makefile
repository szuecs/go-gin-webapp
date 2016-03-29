default: build.test

clean:
	rm -rf build

config:
	@test -d ~/.config/go-gin-webapp || mkdir -p ~/.config/go-gin-webapp
	@test -e ~/.config/go-gin-webapp/config.yaml || cp config.yaml.sample ~/.config/go-gin-webapp/config.yaml
	@echo "modify ~/.config/go-gin-webapp/config.yaml as you need"

all: build.linux build.osx

build.test:
	go test -v ./...
	go vet -v ./...

prepare:
	mkdir -p build/linux
	mkdir -p build/osx

# release
build.linux.release: prepare
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/linux/go-gin-webapp -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=`git describe --tags`" -tags zalandoValidation ./cmd/go-gin-webapp

# dev builds
build.local: prepare
	godep go build -o build/go-gin-webapp -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation  ./cmd/go-gin-webapp

build.linux: prepare
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/linux/go-gin-webapp -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation ./cmd/go-gin-webapp

build.osx: prepare
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/osx/go-gin-webapp -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation ./cmd/go-gin-webapp

dev.install:
	godep go install -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -X main.Version=HEAD" -tags zalandoValidation github.com/zalando-techmonkeys/go-gin-webapp-zmon-agg/...
