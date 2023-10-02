package interfaces

import (
	"sync"
	"time"

	cred "github.com/oappi/awsroleswitcher/credentialFileLogic"
	opLogic "github.com/oappi/awsroleswitcher/onePasswordLogic"
)

type SettingsInterface interface {
	GetMFA() (string, error)
	GetAlias() string
	GetAccesskey() string
	GetSecretAccessKey() string
	GetRegion() (string, error)
	GetAccounts() ([]string, error)
	GetMFADevice() (string, error)
	SetLongtermAccessKeys(string, string) error
	AdvancedFeaturesEnabled() bool
}

type Onepassword struct {
	Lock                     *sync.Mutex //used as writelock for local aws credentials file
	Uuid, OPDomain, Password string
	Profilename              string
}

type LocalSettings struct {
	Lock                *sync.Mutex
	MFASeed             string
	MFADevice           string
	Region              string
	AccessKey           string
	SecretAccessKey     string
	UserAlias           string
	AccountListLocation string
	AWSFolderLocation   string
	LocalWriter         LocalWriter
}

func (op Onepassword) GetMFA() (string, error) {
	signincmd := opLogic.SigninCommand(op.Password, op.OPDomain)
	return opLogic.GetMFA(signincmd, op.Uuid)
}

func (local LocalSettings) GetMFA() (string, error) {
	now := time.Now()
	interval := int64(now.Unix() / 30)
	return cred.GetOTP(local.MFASeed, interval)
}

func (op Onepassword) GetAlias() string {
	signincmd := opLogic.SigninCommand(op.Password, op.OPDomain)
	output, _ := opLogic.GetAlias(signincmd, op.Uuid)
	return output
}

func (local LocalSettings) GetAlias() string {

	return local.UserAlias
}

func (op Onepassword) GetAccesskey() string {
	signincmd := opLogic.SigninCommand(op.Password, op.OPDomain)
	output, _ := opLogic.GetAccesskey(signincmd, op.Uuid)
	return output
}

func (local LocalSettings) GetAccesskey() string {
	return local.AccessKey
}

func (op Onepassword) GetSecretAccessKey() string {
	signincmd := opLogic.SigninCommand(op.Password, op.OPDomain)
	output, _ := opLogic.GetSecretAccesskey(signincmd, op.Uuid)
	return output
}

func (local LocalSettings) GetSecretAccessKey() string {
	return local.SecretAccessKey
}

func (op Onepassword) SetLongtermAccessKeys(accesskey, secretaccesskey string) error {
	signincmd := opLogic.SigninCommand(op.Password, op.OPDomain)
	error := opLogic.SaveAccesskeys(signincmd, op.Uuid, accesskey, secretaccesskey)
	return error
}

func (local LocalSettings) SetLongtermAccessKeys(accesskey, secretaccesskey string) error {

	return LocalWriter.UpdateLongTermKeys(local.LocalWriter, accesskey, secretaccesskey, local.AWSFolderLocation)
}

func (op Onepassword) GetRegion() (string, error) {
	signincmd := opLogic.SigninCommand(op.Password, op.OPDomain)
	output, error := opLogic.GetRegion(signincmd, op.Uuid)
	return output, error
}

func (local LocalSettings) GetRegion() (string, error) {

	return local.Region, nil
}

func (op Onepassword) GetAccounts() ([]string, error) {
	signincmd := opLogic.SigninCommand(op.Password, op.OPDomain)
	fetchCommand := opLogic.GetFetchAccountCommand(op.Uuid)
	accountsRaw, err := opLogic.GetAccounts(signincmd, fetchCommand)
	return accountsRaw, err
}

func (local LocalSettings) GetAccounts() ([]string, error) {

	return cred.GetAccountList(local.AccountListLocation)

}

func (op Onepassword) GetMFADevice() (string, error) {
	signincmd := opLogic.SigninCommand(op.Password, op.OPDomain)
	output, error := opLogic.GetMFADevice(signincmd, op.Uuid)
	return output, error
}

func (local LocalSettings) GetMFADevice() (string, error) {
	return local.MFADevice, nil
}

/*
*
Do we enable things like password rotation. With 1password we check that new key works before removing old one
and if rotation fails for any reason 1password stores old values
*/
func (op Onepassword) AdvancedFeaturesEnabled() bool {
	return true
}

/*
*
Do we enable things like password rotation. Since it is possible, although unlikely,
to get in state where write corrupts localfile I prefer to disable this feature. As for 1password there
is possibility to check historic versions where user can still use old key.
*/
func (local LocalSettings) AdvancedFeaturesEnabled() bool {
	return false
}
