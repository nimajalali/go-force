package force

import (
	"testing"
)

func TestLimits(t *testing.T) {
	forceApi := createTest()
	limits, err := forceApi.GetLimits()
	if err != nil {
		// Developer Accounts, which the testbed uses, do not have access to the limits API. So this will always fail.
		// t.Fatalf("Failed to get Limits: %v", err)
		t.Logf("Failed to get Limits, this is expected due to the developer account: %v", err)
	}

	t.Log(limits)
}
