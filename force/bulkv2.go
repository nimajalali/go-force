//
// Salesforce BULK API V 2.0 is a version of the Buklk API that uses REST
// endpoints and only accepts CSV data. It is an alternative to the normal
// API and misses many features.
//
// This package has not yet been extensively tested against salesforce api.
// i wrote it as a sane replacement for the current bulk api, this version uses
// Salesforce BULK API V 2.0 which runs against the Salesforce REST API and
// expects CSV as input.

package force

import (
	"fmt"
	"net/http"
	"strings"
)

// MaxBatchv2SizeMiB is the maximum a batch can hold
const MaxBatchv2SizeMiB = 1024 * 1024 * 100

// SaveBatchv2SizeMiB a bit smaller than MaxBatchSizeMiB, to be on the save side
const SaveBatchv2SizeMiB = 1024 * 1024 * 80

// Bulkv2JobReq is used to create a new Bulk API 2.0 Job.
type Bulkv2JobReq struct {
	ColumnDelimiter     string      `force:"columnDelimiter,omitempty"`     // The column delimiter used for CSV job data. The default value is COMMA.
	ContentType         ContentType `force:"contentType,omitempty"`         // The content type for the job. The only valid value (and the default) is CSV.
	ExternalIdFieldName string      `force:"externalIdFieldName,omitempty"` // The external ID field in the object being updated. Only needed for upsert operations. Field values must also exist in CSV job data.
	LineEnding          LineEnding  `force:"lineEnding,omitempty"`          // The line ending used for CSV job data, marking the end of a data row. The default is LF. Valid values are:
	Object              string      `force:"object"`                        // The object type for the data being processed. Use only a single object type per job.
	Operation           Operation   `force:"operation"`                     // The processing operation for the job. Valid values are: insert delete update upsert
}

// BulkV2Job as represented by the API.
type BulkV2Job struct {
	API        *API
	ApiVersion float32 `force:"apiVersion,omitempty"` // The API version that the job was created in.

	ID             string    `force:"id,omitempty"`             // Unique ID for this job.
	Operation      Operation `force:"operation,omitempty"`      // The processing operation for the job. Values include: insert delete update upsert
	Object         string    `force:"object,omitempty"`         // The object type for the data being processed.
	CreatedById    string    `force:"createdById,omitempty"`    // The ID of the user who created the job.
	CreatedDate    string    `force:"createdDate,omitempty"`    // The date and time in the UTC time zone when the job was created.
	SystemModstamp string    `force:"systemModstamp,omitempty"` // Date and time in the UTC time zone when the job finished.

	State               JobState        `force:"state,omitempty"`               // The current state of processing for the job.
	ExternalIdFieldName string          `force:"externalIdFieldName,omitempty"` // The name of the external ID field for an upsert.
	ConcurrencyMode     ConcurrencyMode `force:"concurrencyMode,omitempty"`     // The concurrency mode for the job.
	ContentType         string          `force:"contentType,omitempty"`         // The format of the data being processed. In Bulk API v2.0 only CSV is supported.
	ContentUrl          string          `force:"contentUrl,omitempty"`          // The URL to use for Upload Job Data requests for this job. Only valid if the job is in Open state.

	ColumnDelimiter ColumnDelimiter `force:"columnDelimiter,omitempty"` // The column delimiter used for CSV job data.
	LineEnding      LineEnding      `force:"lineEnding,omitempty"`      // The line ending used for CSV job data. Values include: LF—linefeed character CRLF—carriage return character followed by a linefeed character

	NumberRecordsFailed    int `force:"numberRecordsFailed,omitempty"`    // The number of records that were not processed successfully in this job.
	NumberRecordsProcessed int `force:"numberRecordsProcessed,omitempty"` // The number of records already processed.

	JobType                 JobType `force:"jobType,omitempty"`                  // The job’s type. Values include: BigObjectIngest—BigObjects job Classic—Bulk API 1.0 job V2Ingest—Bulk API 2.0 job
	ApiActiveProcessingTime int64   `force:"apiActiveProcessingTime,omitempty"`  // The number of milliseconds taken to actively process the job and includes apexProcessingTime
	ApexProcessingTime      int64   `force:"apexProcessingTime,omitempty"`       // The number of milliseconds taken to process triggers and other processes related to the job data.
	TotalProcessingTime     int64   `force:"totaforcelProcessingTime,omitempty"` // The number of milliseconds taken to process the job.
	Retries                 int     `force:"retries,omitempty"`                  // The number `json:" of times that Salesforce attempted to save the results of an operation.

}

// CreateBulkv2Job creates an Bulk API 2.0 Job.
func (forceAPI *API) CreateBulkv2Job(sObjectname string, operation Operation) (*BulkV2Job, error) {
	uri := fmt.Sprintf(`/services/data/%s/jobs/ingest`, forceAPI.apiVersion)

	req := &Bulkv2JobReq{
		ContentType: ContentTypeCSV,
		Object:      sObjectname,
		Operation:   operation,
	}

	var resp BulkV2Job
	err := forceAPI.Post(uri, nil, nil, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("Bulk API 2.0: Can't create new Job for '%s': %s", sObjectname, err)
	}
	resp.API = forceAPI

	return &resp, nil
}

// BulkV2JobList us used to unmarshal a list of BulkV2Jobs.
type BulkV2JobList struct {
	Done           bool        `json:"done"`
	NextRecordsUrl string      `json:"nextRecordsUrl"`
	Records        []BulkV2Job `json:"records"`
}

