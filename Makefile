gofiles = main.go page.go

moni: $(gofiles)
	go build -v -o moni $(gofiles)

build: $(gofiles)
	go -v build $(gofiles)

run: $(gofiles)
	go -v run $(gofiles)

test:
	go test -v *.go

install:
	go install -o moni $(gofiles)

.PHONY: run build
