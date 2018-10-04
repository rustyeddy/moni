cmd = inv
mapimg = ~/Desktop/home.png

dotfile = etc/home.gv

all:build

go.mod:
	mod init github.com/rustyeddy/inv

build: 
	go build 

run:
	go run *.go

install:
	go install

mapimg:
	dot -Tpng -o $(mapimg) $(dotfile)

clean:
	go clean
	rm -rf *~
	rm -rf moni/*~
