mapimg = ~/Desktop/home.png

dotfile = etc/home.gv

all:build
	make -C moni run

go.mod:
	mod init github.com/rustyeddy/moni

build: 
	go build -o $(cmd)

buildv: 
	go build -v $(cmd)

run:
	go run *.go

test:
	go test

testv:
	go test -v

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
