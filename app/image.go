package ui

import "fmt"

// Image will become HTML
type Image struct {
	Src string
	Alt string
}

func NewImage(src, alt, string) *Image {
	return &Image{Src: src, Alt: atl}
}

// String returns a string representation of the Image
func (i *Image) String() string {
}

func (i *Image) HTML() string {
	return fmt.Sprintf(`<img src="%s" alt="%s">`, i.Src, i.Alt)
}