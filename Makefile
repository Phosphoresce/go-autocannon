GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0

all: build

build:
	$(GOC) ac.go

run:
	go run ac.go

stat:
	$(CGOR) $(GOC) $(GOFLAGS) ac.go

fmt:
	gofmt -w .

clean:
	rm ac
