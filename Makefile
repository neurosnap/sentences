BASE_DIR=$(shell echo $$GOPATH)/src/github.com/neurosnap/sentences
VERSION_FILE=$(BASE_DIR)/VERSION
CURRENT_VERSION=$(shell cat $(VERSION_FILE))
MAJOR=$(word 1, $(subst ., , $(CURRENT_VERSION)))
MINOR=$(word 2, $(subst ., , $(CURRENT_VERSION)))
PATCH=$(word 3, $(subst ., , $(CURRENT_VERSION)))
VER ?= $(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))

build:
	go build -ldflags "-X main.VERSION=$(CURRENT_VERSION)"

install:
	go install -ldflags "-X main.VERSION=$(CURRENT_VERSION)"

bump:
	echo $(VER) > $(VERSION_FILE)

# Compile language specific training data
czech:
	go-bindata -o data/czech.go data/czech.json

danish:
	go-bindata -o data/danish.go data/danish.json

dutch:
	go-bindata -o data/dutch.go data/dutch.json

english:
	go-bindata -o data/english.go data/english.json

estonian:
	go-bindata -o data/estonian.go data/estonian.json

finnish:
	go-bindata -o data/finnish.go data/finnish.json

french:
	go-bindata -o data/french.go data/french.json

german:
	go-bindata -o data/german.go data/german.json

greek:
	go-bindata -o data/greek.go data/greek.json

italian:
	go-bindata -o data/italian.go data/italian.json

norwegian:
	go-bindata -o data/norwegian.go data/norwegian.json

polish:
	go-bindata -o data/polish.go data/polish.json

portugues:
	go-bindata -o data/potuguese.go data/portuguese.json

slovene:
	go-bindata -o data/slovene.go data/slovene.json

spanish:
	go-bindata -o data/spanish.go data/spanish.json

turkish:
	go-bindata -o data/turkish.go data/turkish.json


