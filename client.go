package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	Baseurl string
	Format  string

	io.Reader
	io.Writer
}

func NewClient(baseurl string) (c *Client) {
	c = &Client{
		Baseurl: baseurl,
		Reader:  os.Stdin, // Bydefault, of course
		Writer:  os.Stdout,
		Format:  "text",
	}
	return c
}

func (cli *Client) Get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Errorln("http get ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("ioutil.ReadAll body ", err)
		return
	}
	n, err := fmt.Fprint(cli.Writer, body)
	if err != nil {
		log.Errorln("cli.w.Reader error ", err)
		return
	}
	log.Infof("Wrote %d to reader", n)
}

func (cli *Client) Start() {
	reader := bufio.NewReader(os.Stdin)
	running := true

	for running {

		fmt.Print("prompt~> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("failed to read string %+v", err)
			continue
		}

		// Remove whitespace around line
		input := strings.Trim(line, " \t\n")
		if input == "" {
			continue // ignore blank lines
		}

		// Chop string into an array of words
		cmds := strings.Fields(input)
		if cmds == nil || len(cmds) == 0 {
			// No commands, nothing to do
			log.Warnln("nothing to work with, continuing...")
			continue
		}

		cmd := cmds[0]
		switch cmd {
		case "crawl":
			cli.CrawlUrls(cmds[1:])
		case "home":
			cli.Get("/")
		case "exit":
			fmt.Println("  exiting ..., breaking read loop!")
			running = false
		default:
			fmt.Println("unknown command: ", cmd)
		}
	}
	fmt.Println("client is exiting")
}

func (cli *Client) CrawlUrls(urls []string) {
	r := httpServer().Handler
	for _, url := range urls {

		// Prepare a request
		req, err := http.NewRequest("GET", "/crawl/"+url, nil) // nil is io.Reader (body)
		if err != nil {
			log.Errorln("Client ~ failed to create an http.NewRequest")
			return
		}
		w := httptest.NewRecorder()

		// CrawlHandler is the same function called by the HTTP server!
		// which takes care of sanitizing the URL(s) and other house
		// keeping functions, we will just reuse it from the command line.
		r.ServeHTTP(w, req)
		// get called with vars set properly -> CrawlHandler(w, req)

		// Let Us handle the result
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
	}
}

func (cli *Client) GetHome(url string, args []string) {
	fmt.Fprintln(cli.Writer, "HOME: "+url)
}
