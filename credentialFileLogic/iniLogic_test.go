package credentialFileLogic

import (
	"os"
	"testing"
)

func TestUpdateAWSKeys(t *testing.T) {
	var testcredentialfile = UpdateShortTermAWSKeys("1234567", "098765", "5453653324", os.Getenv("HOME")+"/.aws/credentials")
	if testcredentialfile != nil {
		t.Fail()
	}
}
