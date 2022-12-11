package interfaces

import (
	creds "github.com/oappi/awsroler/credentialFileLogic"
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
	creds.UpdateLongTermAWSKeys(accesskey, secretAccessKey, token, ini.AWSFolderLocation)
	return nil
}

func (ini IniLogic) Get1PasswordSettings() (string, string) {
	return creds.Get1PasswordSettings(ini.AWSFolderLocation)
}

func (ini IniLogic) Set1PasswordSettings(domain, entity string) {
	creds.Set1PasswordSettings(ini.AWSFolderLocation, domain, entity)
}
