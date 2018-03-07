//
// Salesforce BULK API is used to upload larger data sets, it supports CSV,
// XML and JSON as input format. It consists of jobs which contain one ore more
// data batches.
//
// Batches limits are not taken care of by the package, users are responsible
// to check the MaxBatch constants against there batch contents
package force

import (
	"fmt"
	"strings"
)

const (
	// MaxBatchSizeInRecords is the maimum a batch can hold
	MaxBatchSizeInRecords = 1000

	// SaveBatchSizeInRecords a bit smaller than MaxBatchSizeMiB,
	// to be on the save side
	SaveBatchSizeInRecords = 5000

	// MaxBatchSizeMiB is the maimum a batch can hold
	MaxBatchSizeMiB = 1024 * 1024 * 10

	// SaveBatchSizeMiB a bit smaller than MaxBatchSizeMiB,
	// to be on the save side
	SaveBatchSizeMiB = 1024 * 1024 * 8

	// MaxBatchCharsPerBatch a batch can contain a maximum of
	// 10,000,000 characters for all the data in a batch.
	MaxBatchCharsPerBatch = 10000000

	// MaxBatchCharsPerRecord a record can contain a maximum
	// of 400,000 characters for all its fields.
	MaxBatchCharsPerRecord = 400000

	// MaxBatchCharsPerField a field can contain a maximum of
	// 32,000 characters.
	MaxBatchCharsPerField = 32000

	// MaxBatchFields a record can contain a maximum of
	// 5,000 fields.
	MaxBatchFields = 5000
)

// ContentType within the Batch
type ContentType string

const (

	// ContentTypeCSV data in CSV format (default and only supported content type for Bulk V2 type jobs)
	ContentTypeCSV ContentType = `CSV`

	// ContentTypeJSON data in JSON format
	ContentTypeJSON ContentType = `JSON`

	// ContentTypeXML data in XML format (default option for Bulk V1 type jobs)
	ContentTypeXML ContentType = `XML`

	// ContentTypeZIPXML data in CSV format in a zip file containing binary attachments
	ContentTypeZIPXML ContentType = `ZIP_XML`

	// ContentTypeZIPJSON data in JSON format in a zip file containing binary attachments
	ContentTypeZIPJSON ContentType = `ZIP_JSON`

	// ContentTypeZIPCSV data in XML format in a zip file containing binary attachments
	ContentTypeZIPCSV ContentType = `ZIP_CSV`
)

// ColumnDelimiter mode for ContentTypeCSV
type ColumnDelimiter string

const (
	// The column delimiter used for CSV job data.
	// The default value is COMMA. Valid values are:
	//
	// BACKQUOTE—backquote character (`)
	// CARET—caret character (^)
	// COMMA—comma character (,) which is the default delimiter
	// PIPE—pipe character (|)
	// SEMICOLON—semicolon character (;)
	// TAB—tab character

	// ColumnDelimiterBackquote `
	ColumnDelimiterBackquote ColumnDelimiter = `BACKQUOTE`

	// ColumnDelimiterCaret ^
	ColumnDelimiterCaret ColumnDelimiter = `CARET`

	// ColumnDelimiterComma ,
	ColumnDelimiterComma ColumnDelimiter = `COMMA`

	// ColumnDelimiterPipe |
	ColumnDelimiterPipe ColumnDelimiter = `PIPE`

	// ColumnDelimiterSemicolon :
	ColumnDelimiterSemicolon ColumnDelimiter = `SEMICOLON`

	// ColumnDelimiterTab tab
	ColumnDelimiterTab ColumnDelimiter = `TAB`
)

// LineEnding mode
type LineEnding string

const (

	// The line ending used for CSV job data, marking the end of a data row.
	// The default is LF. Valid values are:
	//
	// LF—linefeed character
	// CRLF—carriage return character followed by a linefeed character

	// LineEndingLF LF
	LineEndingLF LineEnding = `LF`

	// LineEndingCRLF CRLF
	LineEndingCRLF LineEnding = `CRLF`
)

