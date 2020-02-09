# Moni the Web Watcher

Moni is a light weight, self contianed (single executable) website
walker and analyzer. It is basically like a "ping" for an entire
website. 

Just point at one or more websites that you manage, moni will grab
every page of your website, record errors it found some, otherwise the
pages response time will also be recorded.

You will be able to tell pretty quickly if there are any problems with
your website, including its responsiveness. Moni was written in go
with concurrency in mind, moni is light, tight and fast.

## What Moni Find for You

- Websites
  - Walk all pages it can find
  - Are all of your web pages accessible?
  - How fast are they responding?
  - What resources do they require (css, javascript and images)?
  - What internal links does this page have?
  - What external links does this page have?
  - Are ANY links broken? 
  - ACL limits access to webpages we are monitoring

### 

Moni has no external dependencies, just provide one or more URLs and
determine if your website is behaving like you expect it to.
 
## How to Use

Crawl is both a _daemon_ and an _angel_ (JK :) and a _command line
tool_.  When run as a command line tool crawl can be used to diagnose
website troubles in real time.

When run as a _daemon_ or _background service_ crawl monitors the
availability and performance of the list of websites.

## Command Line

```bash
% ./moni rustyeddy.com
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

Or watch a group of websites with a quick and easy single line output:

```bash
                  http://rustyeddy.com: resp: 659.999µs links: 7
                  http://oclowvision.com: resp: 986.158µs links: 10
        http://orangecountylowvision.com: resp: 888.039µs links: 10
          http://sierrahydrographics.com: resp: 1.644021ms links: 16
               http://gardenpassages.com: resp: 1.136244ms links: 22
                   http://skillcraft.com: resp: 1.519748ms links: 14
                   http://mobilerobot.io: resp: 1.042628ms links: 7
```

This is handy if you manage a lot of websites, or just want to get a
quick summary of how your web properties are doing..

## Crawling Websites

Moni walks every page of a website, that it can find, and reports back
with the stats gathered above. The output will be produced according
to your desires!  Text or JSON? You got it!

If desired, Moni will also walk every sub-page of the website, however
all external links with NOT be walked otherwise we will quickly be
walking every website in the world, and be the next Google.

We don't want to do that, we just want to ensure the good health of
our own web properties, and find out immediately when something goes
wrong (and it will, this is tech after all).

## Moni is Built on Colly!

Moni uses the go package [colly](http://getcolly.io) to monitor given
websites. Colly makes a standard _get_ request to a given URL, gathers
the links from that web page, as well as text and other elements, then
returns. The return gets a timestamp that we use to record the elapsed
time.

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