// GetAllBulkv2Jobs retrieves all jobs in the org
func (forceAPI *API) GetAllBulkv2Jobs() (resp []*BulkV2Job, err error) {
	uri := fmt.Sprintf(`/services/data/%s/jobs/ingest`, forceAPI.apiVersion)

	var jobList BulkV2JobList
	for {
		err = forceAPI.Get(uri, nil, nil, &jobList)
		if err != nil {
			return nil, fmt.Errorf("Bulk API 2.0: Can't retireve info for all Jobs: %s", err)
		}

		for i := range jobList.Records {
			resp = append(resp, &jobList.Records[i])
		}

		if jobList.Done {
			break
		}
		uri = strings.TrimPrefix(jobList.NextRecordsUrl, forceAPI.oauth.InstanceURL)
	}

	return
}

// AddBatch uploads a data batch of up to 150MB.
//
// A request can provide CSV data that does not in total exceed 150 MB of base64 encoded content.
// When job data is uploaded, it is converted to base64. This conversion can increase the data
// size by approximately 50%. To account for the base64 conversion increase, upload data that does not exceed 100 MB.
func (b *BulkV2Job) AddBatch(batch []byte) (err error) {
	uri := fmt.Sprintf(`/services/data/v%.1f/jobs/ingest/%s/batches`, b.ApiVersion, b.ID)

	headers := http.Header{"Content-Type": []string{"text/csv"}}

	err = b.API.Put(uri, nil, headers, string(batch), nil)
	if err != nil {
		return fmt.Errorf("Bulk API 2.0: Can't add Batch to Job '%s': %s", b.ID, err)
	}
	return nil
}

// Close the job, starts executing the batches.
func (b *BulkV2Job) Close() (err error) {
	uri := fmt.Sprintf(`/services/data/v%.1f/jobs/ingest/%s`, b.ApiVersion, b.ID)

	req := fmt.Sprintf(`{"state":"%s"`, JobStateUploadComplete)

	err = b.API.Patch(uri, nil, nil, req, b)
	if err != nil {
		return fmt.Errorf("Bulk API 2.0: Can't close Job '%s': %s", b.ID, err)
	}
	return nil
}

// Abort the job.
func (b *BulkV2Job) Abort() (err error) {
	uri := fmt.Sprintf(`/services/data/v%.1f/jobs/ingest/%s`, b.ApiVersion, b.ID)

	req := fmt.Sprintf(`{"state":"%s"}`, JobStateAborted)

	err = b.API.Patch(uri, nil, nil, req, b)
	if err != nil {
		return fmt.Errorf("Bulk API 2.0: Can't abort Job '%s': %s", b.ID, err)
	}
	return nil
}

// Delete the job.
func (b *BulkV2Job) Delete() (err error) {
	uri := fmt.Sprintf(`/services/data/v%.1f/jobs/ingest/%s`, b.ApiVersion, b.ID)

	err = b.API.Delete(uri, nil, nil)
	if err != nil {
		return fmt.Errorf("Bulk API 2.0: Can't abort Job '%s': %s", b.ID, err)
	}
	return nil
}

// Info of the current job.
func (b *BulkV2Job) Info() (err error) {
	uri := fmt.Sprintf(`/services/data/v%.1f/jobs/ingest/%s`, b.ApiVersion, b.ID)

	err = b.API.Get(uri, nil, nil, b)
	if err != nil {
		return fmt.Errorf("Bulk API 2.0: Can't retrieve Info for Job '%s': %s", b.ID, err)
	}
	return nil
}

// SuccessfulResults of job execution.
//
// Returns a CSV with all successful records.
// CSV fields are as provided in the original job data.
// Except two additional fields: "sf__Id" and "sf__Created"
func (b *BulkV2Job) SuccessfulResults() (resp []byte, err error) {
	uri := fmt.Sprintf(`/services/data/v%.1f/jobs/ingest/%s/successfulResults/`, b.ApiVersion, b.ID)

	err = b.API.Get(uri, nil, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("Bulk API 2.0: Can't retrieve Successful Results for Job '%s': %s", b.ID, err)
	}
	return resp, nil
}

// FailedResults of job execution.
//
// Returns a CSV with all failed records.
// CSV fields are as provided in the original job data.
// Except two additional fields: "sf__Id" and "sf__Created"
func (b *BulkV2Job) FailedResults() (resp []byte, err error) {
	uri := fmt.Sprintf(`/services/data/v%.1f/jobs/ingest/%s/failedResults/`, b.ApiVersion, b.ID)

	err = b.API.Get(uri, nil, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("Bulk API 2.0: Can't Failed Results for Job '%s': %s", b.ID, err)
	}
	return resp, nil
}

// UnprocessedRecords of job execution.
//
// A job that is interrupted or otherwise fails to complete can result in rows that aren’t processed.
// Unprocessed rows are not the same as failed rows. Failed rows are processed but encounter an error during processing.
//
// Returns a CSV with all failed records.
// CSV fields are as provided in the original job data.
func (b *BulkV2Job) UnprocessedRecords() (resp []byte, err error) {
	uri := fmt.Sprintf(`/services/data/v%.1f/ingest/%s/unprocessedrecords/`, b.ApiVersion, b.ID)

	err = b.API.Get(uri, nil, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("Bulk API 2.0: Can't retrieve Unprocessed Records for Job '%s': %s", b.ID, err)
	}
	return resp, nil
}
