package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

func (cli *Client) Start(done chan<- bool) {
	reader := bufio.NewReader(os.Stdin)
	running := true
	for running {
		fmt.Print("prompt~> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("failed to read string %+v", err)
			done <- true
			return
		}

		// Remove whitespace around line
		input := strings.Trim(line, " \t\n")
		if input == "" {
			// Ignore blank lines
			continue
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
		case "walk":
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
	done <- true
}

func (cli *Client) CrawlUrls(urls []string) {
	for _, url := range urls {
		page := Crawl(url)
		if page == nil {
			log.Errorln("Crawl failed ", url)
		} else {
			fmt.Printf("%s: %+v\n", page)
		}
		// Page is not alone
	}
}

func (cli *Client) GetHome(url string, args []string) {
	fmt.Fprintln(cli.Writer, "HOME: "+url)
}
