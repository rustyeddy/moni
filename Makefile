prog = moni

gofiles = main.go page.go 

build: moni
	go build -o ${prog} 

run:
	go run ${gofiles}
