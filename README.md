go-force
======

[![wercker status](https://app.wercker.com/status/66ea433103de60e20ce0f96340a75828/m "wercker status")](https://app.wercker.com/project/bykey/66ea433103de60e20ce0f96340a75828)

[Golang](http://golang.org/) API wrapper for [Force.com](http://www.force.com/), [Salesforce.com](http://www.salesforce.com/)

Installation
============
	go get github.com/nimajalali/go-force/force

Example
============

	package main

	import (
		"fmt"
		"log"
	
		"github.com/nimajalali/go-force/force"
		"github.com/nimajalali/go-force/sobjects"
	)
	
	type SomeCustomSObject struct {
		sobjects.BaseSObject
		Active    bool   `force:"Active__c"`
		AccountId string `force:"Account__c"`
	}
	
	func (t *SomeCustomSObject) ApiName() string {
		return "SomeCustomObject__c"
	}

	func main() {
		// Init the force
		err := force.Init(
			"YOUR-API-VERSION", 
			"YOUR-CLIENT-ID", 
			"YOUR-CLIENT-SECRET", 
			"YOUR-USERNAME", 
			"YOUR-PASSWORD", 
			"YOUR-SECURITY-TOKEN", 
			"YOUR-ENVIRONMENT",
		)
		if err != nil {
			log.Fatal(err)
		}

		someCustomSObject := &SomeCustomSObject{}
		err := force.GetSObject("Your-Object-ID", someCustomSObject)
		if err != nil {
			fmt.Println(err)
		}
		
		fmt.Printf("%#v", someCustomSObject)
	}

Documentation
=======

* [Package Reference](http://godoc.org/github.com/nimajalali/go-force/force)
* [Force.com API Reference](http://www.salesforce.com/us/developer/docs/api_rest/)

TODO
=================
* Write tests for externalId based api calls
* Implement all standard objects. Pull requests welcome.