// Operation on the individual records
type Operation string

const (

	// To ensure referential integrity, the delete operation supports cascading deletions.
	// If you delete a parent record, you delete its children automatically,
	// as long as each child record can be deleted.

	// OperationInsert inserts a new record
	OperationInsert Operation = `insert`

	// OperationDelete deletes a record
	OperationDelete Operation = `delete`

	// OperationUpdate updates a recored
	OperationUpdate Operation = `update`

	// OperationUpsert upserts a record
	OperationUpsert Operation = `upsert`

	// When the hardDelete value is specified, the deleted records aren't stored in the Recycle Bin.
	// Instead, they become immediately eligible for deletion. The permission for this operation,
	// “Bulk API Hard Delete,” is disabled by default and must be enabled by an administrator

	// OperationHardDelete (Bulk V1 type jobs only)
	OperationHardDelete Operation = `hardDelete` // (Bulk V1 type jobs only)

	// OperationQuery (Bulk V1 type jobs only)
	OperationQuery Operation = `query` // (Bulk V1 type jobs only)

	// OperationQyeryAll (Bulk V1 type jobs only)
	OperationQyeryAll Operation = `queryall` //(Bulk V1 type jobs only)

)

// ConcurrencyMode mode
type ConcurrencyMode string

const (
	// Parallel: Process records in parallel mode. This is the default value.
	// Serial: Process records in serial mode.
	//
	// Processing in parallel can cause database contention.
	// When this is severe, the job can fail. If you’re experiencing this issue,
	// submit the job with serial concurrency mode. This mode guarantees that
	// records are processed serially, but can significantly increase the processing time.

	// ConcurrencyModeParallel runs in parallel
	ConcurrencyModeParallel ConcurrencyMode = `Parallel`

	// ConcurrencyModeSerial runs in serial
	ConcurrencyModeSerial ConcurrencyMode = `Serial`
)

// JobType description
type JobType string

const (

	// JobTypeBigObjectIngest ...
	JobTypeBigObjectIngest JobType = `BigObjectIngest`

	// JobTypeClassic ...
	JobTypeClassic JobType = `Classic`

	// JobTypeV2Ingest ...
	JobTypeV2Ingest JobType = `V2Ingest`
)

// JobState used to control and monitor jobs
type JobState string

const (

	// Open—The job has been created, and data can be added to the job.
	// UploadComplete—No new data can be added to this job. You can’t edit or save a closed job.
	// Aborted—The job has been aborted. You can abort a job if you created it or if you have the “Manage Data Integrations” permission.
	// JobComplete—The job was processed by Salesforce. Failed—The job has failed.
	//
	// Job data that was successfully processed isn’t rolled back.

	// JobStateOpen job is open to add batches
	JobStateOpen JobState = `Open`

	// JobStateClosed job is open to add batches
	JobStateClosed JobState = `Closed`

	// JobStateUploadComplete job is queued or running
	JobStateUploadComplete JobState = `UploadComplete`

	// JobStateAborted job was aborted
	JobStateAborted JobState = `Aborted`

	// JobStateJobComplete job is completed
	JobStateJobComplete JobState = `JobComplete`

	// JobStateFailed job has failed
	JobStateFailed JobState = `Failed`
)

// BulkJobReq is used to create a new job
type BulkJobReq struct {
	ContentType         ContentType `force:"contentType,omitempty"`         // The content type for the job. The only valid value (and the default) is CSV.
	ExternalIdFieldName string      `force:"externalIdFieldName,omitempty"` // The external ID field in the object being updated. Only needed for upsert operations. Field values must also exist in CSV job data.
	Object              string      `force:"object"`                        // The object type for the data being processed. Use only a single object type per job.
	Operation           Operation   `force:"operation"`                     // The processing operation for the job. Valid values are: insert delete update upsert
}

