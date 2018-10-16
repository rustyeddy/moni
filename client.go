package main

import (
	"bufio"
	"encoding/json"
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

func NewClient(wr io.Writer, baseurl string) (c *Client) {
	c = &Client{
		Baseurl: baseurl,
		Reader:  os.Stdin, // Bydefault, of course
		Writer:  wr,
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
		args := strings.Fields(input)
		if args == nil || len(args) == 0 {
			// No commands, nothing to do
			log.Warnln("nothing to work with, continuing...")
			continue
		}

		// Both of these need to be checked
		verb := strings.ToLower(args[0])
		switch verb {
		case "get", "put", "delete", "post":
		default:
			log.Errorln("unsupported HTTP verb ", verb)
			continue
		}
		url := args[1]
		resp := cli.Do(verb, url)
		fmt.Printf("resp: %+v\n", resp)
	}
}

func (cli *Client) Do(cmd string, url string) (resp *http.Response) {
	r := httpServer().Handler

	// Prepare a request
	req, err := http.NewRequest(cmd, url, nil)
	if err != nil {
		log.Errorln("Client ~ failed to create an http.NewRequest")
		return nil
	}
	w := httptest.NewRecorder()

	// CrawlHandler is the same function called by the HTTP server!
	// which takes care of sanitizing the URL(s) and other house
	// keeping functions, we will just reuse it from the command line.
	r.ServeHTTP(w, req)

	// get called with vars set properly -> CrawlHandler(w, req)
	// Let Us handle the result
	resp = w.Result()
	if resp == nil {
		log.Errorln("failed to get a response")
		return nil
	}
	return resp
}

func GetBody(resp *http.Response) (b []byte) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("reading http response body %v", err)
		return nil
	}
	return body
}

func (cli *Client) CrawlUrl(url string) {
	resp := cli.Do("post", "/crawl/"+url)
	body := GetBody(resp)

	// TODO ~ Turn this into some pretty print stuff
	// Accept a writer
	fmt.Fprintln(cli.Writer, "Crawl URL ", url)
	fmt.Fprintln(cli.Writer, "  status ", resp.StatusCode)
	fmt.Fprintln(cli.Writer, resp.Header.Get("Content-Type"))
	fmt.Fprintln(cli.Writer, string(body))
}

func (cli *Client) CrawlList() (cl []string) {
	resp := cli.Do("get", "/crawlids")
	body := GetBody(resp)

	err := json.Unmarshal(body, &cl)
	if err != nil {
		log.Errorf("failed unraveling JSON in CrawlList %+v ", err)
		return nil
	}

	fmt.Fprintln(cli.Writer, "Recent CrawlIds ")
	for _, idx := range cl {
		fmt.Fprintln(cli.Writer, "\t", idx)
	}
	return cl
}

func (cli *Client) CrawlId(cid string) (p *Page) {
	var pg *Page
	resp := cli.Do("get", "/crawlid/"+cid)
	body := GetBody(resp)

	err := json.Unmarshal(body, pg)
	if err != nil {
		log.Errorf("failed unraveling JSON in CrawlId %+v ", err)
		return nil
	}

	fmt.Fprintln(cli.Writer, "Crawl ", cid)
	fmt.Fprintf(cli.Writer, "%+v\n", pg)
	return pg
}

func (cli *Client) GetHome(url string, args []string) {
	fmt.Fprintln(cli.Writer, "HOME: "+url)
}
