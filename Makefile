BINARY_DIR=./binary
CMD_DIR=./cmd/sentences

VERSION_FILE=./VERSION
CURRENT_VERSION=$(shell cat $(VERSION_FILE))

COMMITHASH=$(shell git rev-parse --short HEAD)

PREFIX?=/usr/local
BINDIR?=$(PREFIX)/bin

test:
	go test ./...
.PHONY: test

get:
	go get ./...
.PHONY: get

build:
	go build -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)" ${CMD_DIR}
.PHONY: build

# used to package language files with go binary
bindata:
	go get github.com/go-bindata/go-bindata/...
	go install github.com/go-bindata/go-bindata/...
.PHONY: bindata

# install:
# 	install -m755 sentences $(DESTDIR)$(BINDIR)/sentences
# .PHONY: install
#
# uninstall:
# 	rm -f $(DESTDIR)$(BINDIR)/sentences
# .PHONY: uninstall

cross:
	mkdir -p $(BINARY_DIR)

	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)" ${CMD_DIR}
	tar -czvf $(BINARY_DIR)/sentences_linux-amd64.tar.gz ./sentences

	GOOS=linux GOARCH=arm go build -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)" ${CMD_DIR}
	tar -czvf $(BINARY_DIR)/sentences_linux-arm.tar.gz ./sentences

	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)" ${CMD_DIR}
	tar -czvf $(BINARY_DIR)/sentences_linux-arm64.tar.gz ./sentences

	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)" ${CMD_DIR}
	tar -czvf $(BINARY_DIR)/sentences_darwin-amd64.tar.gz ./sentences

	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)" ${CMD_DIR}
	tar -czvf $(BINARY_DIR)/sentences_darwin-arm64.tar.gz ./sentences

	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)" ${CMD_DIR}
	tar -czvf $(BINARY_DIR)/sentences_windows-amd64.tar.gz ./sentences
.PHONY: cross

deploy: cross
	source ~/virtualenvs/aws/bin/activate
	aws s3 cp ./binary s3://sentence-binaries/ --recursive --exclude "*" --include "*.tar.gz"
.PHONY: deploy

install:
	go install -ldflags "-X main.VERSION=$(CURRENT_VERSION) -X main.COMMITHASH=$(COMMITHASH)" ${CMD_DIR}
.PHONY: install

bump:
	MAJOR=$(word 1, $(subst ., , $(CURRENT_VERSION)))
	MINOR=$(word 2, $(subst ., , $(CURRENT_VERSION)))
	PATCH=$(word 3, $(subst ., , $(CURRENT_VERSION)))
	VER ?= $(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))

	echo $(VER) > $(VERSION_FILE)
.PHONY: bump

# Compile language specific training data
czech:
	go-bindata -pkg="data" -o data/czech.go data/czech.json
.PHONY: czech

danish:
	go-bindata -pkg="data" -o data/danish.go data/danish.json
.PHONY: danish

dutch:
	go-bindata -pkg="data" -o data/dutch.go data/dutch.json
.PHONY: dutch

english:
	go-bindata -pkg="data" -o data/english.go data/english.json
.PHONY: english

estonian:
	go-bindata -pkg="data" -o data/estonian.go data/estonian.json
.PHONY: estonian

finnish:
	go-bindata -pkg="data" -o data/finnish.go data/finnish.json
.PHONY: finnish

french:
	go-bindata -pkg="data" -o data/french.go data/french.json
.PHONY: french

german:
	go-bindata -pkg="data" -o data/german.go data/german.json

greek:
	go-bindata -pkg="data" -o data/greek.go data/greek.json
.PHONY: greek

italian:
	go-bindata -pkg="data" -o data/italian.go data/italian.json
.PHONY: italian

norwegian:
	go-bindata -pkg="data" -o data/norwegian.go data/norwegian.json
.PHONY: norwegian

polish:
	go-bindata -pkg="data" -o data/polish.go data/polish.json
.PHONY: polish

portuguese:
	go-bindata -pkg="data" -o data/potuguese.go data/portuguese.json
.PHONY: portuguese

slovene:
	go-bindata -pkg="data" -o data/slovene.go data/slovene.json
.PHONY: slovene

spanish:
	go-bindata -pkg="data" -o data/spanish.go data/spanish.json
.PHONY: spanish

turkish:
	go-bindata -pkg="data" -o data/turkish.go data/turkish.json
.PHONY: turkish
