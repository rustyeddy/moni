package moni

import (
	"fmt"
	"net/http"
)

// Element Writer (HTML) ~~~ Interface
// ====================================================================

// Element will become an HTML element
type ElementWriter interface {
	OpenTag(w http.ResponseWriter) error
	CloseTag(w http.ResponseWriter) error
	Content(w http.ResponseWriter) error
}

// Element ~ Concrete HTML implementation of things
// ====================================================================

// Node is an Element that has other Elements for children
type Element struct {
	Name    string
	Id      string
	Classtr string

	children []*Element
	data     string

	openfmt  string
	closefmt string
}

// NewElement spits out an element representing the given parameters.
// Elements with no closing tags can use alternative formats such as:
// openfmt: <%s%s%s /> and closefmt: "", meaning no close format.
// Note: an empty closing tag also implies that this element can NOT
// have any Content, meaning no children OR Data.
func NewElement(name, id string, classes string) Element {
	return Element{
		Name:    name,
		Id:      id,
		Classtr: classes,

		openfmt:  "<%s%s%s>",
		closefmt: "</%s>", // default could be: "" meaning no closing tab
	}
}

// verify this element is only acting as a branch or leaf
func (e *Element) verify() bool {
	if e.children == nil || e.data == "" {
		// at least one of the contents are empty, we are fine
		return true
	}
	log.Fatalln("Element verification failed. Children and Data present. Data xor Children")
	return false
}

// Content will have all content written out to the caller
func (e *Element) Content(w http.ResponseWriter) error {
	e.verify()

	// If we are a Node in the Tree, recursively get the combined content
	// of our children spit out.
	if e.children != nil {
		return e.childContent(w)
	}
	if e.data != "" {
		if _, err := w.Write([]byte(e.data)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Element) childContent(w http.ResponseWriter) error {
	for _, child := range e.children {
		if err := child.OpenTag(w); err != nil {
			return fmt.Errorf("failed content for %v", e)
		}

		if err := child.Content(w); err != nil {
			return fmt.Errorf("failed childContent for Content %v", err)
		}

		if err := child.CloseTag(w); err != nil {
			return fmt.Errorf("failed close tag for %v", e)
		}
	}
	return nil
}

// OpenTag will assemble the opening according to whether one or
// classes exists, or an Id.  Finally write the opening tag to
// the http.ResponseWriter
func (n *Element) OpenTag(w http.ResponseWriter) error {
	tag := "<" + n.Name
	if n.Id != "" {
		tag += " id='" + n.Id + "'"
	}
	if cls := n.Classtr; cls != "" {
		tag += " class='" + cls + "'"
	}
	tag += ">"
	if _, err := w.Write([]byte(tag)); err != nil {
		return fmt.Errorf("failed writing OpenTag %s %v", tag, err)
	}
	return nil
}

// CloseTag finalizes writing this particular element to output.  The content of
// all children (if any should have been output as well)
func (n *Element) CloseTag(w http.ResponseWriter) error {
	if _, err := fmt.Fprintf(w, "</%s> <!-- %s --> ", n.Name, n.Classtr); err != nil {
		return fmt.Errorf("failed to write CloseTag %s %v", n.Name, err)
	}
	return nil
}

// Node is an HTML element with more children, it has NO data
type Node struct {
	Element
	Children []*Node
}

// Data is an HTML element that has data for the end user,
// it has NO children,
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

// Content will Recursively cause all children content to be printed
func (n *Node) Content() (s string) {
	for _, n := range n.Children {
		s += n.Content()
	}
	return s
}

// Content will return the content of this Data object directly
func (d *Data) Content() (s string) {
	return d.content
}
