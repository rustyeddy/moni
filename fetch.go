package moni

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type PageResponse struct {
	URL string
	*http.Response
	Err     error
	Content []byte
}

func FetchURL(url string, ch chan<- string, processDoc func(li []string, doc *html.Node) []string) {
	var resp *http.Response
	var err error

	start := time.Now()
	if resp, err = http.Get(url); err != nil {
		ch <- err.Error()
		return
	}
	defer resp.Body.Close()

	// Here we go, parse the body now we have find all links
	// and queue them up for a query.
	doc, err := html.Parse(resp.Body)
	if err != nil {
		ch <- err.Error()
	}

	// Should this be a channel?
	log.Debugf("Fetched url %s starting to process ", url)

	rlist := processDoc(nil, doc)
	if rlist != nil {
		log.Errorf("  Left over list %+v\n", rlist)
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("\t%.2fs %s", secs, url)
}

func CheckSites(sites []string, procDoc func(li []string, doc *html.Node) []string) {
	start := time.Now()

	// url channel
	uch := make(chan string)

	// Send the first batch, create a channel back talk
	for _, url := range sites {
		go FetchURL(url, uch, procDoc)
	}

	// Wait for all channels to have be returned, before moving
	// to the next batch.
	for range sites {
		log.Println(<-uch)
	}
	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	time.Sleep(time.Duration(1000) * 10000)
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func outline(stack []string, n *html.Node) []string {
	fmt.Printf("NODE: %+v", n)
	switch n.Type {
	case html.ErrorNode:
		fmt.Println("error node")

	case html.TextNode:
		if n.Data != "" {
			fmt.Println("text data: ", n.Data)
		}

	case html.DocumentNode:
		fmt.Printf("document node: %+v\n", n)

	case html.ElementNode:
		stack = append(stack, n.Data) //  push tag
		if n.Data != "script" {

			// Lets not print big scripts...
			fmt.Printf("elemet: %+v\n", n)
			fmt.Println("element: ", stack)
		}

	case html.CommentNode:
		if n.Data != "" {
			fmt.Println("  comment: ", n.Data)
		}

	case html.DoctypeNode:
		fmt.Printf("doctype node: %+v\n", n)

	default:
		fmt.Printf("unhandled: %+v\n", n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
	return stack
}
