package force

import (
	"testing"
)

func init() {
	initTest()
}

func TestOauth(t *testing.T) {
	// Verify oauth object is valid
	if err := oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}
}
