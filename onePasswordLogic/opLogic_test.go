package onePasswordLogic

import (
	"testing"
)

func getCreds() (string, string) {
	return "nope", "my"
}

func getEntityName() string {
	return "AWSRoleSwitchTest"
}

func TestGetMFA(t *testing.T) {

	scommand := SigninCommand(getCreds())
	otp, err := GetMFA(scommand, "AWSRoleSwitchTest")
	if otp != "" {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}

/*
with correct X this should work, but commented out for now
*/
/*
func TestGetMFASucceeds(t *testing.T) {

	scommand := SigninCommand("X", "my")
	otp, err := GetMFA(scommand, "AWSRoleSwitchTest")
	if otp != "" {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}
*/

func TestOpCMDlogic(t *testing.T) {

	output, err := OpCMDlogic("echo test")
	if output != "test" {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
func TestGetAccounts(t *testing.T) {
	scommand := SigninCommand(getCreds())
	accounts, err := GetAccounts(scommand, getEntityName())

	if len(accounts) == 0 {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
