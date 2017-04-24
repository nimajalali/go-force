package force

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// FileUpload creates a new file upload http request with optional extra params
func (j *SJob) FileUpload(params map[string]string, paramName string, file *os.File) (*SBatch, error) {

	fi, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Cannot get request file info: %s", err)
	}

	URI := fmt.Sprintf("%s/%s/batch", j.forceAPI.oauth.InstanceURL, j.BaseURI)
	req, err := http.NewRequest("POST", URI, file)

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", responseType)
	req.Header.Set("Accept-Encoding", gzipEncodingType)
	req.Header.Set("X-SFDC-Session", j.forceAPI.oauth.AccessToken)
	req.Header.Set("Content-Type", zipJSONContentType)
	req.Header.Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

	fmt.Printf("--------------------FILE REQ:::: %+v\n", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Cannot send file upload request: %s", err)
	}

	batch := &SBatch{}
	err = j.forceAPI.readResponse(resp, "POST", "", nil, nil, "", batch)

	return batch, nil
}
