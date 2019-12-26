prog = moni

gofiles = main.go page.go walker.go

build: $(prog )
	go build

run: $(gofiles)
	go -v run $(gofiles)

test:
	go test -v

install:
	go install -o moni $(gofiles)

.PHONY: run build
