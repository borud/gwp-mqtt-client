ifeq ($(GOPATH),)
GOPATH := $(HOME)/go
endif

ifeq ($(OS),Windows_NT)
	EXTENSION=.exe
else
	EXTENSION=""
endif

all: test lint vet build

build: gmc

gmc:
	@echo "*** building $@"
	@cd cmd/$@ && go build -o ../../bin/$@$(EXTENSIION) --trimpath -tags osusergo,netgo -ldflags="$(LDFLAGS) -s -w"

test:
	@echo "*** $@"
	@go test ./...

vet:
	@echo "*** $@"
	@go vet ./...

lint:
	@echo "*** $@"
	@revive ./... 

count:
	@gocloc --not-match-d=pkg/gwp.v1 .

clean:
	@rm -rf bin
	@rm -rf pkg/gwp.v1
	@rm -rf dist

install-deps:
	@go install github.com/mgechev/revive@latest
