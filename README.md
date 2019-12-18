# Crawl the Website Babby Sitter

Crawl is a light weight website babby sitter.  It constantly monitors
your website(s) for both reliability and performance. Crawl is
written go (very small and fast) as a single executable with no
required external dependencies (just copy the program and use it.)

Crawl uses the excellent _Go package_[colly](http://gocolly.io) to do
most of the hard work.

## How to Use

Crawl is both a _daemon_ and an _angel_ (JK :) and a _command line
tool_.  When run as a command line tool crawl can be used to diagnose
website troubles in real time.

When run as a _daemon_ or _background service_ crawl monitors the
availability and performance of the list of websites.

## Command Line

```bash
% ./crawl rustyeddy.com
http://rustyeddy.com
	https://rustyeddy.com/projects
	https://rustyeddy.com/contact
	https://rustyeddy.com/projects/crawl
	https://rustyeddy.com/
	https://rustyeddy.com/interview
	https://rustyeddy.com/notes
	https://rustyeddy.com/resume
  elapsed time 341.752345ms
```

## Service

As a daemon Crawl provides the following REST API

- POST /api/crawl/_{url}_	= Add the URL to the watch list
- GET  /api/crawl/_{url}_   = Get current watch details
- GET  /api/sites           = Get list of sites on watch list
- GET  /api/site/_{url}_    = Get information about specific site
- GET  /api/healthcheck		= Are we alive?

## Crawling Websites

Crawl walk a given web page (URL) gather all links from the page, it
will very that each link of the given page is reachable, and if so how
responsive is it?

This allows you to quickly determine if something in your web has
become unreachable, or unacceptably slow.

## How it Works

List of crawlable website is loaded into memory, then each of these
sites is periodically crawled and compared to previous crawls.

### Crawl with Colly

Crawl uses the go package colly to monitor given websites. Colly makes
a standard _get_ request to a given URL, gathers the links from that
web page, as well as text and other elements, then returns. The return
gets a timestamp that we use to record the elapsed time.

#### Internal vs External Links

Each link of the webpage is now categorized as either _internal_ or
_external_, internal links may also be crawled, where as external links
will not be _crawled_.

_Internal_ and _External_ links will be tested for reachability.

##### Internal Links

Are links from one page of a given _website_ to another page of the
_same_ website.  We define two web pages being on the the _same
website_ when:

> Two page belong to the same website when their respective
> [canonical] URL's share the same _root domain_.

For simplicity we defer the definition of _root domain_ to DNS, and
take a rather ignorant approach of simply comparing host portion of
the URL.

#### Crawl vs Reachability

When we **crawl** a page we gather all of its contents, in particular
all links on a given page, as well as content.  We then _optionally_
will also _crawl_ any internal links and links to permitted external
sites. 

> The reason for blocking un-allowed external sites is to avoid
> walking sites like amazon.com and github.com just because the page
> has a link to it.

**Reachability** on the other hand is much lighter weight, we simply
verify the page is still online and available. If not, it will be
flagged as a _broken link_.

We do not collect any content from the Reachability test and hence we
do not have any additional links to crawl.

## Software Details

In support of our Crawl capability this software has a _REST API
service_ and a periodic website walking service.

Both services run a Go routines, both run forever in there own little
space. The program is synchronized using a _wait group_.