// BulkJob s represented by the API.
type BulkJob struct {
	API        *API
	APIVersion float32 `json:"apiVersion,omitempty"` // The API version of the job set in the URI when the job was created.

	ID             string    `json:"id,omitempty"`             // Unique ID for this job.
	Operation      Operation `json:"operation,omitempty"`      // The processing operation for the job. Values include: insert delete update upsert hardDelete query queryAll.
	Object         string    `json:"object,omitempty"`         // The object type for the data being processed.
	CreatedById    string    `json:"createdById,omitempty"`    // The ID of the user who created the job.
	CreatedDate    string    `json:"createdDate,omitempty"`    // The date and time in the UTC time zone when the job was created.
	SystemModstamp string    `json:"systemModstamp,omitempty"` // Date and time in the UTC time zone when the job finished.

	State               JobState        `json:"state,omitempty"`               // The current state of processing for the job.
	ExternalIdFieldName string          `json:"externalIdFieldName,omitempty"` // The name of the external ID field for an upsert.
	ConcurrencyMode     ConcurrencyMode `json:"concurrencyMode,omitempty"`     // The concurrency mode for the job.
	ContentType         ContentType     `json:"contentType,omitempty"`         // The format of the data being processed.

	NumberBatchesQueued     int `json:"numberBatchesQueued,omitempty"`     // The number of batches queued for this job.
	NumberBatchesInProgress int `json:"numberBatchesInProgress,omitempty"` // The number of batches that are in progress for this job.
	NumberBatchesCompleted  int `json:"numberBatchesCompleted,omitempty"`  // The number of batches that have been completed for this job.
	NumberBatchesFailed     int `json:"numberBatchesFailed,omitempty"`     // The number of batches that have failed for this job.
	NumberBatchesTotal      int `json:"numberBatchesTotal,omitempty"`      // The number of total batches currently in the job. This value increases as more batches are added to the job.

	NumberRecordsProcessed int `json:"numberRecordsProcessed,omitempty"` // The number of records already processed. This number increases as more batches are processed.
	NumberRetries          int `json:"numberRetries,omitempty"`          // The number of times that Salesforce attempted to save the results of an operation. The repeated attempts are due to a problem, such as a lock contention.
	NumberRecordsFailed    int `json:"numberRecordsFailed,omitempty"`    // The number of records that were not processed successfully in this job.

	ApiActiveProcessingTime int64 `json:"apiActiveProcessingTime,omitempty"` // The number of milliseconds taken to actively process the job.
	ApexProcessingTime      int64 `json:"apexProcessingTime,omitempty"`      // The number of milliseconds taken to process triggers and other processes related to the job data.
	TotalProcessingTime     int64 `json:"totalProcessingTime,omitempty"`     // The number of milliseconds taken to process the job. This is the sum of the total processing times for all batches in the job.

	// The ID of a specific assignment rule to run for a case or a lead.
	// The assignment rule can be active or inactive. The ID can be retrieved
	// by using the SOAP-based SOAP API to query the AssignmentRule object.
	// AssignmentRuleId string `json:"assignmentRuleId,omitempty"` // null string
}

// CreateBulkJob creates an Bulk API 2.0 Job.
// sObjectExtID used for upsert and should be empty if not needed.
func (forceAPI *API) CreateBulkJob(sObjectname string, sObjectExtID string, contentType ContentType, operation Operation, mode ConcurrencyMode) (*BulkJob, error) {
	uri := fmt.Sprintf(`/services/async/%s/job`, strings.TrimPrefix(forceAPI.apiVersion, `v`))

	req := &Bulkv2JobReq{
		ContentType:         contentType,
		ExternalIdFieldName: sObjectExtID,
		Object:              sObjectname,
		Operation:           operation,
	}

	var resp BulkJob
	err := forceAPI.Post(uri, nil, nil, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("Bulk API: Can't create new Job for '%s': %s", sObjectname, err)
	}
	resp.API = forceAPI

	return &resp, nil
}

