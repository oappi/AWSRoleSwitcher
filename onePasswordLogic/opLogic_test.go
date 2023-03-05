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
	var simulatedAnswer1P string = ` "\"[simulatedAccount] 
	aws_account_id = 111112222333
	role_name = exampleRoleName

	[simulatedAccount2] 
	aws_account_id = 111112222334
	role_name = exampleRoleName1\""
	`
	accounts, err := GetAccounts("echo", simulatedAnswer1P)

	if len(accounts) == 0 {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
