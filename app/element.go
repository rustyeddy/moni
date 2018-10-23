package app

import (
	"fmt"
)

// Element will become an HTML element
type ElementIface interface {
	OpenTag() string
	CloseTag() string
}

// Node is an Element that has other Elements for children
type Element struct {
	Name    string
	Id      string
	Classtr string
}

func NewElement(name, id string, classes string) Element {
	return Element{
		Name:    name,
		Id:      id,
		Classtr: classes,
	}
}

type Node struct {
	Element
	Children []*Node
}

type Data struct {
	Element
	content string
}

type Row struct {
	Element
}

func NewRow() (r Row) {
	r.Element = NewElement("div", "", "row")
	return r
}

type Col struct {
	Element
}

func NewCol(colsize string) (c Col) {
	c.Name = "col" + "-" + string(colsize)
	c.Element = NewElement("div", "", "c.Name")
	return c
}

func (n *Element) OpenTag() (tag string) {
	tag = "<" + n.Name
	if n.Id != "" {
		tag += " id='" + n.Id + "'"
	}
	if cls := n.Classtr; cls != "" {
		tag += " class='" + cls + "'"
	}
	tag += ">"
	return tag
}

func (n *Element) CloseTag() (tag string) {
	return fmt.Sprintf("</%s> <!-- %s --> ", n.Name, n.Classtr)
}

func (n *Node) Content() (s string) {
	for _, n := range n.Children {
		s += n.Content()
	}
	return s
}

func (d *Data) Content() (s string) {
	return d.content
}
