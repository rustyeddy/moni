cmd = moni
mapimg = ~/Desktop/home.png

dotfile = etc/home.gv

all:build

go.mod:
	mod init github.com/rustyeddy/inv

build: 
	go build -o $(cmd)

buildv: 
	go build -v $(cmd)

run:
	go run *.go

test:
	go test
	make -C store test

testv:
	go test -v
	make -C store test

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
