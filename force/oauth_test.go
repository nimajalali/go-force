package force

import (
	"testing"
)

func TestOauth(t *testing.T) {
	forceAPI := createTest()
	// Verify oauth object is valid
	if err := forceAPI.oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}
}
