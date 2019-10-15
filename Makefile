mapimg = ~/Desktop/home.png

dotfile = etc/home.gv

all: build
	make -C moni build

go.mod:
	mod init github.com/rustyeddy/moni

build: 
	go build 

buildv: 
	go build -v

run:
	make -C moni run *.go

test:
	go test

testv:
	go test -v

ttv:
	go test -v -trace=test.out

serve:
	go run *.go -http

install:
	go install

mapimg:
	dot -Tpng -o $(mapimg) $(dotfile)

clean:
	go clean
	rm -rf *~
	rm -rf moni/*~
