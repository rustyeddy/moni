package main

/*
Access List (ACL) is a structure that wraps a map of URLs with boolean
values indicating if that URL will be scanned.  We want to avoid doing
long walks on large external websites like google, github and so on.
*/

// ACL is the AccessControl List that
type ACL struct {
	AccessList map[string]bool

	Allowed  int
	Rejected int
}

// NewACL creates a new access control list
func init() {
	acl = ACL{make(map[string]bool), 0, 0}

	// reject these websites out right
	acl.Add("google.com", false)
	acl.Add("github.com", false)
	acl.Add("wpengine.com", false)

	// accept these websites
	acl.Add("gumsole.com", true)

	acl.Add("rustyeddy.com", true)
	acl.Add("oclowvision.com", true)
	acl.Add("mobilerobot.io", true)

	acl.Add("sierrahydrographics.com", true)
	acl.Add("gardenpassages.com", true)
}

// Add adds another rule to the ACL
func (a *ACL) Add(urlstr string, allow bool) {
	a.AccessList[urlstr] = allow
}

// Allow returns true if the url exists and is set to true, otherwise
// false is returned
func (a *ACL) Allow(urlstr string) bool {
	allow, exists := a.AccessList[urlstr]
	if exists == false || allow == false {
		a.Rejected++
		return false
	}
	a.Allowed++
	return true
}
