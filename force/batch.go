package force

import (
	"fmt"
)

const (
	jsonType = "JSON"
)

var (
	// job states
	abortJob = jobState{"Aborted"}
	closeJob = jobState{"Closed"}
)

type jobState struct {
	State string `force:"state"`
}

// BaseRequest is the interface for job and batch requests.
type BaseRequest interface {
	GetID() string
}

// SJobRequest represents an job request.
type SJobRequest struct {
	baseRequest         BaseRequest
	id                  string
	Operation           string `force:"operation"`
	Object              string `force:"object"`
	ContentType         string `force:"contentType"`
	ExternalIDFieldName string `force:"externalIdFieldName"`
}

// BaseResponse is the interface for job and batch responses.
type BaseResponse interface {
	GetID() string
	GetState() string
}

// SJob represents a response for a given job request.
type SJob struct {
	BaseResponse
	ApexProcessingTime      int     `json:"apexProcessingTime"`
	APIVersion              float64 `json:"apiVersion"`
	APIActiveProcessingTime int     `json:"apiActiveProcessingTime"`
	AssignmentRuleID        string  `json:"assignmentRuleId"`
	baseURI                 string
	bulkAPIVersion          string
	ConcurrencyMode         string `json:"concurrencyMode"`
	ContentType             string `json:"contentType"`
	CreatedByID             string `json:"createdById"`
	CreatedDate             string `json:"createdDate"`
	ExternalIDFieldName     string `json:"externalIdFieldName"`
	FastPathEnabled         bool   `json:"fastPathEnabled"`
	forceAPI                *API
	ID                      string `json:"id"`
	NumberBatchesCompleted  int    `json:"numberBatchesCompleted"`
	NumberBatchesFailed     int    `json:"numberBatchesFailed"`
	NumberBatchesInProgress int    `json:"numberBatchesInProgress"`
	NumberBatchesQueued     int    `json:"numberBatchesQueued"`
	NumberBatchesTotal      int    `json:"numberBatchesTotal"`
	NumberRecordsFailed     int    `json:"numberRecordsFailed"`
	NumberRecordsProcessed  int    `json:"numberRecordsProcessed"`
	NumberRetries           int    `json:"numberRetries"`
	Object                  string `json:"object"`
	Operation               string `json:"operation"`
	State                   string `json:"state"`
	SystemModstamp          string `json:"systemModstamp"`
	TotalProcessingTime     int    `json:"totalProcessingTime"`
}

// SBatch represents a batch of records.
type SBatch struct {
	ApexProcessingTime      int    `json:"apexProcessingTime"`
	APIActiveProcessingTime int    `json:"apiActiveProcessingTime"`
	CreatedDate             string `json:"createdDate"`
	ID                      string `json:"id"`
	JobID                   string `json:"jobId"`
	NumberRecordsFailed     int    `json:"numberRecordsFailed"`
	NumberRecordsProcessed  int    `json:"numberRecordsProcessed"`
	State                   string `json:"state"`
	SystemModstamp          string `json:"systemModstamp"`
	TotalProcessingTime     int    `json:"totalProcessingTime"`
}

// SBatchResponse contains the batch results.
type SBatchResponse struct {
	ID      string    `json:"id"`
	Errors  APIErrors `json:"errors"`
	Success bool      `json:"success"`
	Created bool      `json:"created"`
}

// CreateJob requests a new bulk job.
func (forceAPI *API) CreateJob(bulkAPIVersion, operation, object, externalIDField string) (*SJob, error) {
	jobReq := &SJobRequest{
		Operation:           operation,
		Object:              object,
		ContentType:         jsonType,
		ExternalIDFieldName: externalIDField,
	}

	uri := fmt.Sprintf("/services/async/%s/job", bulkAPIVersion)

	job := &SJob{}
	err := forceAPI.Post(uri, nil, jobReq, job)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new job: %s", err)
	}

	job.bulkAPIVersion = bulkAPIVersion
	job.baseURI = fmt.Sprintf("%s/%s", uri, job.ID)
	job.forceAPI = forceAPI

	return job, nil
}

// jobIsOpen checks whether the job is open.
func (j *SJob) jobIsOpen() error {
	if j.State != "Open" {
		return fmt.Errorf("Job '%s' is in '%s' state", j.ID, j.State)
	}
	return nil
}

// AddBatch adds a batch to job.
func (j *SJob) AddBatch(payload interface{}) (*SBatch, error) {

	if err := j.jobIsOpen(); err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s/batch", j.baseURI)

	batch := &SBatch{}
	err := j.forceAPI.Post(uri, nil, payload, batch)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new batch: %s", err)
	}

	return batch, nil
}

// GetState returns the state of a job.
func (j *SJob) GetState() (string, error) {

	if err := j.jobIsOpen(); err != nil {
		return "", err
	}

	job := &SJob{}

	err := j.forceAPI.Get(j.baseURI, nil, job)
	if err != nil {
		return "", fmt.Errorf("Failed to get job (%s) state: %s", j.ID, err)
	}
	// update the reference to the job
	j = job

	return j.State, nil
}

// Close closes the job.
func (j *SJob) Close() error {

	job := &SJob{}
	err := j.forceAPI.Post(j.baseURI, nil, closeJob, job)
	if err != nil {
		return fmt.Errorf("Failed to close job (%s): %s", j.ID, err)
	}
	// Update the reference
	j = job

	return nil
}

// BatchState returns an array containing the state of each record in the batch.
func (j *SJob) BatchState(id string) ([]SBatchResponse, error) {
	if err := j.jobIsOpen(); err != nil {
		return nil, err
	}
	URI := fmt.Sprintf("%s/batch/%s/result", j.baseURI, id)

	resp := make([]SBatchResponse, j.NumberRecordsProcessed)

	err := j.forceAPI.Get(URI, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to get batch (%s) state: %s", id, err.Error())
	}

	return resp, nil
}

// Abort aborts a job.
func (j *SJob) Abort() error {
	job := &SJob{}
	err := j.forceAPI.Post(j.baseURI, nil, abortJob, job)
	if err != nil {
		return fmt.Errorf("Failed to abort job (%s): %s", j.ID, err)
	}
	// Update the reference
	j = job

	return nil
}
