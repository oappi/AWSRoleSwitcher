package onePasswordLogic

/*
We are using 1password cliV2 to query entities stored in 1password. These include credentials, accountlist and more
*/

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"gopkg.in/ini.v1"
)

/*
signinCommand gives signin command for 1password client.
This is always needed when we query something from 1password before we issue the actual command we want to run
*/
func SigninCommand(password string, domain string) string {
	var signin = "eval $(echo " + password + " | op signin --account " + domain + ") && "
	return signin
}

func OpCMDlogic(command string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		cmdError := errors.New(stderr.String())
		return "", cmdError
	} else {
		return strings.TrimSpace(out.String()), nil
	}
}

/*
Here we fetch MFA by running signincommand and then mfa fetch command
*/
func GetMFA(signInCommand string, entityName string) (string, error) {
	var mfaFetchCommand = "op item get " + entityName + " --otp"
	output, err := OpCMDlogic(signInCommand + mfaFetchCommand)
	return output, err
}

/*
Rest work similarly to GetMFA function
*/

func GetAccesskey(signInCommand string, entityName string) (string, error) {
	var akFetchCommand = "op item get " + entityName + " --fields access_key"
	return OpCMDlogic(signInCommand + akFetchCommand)
}

func GetSecretAccesskey(signInCommand string, entityName string) (string, error) {
	var sakFetchCommand = "op item get " + entityName + " --fields secret_access_key"
	return OpCMDlogic(signInCommand + sakFetchCommand)
}
func SaveAccesskeys(signInCommand, entityName, accesskey, secretaccesskey string) error {
	var accessKeyUpdateCommand = "op item edit " + entityName + " access_key=" + accesskey
	_, errors := OpCMDlogic(signInCommand + accessKeyUpdateCommand)
	if errors != nil {
		return errors
	}
	var secretAccessKeyUpdateCommand = "op item edit " + entityName + " secret_access_key=" + secretaccesskey
	_, error := OpCMDlogic(signInCommand + secretAccessKeyUpdateCommand)
	if errors != nil {
		return error
	}
	return nil
}

func GetAlias(signInCommand string, entityName string) (string, error) {
	var settingsFetchCommand = "op item get " + entityName + " --fields Alias"
	return OpCMDlogic(signInCommand + settingsFetchCommand)
}

func GetFetchAccountCommand(entityName string) string {
	return "op item get " + entityName + " --fields Accounts"
}

func GetAccounts(signInCommand string, fetchAccountCommand string) ([]string, error) {
	accountStringListRaw, err := OpCMDlogic(signInCommand + fetchAccountCommand)
	if err != nil {
		return nil, err
	}
	accountsParsed, _ := OPAccountListParser(accountStringListRaw)
	accountList, error := convert1PasswordAccountStringToList(accountsParsed)

	return accountList, error
}

func GetMFADevice(signInCommand string, entityName string) (string, error) {
	var accountsFetchCommand = "op item get " + entityName + " --fields MFADevice"
	return OpCMDlogic(signInCommand + accountsFetchCommand)
}

func GetRegion(signInCommand string, entityName string) (string, error) {
	var accountsFetchCommand = "op item get " + entityName + " --fields Region"
	return OpCMDlogic(signInCommand + accountsFetchCommand)
}

func convert1PasswordAccountStringToList(accounts string) ([]string, error) {
	accountList := []string{}
	accountByteArray := []byte(accounts)
	accountsList, _ := ini.Load(accountByteArray) //error on not able to read string
	for _, s := range accountsList.Sections() {
		var credentialElementName = s.Name()
		accountId := s.Key("aws_account_id").String()
		roleName := s.Key("role_name").String()
		if accountId != "" && roleName != "" {
			accountList = append(accountList, credentialElementName+"|"+accountId+"|"+roleName)
		}
	}
	return accountList, nil
}

func OPAccountListParser(accountStringListRaw string) (string, error) {
	var items = len(accountStringListRaw)
	var accountsParsed = accountStringListRaw[:items-1]
	accountsParsed = accountsParsed[1:]

	return accountsParsed, nil
}
