mapimg = ~/Desktop/home.png 

dotfile = etc/home.gv

build:
	make -C cmd build

run:
	go run cli/main.go serve

mapimg:
	dot -Tpng -o $(mapimg) $(dotfile)

