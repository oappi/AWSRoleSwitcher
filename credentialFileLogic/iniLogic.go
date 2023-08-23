package credentialFileLogic

import (
	"os"

	"gopkg.in/ini.v1"
)

func GetLocalSettings(AWSFolderlocation string) (string, string, string, string, string, string, error) {
	cfg, err := ini.Load(AWSFolderlocation + "awsroleswitcher")
	if err != nil {
		return "", "", "", "", "", "", err
	}
	mfaDevice := cfg.Section("localSettings").Key("MFADevice").String()
	mfaSeed := cfg.Section("localSettings").Key("MFASeed").String()
	accesskey := cfg.Section("localSettings").Key("access_key").String()
	secretaccesskey := cfg.Section("localSettings").Key("secret_access_key").String()
	alias := cfg.Section("localSettings").Key("alias").String()
	region := cfg.Section("localSettings").Key("Region").String()

	return mfaDevice, mfaSeed, accesskey, secretaccesskey, alias, region, nil
}

func saveWithReducedPriviliges(fullFilePath string, cfg *ini.File) error {
	err := cfg.SaveTo(fullFilePath)
	if err != nil {
		return err
	}
	cerr := os.Chmod(fullFilePath, 0600)
	if cerr != nil {
		return cerr
	}
	return nil
}

func SetLocalSettings(AWSFolderlocation, MFADevice, MFASeed, access_key, secret_access_key, alias string) error {
	cfg, err := ini.Load(AWSFolderlocation + "awsroleswitcher")
	if err != nil {
		cfg = ini.Empty()
	}
	cfg.Section("localSettings").Key("MFADevice").SetValue(MFADevice)
	cfg.Section("localSettings").Key("MFASeed").SetValue(MFASeed)
	cfg.Section("localSettings").Key("access_key").SetValue(access_key)
	cfg.Section("localSettings").Key("secret_access_key").SetValue(secret_access_key)
	cfg.Section("localSettings").Key("alias").SetValue(alias)
	return saveWithReducedPriviliges(AWSFolderlocation+"awsroleswitcher", cfg)
}

func Get1PasswordSettings(AWSFolderlocation string) (string, string) {
	cfg, err := ini.Load(AWSFolderlocation + "awsroleswitcher")
	if err != nil {
		return "", ""
	}
	domain := cfg.Section("1password").Key("domain").String()
	entity := cfg.Section("1password").Key("entity").String()
	return domain, entity
}

func Set1PasswordSettings(AWSFolderlocation, domain, entity string) error {
	cfg, err := ini.Load(AWSFolderlocation + "awsroleswitcher")
	if err != nil {
		cfg = ini.Empty()
	}
	cfg.Section("1password").Key("domain").SetValue(domain)
	cfg.Section("1password").Key("entity").SetValue(entity)
	return saveWithReducedPriviliges(AWSFolderlocation+"awsroleswitcher", cfg)
}

func UpdateShortTermAWSKeys(accesskey, secretaccesskey, token, AWSFolderlocation string) error {
	cfg, err := ini.Load(AWSFolderlocation + "credentials")
	if err != nil {
		return err
	}
	cfg.Section("default").Key("aws_access_key_id").SetValue(accesskey)
	cfg.Section("default").Key("aws_secret_access_key").SetValue(secretaccesskey)
	cfg.Section("default").Key("aws_session_token").SetValue(token)
	println("Note, Saving to AWS credentials file")
	return saveWithReducedPriviliges(AWSFolderlocation+"credentials", cfg)
}

func UpdateLongTermAWSKeys(accesskey, secretaccesskey, token, AWSFolderlocation string) error {
	settingsFile := AWSFolderlocation + "awsroleswitcher"
	cfg, err := ini.Load(settingsFile)
	if err != nil {
		return err
	}
	cfg.Section("localSettings").Key("aws_access_key_id").SetValue(accesskey)
	cfg.Section("localSettings").Key("aws_secret_access_key").SetValue(secretaccesskey)
	if token != "" {
		cfg.Section("localSettings").Key("aws_session_token").SetValue(secretaccesskey)
	}
	println("Note, Saving to AWS credentials file")
	return saveWithReducedPriviliges(settingsFile, cfg)
}

/*
gets long term AWS keys from local settings file
*/
func GetAWSKeys(accesskey, secretaccesskey, token, AWSFolderlocation string) (string, string, string, error) {
	cfg, err := ini.Load(AWSFolderlocation)
	if err != nil {
		return "", "", "", err
	}
	return cfg.Section("default-long-term").Key("aws_access_key_id").Value(), cfg.Section("default-long-term").Key("aws_secret_access_key").Value(), cfg.Section("default-long-term").Key("aws_session_token").Value(), nil
}

/*
saves long term AWS keys from local settings file
*/
func SetAWSKeys(accesskey string, secretaccesskey string, AWSFolderlocation string) error {
	cfg, _ := ini.Load(AWSFolderlocation) //
	cfg.Section("default-long-term").Key("aws_access_key_id").SetValue(accesskey)
	cfg.Section("default-long-term").Key("aws_secret_access_key").SetValue(secretaccesskey)
	return nil

}

func GetAccountList(AWSFolderlocation string) ([]string, error) {
	accountList := []string{}
	accountsList, err := ini.Load(AWSFolderlocation) //error on not able to read string
	if err != nil {
		return accountList, err
	}
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
