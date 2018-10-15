# Monti is a Site Management and Monitor Tool

## What does Monti do?

Monti is a Site Inventory and Monitoring tool.  

By _Site_ I basically mean _Website_ and all the compute and network
infrastructure to support the production site.

By _Inventory_ I mean *lots of things*, all the things that go into
running a production *site*.

Even the smallest of sites (a modest WordPress install) have lots of
things that need to be keep track of, as we'll see in the list below. 

Here are some of the items that need to be tracked (and monitored) to
ensure the wholistic health of the *site*.


## Inventory Management 

- domain 
  - registerar 
  - nameserver
  - account info (user / passwd keep separate)
	- host names and sub-domains
    - expiration
	- owner / contact

  - hosts & servers
    - host ip
	- services: ports and applications

  - apps & services
    website: example.com
	database: db.example.com
	log: log.example.com

## Monitoring

- Site availablility ~ time series of uptimes
- Site performance ~ timeseries walk every page recording responses
- Site map checking ~ all pages accounted for?  Any unexpected page or
  pages show up? 
- Content scanning ~ ensure all content is as expected, nothing new,
  nothing lost
- Maleware scanning ~ part of scanning

- Monitor Google Analytics
- Monitor Email campaign

- Monitor Shopping Carts

- Monitor ...


Inv can answer the following questions:

- Websites (Every page of Website)?
- Are my web pages all accessible and fast?
- Has my website been hacked and returning junk?
- Are all my networked computers alive?
- Services are they running and accessible?
- What *exactly* is connected to my network?
- How many computers and servers do I have?
- How many devices do I have on my network that are NOT computers?
- What are these devices doing?
