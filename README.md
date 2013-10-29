go-force
======

[Golang](http://golang.org/) API wrapper for [Force.com](http://www.force.com/), [Salesforce.com](http://www.salesforce.com/)

Installation
============
	go get github.com/nimajalali/go-force/force

Example
============

	package main

	import (
		"fmt"
		"github.com/nimajalali/go-force/force"
	)

	func main() {
		// Initialize with your Login information
		force.Init("YOUR-API-VERSION", "YOUR-CLIENT-ID", "YOUR-CLIENT-SECRET", "YOUR-USERNAME", "YOUR-PASSWORD", "YOUR-SECURITY-TOKEN", "YOUR-ENVIRONMENT")

		
	}

Documentation
=======

* [Package Reference](http://godoc.org/github.com/nimajalali/go-force/force)
* [Force.com API Reference](http://www.salesforce.com/us/developer/docs/api_rest/)

TODO
=================
* Write tests for externalId based api calls
* Implement all standard objects. Pull requests welcome.
