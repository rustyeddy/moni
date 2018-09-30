mapimg = ~/Desktop/home.png 

dotfile = etc/home.gv

$(mapimg):
	dot -Tpng -o $(mapimg) $(dotfile)
