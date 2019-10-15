# Monti is a Site Monitoring Tool

Moni is a Site Monitoring tool.  I use it to track the all the
websites that I am responsible for managing.  It keeps "near" real
time stats on most important things regarding the sites _health_ and
_performance_, it has a dashboard that will alert me of anything
requiring immediate or eventual attention.

## 2 Phases ~ Collection and Analysis

For a list of _sites_ give moni, each site will be run through both
the collection phase and analysis phase.

### Collection is responsible for:

The collection phase is mostly mechanical, though the details are
rather complex, and the structure of websites is mostly
non-deterministic, there is not much difference in how collections are
performed. 

1. Determining what to collect (and what NOT to collect)
2. Scheduling (and dampening) collections requests
3. Performing the collections
4. Storaging collected data (if desirable)

### Analysis Can do various things, like

Analysis however, has an endless array of possibilities, not the least
of which are some pretty obvious benefits from walking or "monitoring"
sites, including (but not limited to):

1. Measure performance of pages and all interactions that result in
   the delivery of a single web page.

2. Verify accessability of links

3. Search for required items (Google Analytics, Optin Forms, etc.) or
   unwanted items (anything we did not put in).

## Informal Description of the API

### Sites

 - GET		/site
 - GET		/site/{url}
 - POST		/site/{url}
 - DELETE	/site/{url}

Currently, when a page is submitted to be monitored, it is
automagically added to the ACL.

### Pages

 - GET		/page
 - GET		/page/{url}
 - POST		/page/{url}
 - DELETE	/page/{url}

### Access Lists 

Access Lists control the sites that will be crawled, and a list of
sites that will not be crawled.  Sites may include subdomain
(sub-sites?), 

For example you may want to monitor *example.com* but want to ignore
everything on warriorlist.com.  AccessLists provide that ability.

 - GET		/acl
 - GET		/acl/{url}
 - POST		/acl/{url}
 - DELETE	/acl/{url}

