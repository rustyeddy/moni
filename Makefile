cmd = moni
cmd/src = moni/main.go
mapimg = ~/Desktop/home.png 

dotfile = etc/home.gv

build:
	make -C $(cmd) build

run:
	go run $(cmdsrc) 

mapimg:
	dot -Tpng -o $(mapimg) $(dotfile)

