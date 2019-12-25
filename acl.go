package main

func init() {
	acl = make(map[string]bool)

	// reject
	acl["localhost"] = false
	acl["google.com"] = false
	acl["github.com"] = false
	acl["wpengine.com"] = false

	// accept
	acl["gumsole.com"] = true

	acl["rustyeddy.com"] = true
	acl["oclowvision.com"] = true
	acl["mobilerobot.io"] = true

	acl["sierrahydrographics.com"] = true
	acl["gardenpassages.com"] = true
}
