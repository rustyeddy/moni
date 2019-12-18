prog = moni

gofiles = main.go page.go walker.go

build: $(prog )
	go build -o ${prog} 

run: $(gofiles)
	go -v run $(gofiles)

test:
	go test -v *.go

install:
	go install -o moni $(gofiles)

.PHONY: run build