// AddBatch adds a new batch to a job by sending a POST request to the following URI.
// The request body contains a list of records for processing.
func (b *BulkJob) AddBatch(batch []byte) (*Batch, error) {
	uri := fmt.Sprintf(`/services/async/%.1f/job/%s/batch`, b.APIVersion, b.ID)

	var resp Batch
	err := b.API.Post(uri, nil, nil, string(batch), &resp)
	if err != nil {
		return nil, fmt.Errorf("Bulk API: Can't add Batch to Job '%s': %s", b.ID, err)
	}

	resp.API = b.API
	resp.APIVersion = b.APIVersion

	return &resp, nil
}

// Close the job, starts executing the batches.
func (b *BulkJob) Close() (err error) {
	uri := fmt.Sprintf(`/services/async/%.1f/job/%s`, b.APIVersion, b.ID)

	req := fmt.Sprintf(`{"state":"%s"}`, JobStateClosed)

	err = b.API.Post(uri, nil, nil, req, b)
	if err != nil {
		return fmt.Errorf("Bulk API: Can't close Job '%s': %s", b.ID, err)
	}
	return nil
}

// Abort the job.
func (b *BulkJob) Abort() (err error) {
	uri := fmt.Sprintf(`/services/async/%.1f/job/%s`, b.APIVersion, b.ID)

	req := fmt.Sprintf(`{"state":"%s"}`, JobStateAborted)

	err = b.API.Patch(uri, nil, nil, req, b)
	if err != nil {
		return fmt.Errorf("Bulk API: Can't abort Job '%s': %s", b.ID, err)
	}
	return nil
}

// Info of the current job.
func (b *BulkJob) Info() (err error) {
	uri := fmt.Sprintf(`/services/async/%.1f/job/%s`, b.APIVersion, b.ID)

	err = b.API.Get(uri, nil, nil, b)
	if err != nil {
		return fmt.Errorf("Bulk API: Can't retrieve Info for Job '%s': %s", b.ID, err)
	}
	return nil
}

// BulkJobBatches as represented by the API.
type BulkJobBatches struct {
	BatchInfo []*Batch `json:"batchInf,omitempty"`
}

// GetBatchesInfo gets information about all batches in a job.
func (b *BulkJob) GetBatchesInfo() (batchInfo *BulkJobBatches, err error) {
	uri := fmt.Sprintf(`/services/async/%.1f/job/%s/batch`, b.APIVersion, b.ID)

	err = b.API.Get(uri, nil, nil, batchInfo)
	if err != nil {
		return nil, fmt.Errorf("Bulk API: Can't retrieve Info for Job '%s': %s", b.ID, err)
	}
	return batchInfo, nil
}

// BatchState as represented by the API.
type BatchState string

const (

	// BatchStateQueued Processing of the batch has not started yet.
	// If the job associated with this batch is aborted, the batch isn’t
	// processed and its state is set to
	BatchStateQueued BatchState = `Queued`

	// BatchStateInProgress the batch is being processed.
	// If the job associated with the batch is aborted, the batch is still
	// processed to completion. You must close the job associated
	// with the batch so that the batch can finish processing.
	BatchStateInProgress BatchState = `InProgress`

	// BatchStateCompleted the batch has been processed completely,
	// and the result resource is available. The result resource indicates if
	// some records have failed. A batch can be completed even if some or all
	// the records have failed. If a subset of records failed, the successful
	// records aren’t rolled back.
	BatchStateCompleted BatchState = `Completed`

	// BatchStateFailed the batch failed to process the full request due to an
	// unexpected error, such as the request is compressed with an unsupported
	// format, or an internal server error. The stateMessage element could
	// contain more details about any failures. Even if the batch failed, some
	// records could have been completed successfully. The numberRecordsProcessed
	// field tells you how many records were processed. The numberRecordsFailed
	// field contains the number of records that were not processed successfully.
	BatchStateFailed BatchState = `Failed`

	// BatchStateNotProcessed The batch won’t be processed. This state is assigned
	// when a job is aborted while the batch is queued. For bulk queries, if the job
	// has PK chunking enabled, this state is assigned to the original batch that
	// contains the query when the subsequent batches are created. After the original
	// batch is changed to this state, you can monitor the subsequent batches and
	// retrieve each batch’s results when it’s completed. Then you can safely close the job.
	BatchStateNotProcessed BatchState = `NotProcessed`
)

