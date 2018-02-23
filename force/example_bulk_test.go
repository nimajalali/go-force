package force_test

import (
	"time"

	"github.com/prometheus/common/log"
	"github.com/simplesurance/go-force/force"
	"github.com/simplesurance/go-force/forcejson"
	"github.com/simplesurance/go-force/sobjects"
)

const (
	restAPIVersion    = `v40.0`
	oauthClientID     = `0000`
	oauthClientSecret = `1111`
	username          = `test@salesforce.com`
	pssword           = `secret`
	environemnt       = `sandbox`
)

type Event struct {
	sobjects.BaseSObject

	ActivityDateTime  string
	Description       string
	DurationInMinutes int
	Subject           string
	Type              string
	WhatId            string
	WhoId             string
	OwnerId           string
}

func (e Event) APIName() string {
	return "Event"
}

func Example_Bulk() {

	forceAPI, err := force.Create(
		restAPIVersion,
		oauthClientID,
		oauthClientSecret,
		username,
		pssword,
		"",
		environemnt,
		[]string{},
	)
	if err != nil {
		log.Fatal(err)
	}

	data := []Event{
		{
			Subject: `some event`,
			Type:    `Call`,
		},
		{
			Subject: `some other event`,
			Type:    `Meeting`,
		},
	}

	jsonData, err := forcejson.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	job, err := forceAPI.CreateBulkJob(
		data[0].APIName(),
		"",
		force.ContentTypeJSON,
		force.OperationInsert,
		force.ConcurrencyModeSerial)
	if err != nil {
		log.Fatal(err)
	}

	batch, err := job.AddBatch(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	err = job.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		err = batch.Info()
		if err != nil {
			log.Fatal(err)
		}
		switch batch.State {
		case force.BatchStateCompleted:
			if batch.NumberRecordsFailed > 0 {
				log.Errorf("some records failed during processing: %d",
					job.NumberRecordsFailed)
			}
		case force.BatchStateFailed:
			log.Error("job failed")
		}

		time.Sleep(time.Second * 15)
	}
}
