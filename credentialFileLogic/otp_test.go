package credentialFileLogic

import (
	"testing"
)

func TestGetOTP(t *testing.T) {

	interval := int64(1658177010 / 30)
	otp, _ := GetOTP("testcasetestcase", interval)
	if otp != "786411" {
		t.Fail()
	}
}

func TestGetOTPBase32Fail(t *testing.T) {

	interval := int64(1658177010 / 30)
	_, err := GetOTP("testcasetestcas]", interval)
	if err == nil {
		t.Fail()
	}
}

func TestAddPadding(t *testing.T) {

	otp := AddPadding("1234")
	if otp != "001234" {
		t.Fail()
	}
}