// Batch holds the batch structure as represented by the API.
type Batch struct {
	API        *API
	APIVersion float32 `json:"apiVersion,omitempty"` // The API version of the job set in the URI when the job was created.

	ID    string     `json:"id,omitempty"`    // The ID of the batch.
	JobID string     `json:"jobId,omitempty"` // The unique, 18–character ID for the job associated with this batch.
	State BatchState `json:"state,omitempty"` // The current state of processing for the batch.

	// StateMessage Contains details about the state.
	// For example, if the state value is Failed, this field contains the reasons for failure.
	// If there are multiple failures, the message may be truncated. If so, fix the known errors
	// and re-submit the batch. Even if the batch failed, some records could have been completed successfully.
	StateMessage            string `json:"stateMessage,omitempty"`
	CreatedDate             string `json:"createdDate,omitempty"`             // The date and time in the UTC time zone when the batch was created.
	SystemModstamp          string `json:"systemModstamp,omitempty"`          // The date and time in the UTC time zone that processing ended. This is only valid when the state is Completed.
	NumberRecordsProcessed  int    `json:"numberRecordsProcessed,omitempty"`  // The number of records processed in this batch at the time the request was sent.
	NumberRecordsFailed     int    `json:"numberRecordsFailed,omitempty"`     // The number of records that were not processed successfully in this batch.
	TotalProcessingTime     int64  `json:"totalProcessingTime,omitempty"`     // The number of milliseconds taken to process the batch. This excludes the time the batch waited in the queue to be processed.
	ApiActiveProcessingTime int64  `json:"apiActiveProcessingTime,omitempty"` // The number of milliseconds taken to actively process the batch, and includes apexProcessingTime.
	ApexProcessingTime      int64  `json:"apexProcessingTime,omitempty"`      // The number of milliseconds taken to process triggers and other processes related to the batch data.
}

// Info gets information about an existing batch.
func (b *Batch) Info() (err error) {
	uri := fmt.Sprintf(`/services/async/%.1f/job/%s/batch/%s`, b.APIVersion, b.JobID, b.ID)

	err = b.API.Get(uri, nil, nil, b)
	if err != nil {
		return fmt.Errorf("Bulk API: Can't retrieve Info for Batch '%s': %s", b.ID, err)
	}
	return nil
}

// BatchRequest is a json represantion of the request.
type BatchRequest map[string]interface{}

// Request gets the request of a batch.
func (b *Batch) Request() ([]BatchRequest, error) {
	uri := fmt.Sprintf(`/services/async/%.1f/job/%s/batch/%s/request`, b.APIVersion, b.JobID, b.ID)

	batchRequest := []BatchRequest{}

	err := b.API.Get(uri, nil, nil, &batchRequest)
	if err != nil {
		return batchRequest, fmt.Errorf("Bulk API: Can't retrieve Request for Batch '%s': %s", b.ID, err)
	}

	return batchRequest, nil
}

// BatchError as represented within the API.
type BatchError struct {
	Fields     []string `json:"fields"`
	Message    string   `json:"message"`
	StatusCode string   `json:"statusCode"`
}

// BatchResult as represented within the API.
type BatchResult struct {
	Success bool          `json:"success"`
	Created bool          `json:"created"`
	Id      string        `json:"id"`
	Errors  []*BatchError `json:"errors"`
}

// Result gets results of a batch that has completed processing.
func (b *Batch) Result() ([]BatchResult, error) {
	uri := fmt.Sprintf(`/services/async/%.1f/job/%s/batch/%s/result`, b.APIVersion, b.JobID, b.ID)

	batchResult := []BatchResult{}

	err := b.API.Get(uri, nil, nil, &batchResult)
	if err != nil {
		return batchResult, fmt.Errorf("Bulk API: Can't retrieve Result for Batch '%s': %s", b.ID, err)
	}
	return batchResult, nil
}
