go-force
======
<p align="center">
  <a href="https://goreportcard.com/report/github.com/nimajalali/go-force"><img src="https://goreportcard.com/badge/github.com/nimajalali/go-force" alt="Go Report Card"></a>
  <a href="https://github.com/nimajalali/go-force/actions?query=workflow%3Abuild"><img src="https://github.com/nimajalali/go-force/workflows/build/badge.svg" alt="build status"></a>
  <a href="https://github.com/nimajalali/go-force/blob/master/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/nimajalali/go-force" alt="Go version"></a>
  <a href="https://github.com/nimajalali/go-force/releases"><img src="https://img.shields.io/github/v/release/nimajalali/go-force.svg" alt="Current Release"></a>
  <a href="https://godoc.org/github.com/nimajalali/go-force"><img src="https://godoc.org/github.com/nimajalali/go-force?status.svg" alt="godoc"></a>
  <a href="https://gocover.io/github.com/nimajalali/go-force/force"><img src="https://gocover.io/_badge/github.com/nimajalali/go-force/force" alt="Coverage"></a>
  <a href="https://github.com/nimajalali/go-force/blob/master/LICENSE"><img src="https://img.shields.io/github/license/nimajalali/go-force" alt="License"></a>
</p>

[Golang](http://golang.org/) API wrapper for [Force.com](http://www.force.com/), [Salesforce.com](http://www.salesforce.com/)

Installation
============
	go get github.com/nimajalali/go-force/force

Example
============
```go
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

type SomeCustomSObjectQueryResponse struct {
	sobjects.BaseQuery

	Records []*SomeCustomSObject `force:"records"`
}

func main() {
	// Init the force
	forceApi, err := force.Create(
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

	// Get somCustomSObject by ID
	someCustomSObject := &SomeCustomSObject{}
	err = forceApi.GetSObject("Your-Object-ID", nil, someCustomSObject)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", someCustomSObject)

	// Query
	someCustomSObjects := &SomeCustomSObjectQueryResponse{}
	err = forceApi.Query("SELECT Id FROM SomeCustomSObject__c LIMIT 10", someCustomSObjects)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", someCustomSObjects)
}
```
Documentation 
=======

* [Package Reference](http://godoc.org/github.com/nimajalali/go-force/force)
* [Force.com API Reference](http://www.salesforce.com/us/developer/docs/api_rest/)
