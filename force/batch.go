package force

import (
	"fmt"
)

var (
	// job states
	abortJob = jobState{"Aborted"}
	closeJob = jobState{"Closed"}
	// TODO: make this part of the client object.
	sfVersion = ""
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
	BaseURI                 string
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

// SBatchInfo is container for the a list of batch status.
type SBatchInfo struct {
	BatchInfo []SBatch `json:"batchInfo"`
}

// CreateJob requests a new bulk job.
func (forceAPI *API) CreateJob(bulkAPIVersion string, req *SJobRequest) (*SJob, error) {

	uri := fmt.Sprintf("/services/async/%s/job", bulkAPIVersion)

	job := &SJob{}
	err := forceAPI.Post(uri, nil, req, job)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new job: %s", err)
	}

	job.bulkAPIVersion = bulkAPIVersion
	sfVersion = bulkAPIVersion

	job.BaseURI = fmt.Sprintf("%s/%s", uri, job.ID)
	job.forceAPI = forceAPI

	forceAPI.openJobMap[job.ID] = *job

	return job, nil
}

// GetSJob return an open SJob.
func (forceAPI *API) GetSJob(id string) (SJob, error) {
	if j, ok := forceAPI.openJobMap[id]; ok {
		return j, nil
	}
	return SJob{}, fmt.Errorf("Job '%s' is not in the list of open jobs", id)
}

// GetOpenSJobs return a map containing the list of open SJobs.
func (forceAPI *API) GetOpenSJobs() map[string]SJob {
	return forceAPI.openJobMap
}

// IsCompleted returns whether a job is completed
// (i.e. no in-progress or queued records).
func (j *SJob) IsCompleted() (bool, error) {

	err := j.forceAPI.Get(j.BaseURI, nil, j)
	if err != nil {
		return false, fmt.Errorf("Failed to get job (%s) state: %s", j.ID, err)
	}
	if _, ok := j.forceAPI.openJobMap[j.ID]; ok {
		j.forceAPI.openJobMap[j.ID] = *j
	}

	if j.NumberBatchesInProgress == 0 && j.NumberBatchesQueued == 0 {
		return true, nil
	}

	return false, nil
}

// IsCompletedByID given a job ID, returns whether a job is completed
// (i.e. no in-progress or queued records).
func (forceAPI *API) IsCompletedByID(ID string) (bool, error) {

	uri := fmt.Sprintf("/services/async/%s/job/%s", sfVersion, ID)
	j := &SJob{}
	err := forceAPI.Get(uri, nil, j)
	if err != nil {
		return false, fmt.Errorf("Failed to check job (%s): %s", ID, err)
	}
	if _, ok := forceAPI.openJobMap[ID]; ok {
		forceAPI.openJobMap[ID] = *j
	}

	if j.NumberBatchesInProgress == 0 && j.NumberBatchesQueued == 0 {
		return true, nil
	}

	return false, nil
}

// AddBatch adds a batch to job.
func (j *SJob) AddBatch(payload interface{}) (*SBatch, error) {
	if !j.IsOpen() {
		return nil, fmt.Errorf("Job '%s' is not in open state", j.ID)
	}
	uri := fmt.Sprintf("%s/batch", j.BaseURI)

	batch := &SBatch{}
	err := j.forceAPI.Post(uri, nil, payload, batch)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new batch: %s", err)
	}
	j.forceAPI.openJobMap[j.ID] = *j

	return batch, nil
}

// IsOpen checks whether the job is open.
func (j *SJob) IsOpen() bool {
	state, _ := j.GetState()

	return state == "Open"
}

// Refresh updates the info of a job.
func (j *SJob) Refresh() error {

	err := j.forceAPI.Get(j.BaseURI, nil, j)
	if err != nil {
		return fmt.Errorf("Failed to refresh job (%s): %s", j.ID, err)
	}
	if _, ok := j.forceAPI.openJobMap[j.ID]; ok {
		j.forceAPI.openJobMap[j.ID] = *j
	}

	return nil
}

// GetState returns the state of a job.
func (j *SJob) GetState() (string, error) {

	err := j.forceAPI.Get(j.BaseURI, nil, j)
	if err != nil {
		return "", fmt.Errorf("Failed to get job (%s) state: %s", j.ID, err)
	}
	if _, ok := j.forceAPI.openJobMap[j.ID]; ok {
		j.forceAPI.openJobMap[j.ID] = *j
	}

	return j.State, nil
}

// Close closes the job.
func (j *SJob) Close() error {
	if !j.IsOpen() {
		return nil
	}
	err := j.forceAPI.Post(j.BaseURI, nil, closeJob, j)
	if err != nil {
		return fmt.Errorf("Failed to close job (%s): %s", j.ID, err)
	}
	delete(j.forceAPI.openJobMap, j.ID)
	return nil
}

// CloseJobByID closes a job by ID.
func (forceAPI *API) CloseJobByID(ID string) error {

	uri := fmt.Sprintf("/services/async/%s/job/%s", sfVersion, ID)

	job := &SJob{}
	err := forceAPI.Post(uri, nil, closeJob, job)
	if err != nil {
		return fmt.Errorf("Failed to close job (%s): %s", ID, err)
	}
	delete(forceAPI.openJobMap, ID)

	return nil
}

// GetBatches returns an array containing the list of batches in a job.
func (j *SJob) GetBatches() (*SBatchInfo, error) {

	URI := fmt.Sprintf("%s/batch", j.BaseURI)

	resp := &SBatchInfo{
		BatchInfo: make([]SBatch, j.NumberRecordsProcessed),
	}

	err := j.forceAPI.Get(URI, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to get list of batches (%s): %s", j.ID, err.Error())
	}

	return resp, nil
}

// BatchState returns an array containing the state of each record in the batch.
func (j *SJob) BatchState(id string) ([]SBatchResponse, error) {
	URI := fmt.Sprintf("%s/batch/%s/result", j.BaseURI, id)

	resp := make([]SBatchResponse, j.NumberRecordsProcessed)

	err := j.forceAPI.Get(URI, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to get batch (%s) state: %s", id, err.Error())
	}

	return resp, nil
}

// GetBatchRecordIDs returns a list of records in a batch. If 'all' is false
// only the successful records are returned. Otherwise all records are returned.
func (j *SJob) GetBatchRecordIDs(ID string, all bool) ([]string, error) {

	var recordList []string
	records, err := j.BatchState(ID)
	if err != nil {
		errStr := "CRM cannot get batch state for batch '%s' in job '%s': %s"
		return nil, fmt.Errorf(errStr, ID, j.ID, err.Error())
	}

	for _, record := range records {
		if all {
			recordList = append(recordList, record.ID)
		} else if record.Success || record.Created {
			recordList = append(recordList, record.ID)
		}
	}
	return recordList, nil
}

// Abort aborts a job.
func (j *SJob) Abort() error {
	err := j.forceAPI.Post(j.BaseURI, nil, abortJob, j)
	if err != nil {
		return fmt.Errorf("Failed to abort job (%s): %s", j.ID, err)
	}

	return nil
}
