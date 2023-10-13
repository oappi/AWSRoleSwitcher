package interfaces

import (
	creds "github.com/oappi/awsroleswitcher/credentialFileLogic"
)

type LocalWriter interface {
	UpdateShortTermKeys(accesskey, secretAccessKey, token string) error
	UpdateLongTermKeys(accesskey, secretAccessKey, token string) error
}

type IniLogic struct {
	AWSFolderLocation string
}

func (ini IniLogic) UpdateShortTermKeys(accesskey, secretAccessKey, token string) error {
	creds.UpdateShortTermAWSKeys(accesskey, secretAccessKey, token, ini.AWSFolderLocation)
	return nil
}

func (ini IniLogic) UpdateLongTermKeys(accesskey, secretAccessKey, token string) error {
	return creds.UpdateLongTermAWSKeys(accesskey, secretAccessKey, token, ini.AWSFolderLocation)

}

func (ini IniLogic) Get1PasswordSettings() (string, string) {
	return creds.Get1PasswordSettings(ini.AWSFolderLocation)
}

func (ini IniLogic) Set1PasswordSettings(domain, entity string) {
	creds.Set1PasswordSettings(ini.AWSFolderLocation, domain, entity)
}

func (ini IniLogic) GetLocalSettings() (string, string, string, string, string, string, error) {
	return creds.GetLocalSettings(ini.AWSFolderLocation)
}

func (ini IniLogic) SetLocalSettings(MFADevice, MFASeed, access_key, secret_access_key, alias, region string) {
	creds.SetLocalSettings(ini.AWSFolderLocation, MFADevice, MFASeed, access_key, secret_access_key, alias, region)
}
