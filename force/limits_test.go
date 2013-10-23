package force

import (
	"testing"
)

func init() {
	initTest()
}

func TestLimits(t *testing.T) {
	limits, err := GetLimits()
	if err != nil {
		// Developer Accounts, which the testbed uses, do not have access to the limits API. So this will always fail.
		// t.Fatalf("Failed to get Limits: %v", err)
		t.Logf("Failed to get Limits, this is expected due to the developer account: %v", err)
	}

	t.Log(limits)
}
