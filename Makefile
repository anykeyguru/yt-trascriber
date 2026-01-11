APPNAME := ytb


VERSION ?= dev
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X main.AppName=$(APPNAME) \
	-X main.Version=$(VERSION) \
	-X main.Commit=$(COMMIT) \
	-X main.Date=$(DATE)

BUILD_CMD_LINUX := CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -trimpath -ldflags="$(LDFLAGS)"

BUILD_CMD := CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
	go build -trimpath -ldflags="$(LDFLAGS)"\

.PHONY: build tidy clean

build: tidy
	$(BUILD_CMD) -o $(APPNAME) ./cmd/ytranscribe

tidy:
	go mod tidy

clean:
	rm -f $(APPNAME)
