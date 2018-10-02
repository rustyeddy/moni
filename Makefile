cmd = moni
cmd/src = moni/main.go
mapimg = ~/Desktop/home.png 

dotfile = etc/home.gv

all: build

go.mod:
	mod init github.com/rustyeddy/inv

build: go.mod
	make -C $(cmd) build

run:
	go run $(cmdsrc) 

mapimg:
	dot -Tpng -o $(mapimg) $(dotfile)

clean:
	go clean
	rm -rf *~
	rm -rf moni/*~
