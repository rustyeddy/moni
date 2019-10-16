prog = moni

gofiles = main.go page.go walker.go

build: moni
	go build -o ${prog} 

run:
	go run ${gofiles}
