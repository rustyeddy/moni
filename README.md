# Crawl the Website Babby Sitter

Crawl is a light weight website babby sitter.  It constantly monitors
your website(s) for both reliability and performance.

Crawl is written go (very small and fast with single executable) using
the [colly](http://gocolly.io) web scraping package.

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


