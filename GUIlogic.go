// this file includes call logic. Basically supports main.go
// Purpose is to have GUI definitions in main.go file, and actual logic in this file.
// Basically anything that doesnt draw something to GUI should be here
package main

import (
	"strconv"
	"strings"

	"github.com/oappi/awsroler/awsLogic"
	"github.com/oappi/awsroler/interfaces"
	"github.com/oappi/awsroler/sharedStructs"
)

func UpdateList() {
	var empty = ""
	filteredList, _ := filteredListForSelect(&empty)
	gOptionSelection.SetOptions(filteredList)
}

func UpdateUISettings(settings sharedStructs.FederationAccountSettingsObject) {
	accountsList = settings.Accounts
	SettingsObject.MFA = settings.MFA
	SettingsObject.Alias = settings.Alias
	SettingsObject.AccessKey = settings.AccessKey
	SettingsObject.SecretAccessKey = settings.SecretAccessKey
	SettingsObject.Region = settings.Region
	SettingsObject.Accounts = settings.Accounts
	SettingsObject.MFADevice = settings.MFADevice
}

type FederationAccountSettingsObject struct {
	MFA             string
	Alias           string
	AccessKey       string
	SecretAccessKey string
	Region          string
	Accounts        []string
	MFADevice       string
}

func filteredListForSelect(filter *string) (resultList []string, match bool) {
	var exactMatch = false
	if filter == nil || *filter == "" {
		return accountsList, false
	} else {
		for _, option := range accountsList {
			var filterstring = *filter
			if strings.Contains(strings.ToLower(option), strings.ToLower(filterstring)) {
				resultList = append(resultList, option)
				if option == filterstring {
					exactMatch = true
				}
			}
		}
		return resultList, exactMatch
	}
}

func filteredCustomListForSelect(filter *string, list []string) (resultList []string, match bool) {
	var exactMatch = false
	for _, option := range list {
		var filterstring = *filter
		if strings.Contains(strings.ToLower(option), strings.ToLower(filterstring)) {
			resultList = append(resultList, option)
			if option == filterstring {
				exactMatch = true
			}
		}
	}
	return resultList, exactMatch
}

func OverRideSavedIfUserGivesInput(userInput string, savedInput string) string {
	if userInput == "" {
		return userInput
	} else {
		return savedInput
	}
}

func fetchFederationCredentialsAndSettings(settings interfaces.SettingsInterface) (sharedStructs.FederationAccountSettingsObject, error) {
	var emptyArray []string
	var emptySettings = sharedStructs.FederationAccountSettingsObject{MFA: "", Alias: "", AccessKey: "", SecretAccessKey: "", Region: "", Accounts: emptyArray, MFADevice: ""}
	mfa, getMfaError := settings.GetMFA()
	if getMfaError != nil {
		return emptySettings, getMfaError
	}
	accessKey := settings.GetAccesskey()
	secretAccessKey := settings.GetSecretAccessKey()
	alias := settings.GetAlias()
	mfaDevice := settings.GetMFADevice()
	region, regionError := settings.GetRegion()
	if regionError != nil {
		return emptySettings, regionError
	}
	accounts, accountError := settings.GetAccounts()
	if accountError != nil {
		return emptySettings, accountError
	}

	settingsObject := sharedStructs.FederationAccountSettingsObject{MFA: mfa, Alias: alias, AccessKey: accessKey, SecretAccessKey: secretAccessKey, Region: region, Accounts: accounts, MFADevice: mfaDevice}
	return settingsObject, nil
}

func updateSettings(SettingsInterface interfaces.SettingsInterface) error {
	federationSettings, err := fetchFederationCredentialsAndSettings(SettingsInterface)
	if err != nil {
		return err
	}
	UpdateUISettings(federationSettings)
	UpdateList()
	return nil
}

func FetchAndSaveAccountCredentials(stsConfig *sharedStructs.STSConfig, accountToConnect string, accountRole string, region string, sessionTime int64) error {
	mfa, mfaError := SettingsInterface.GetMFA()
	if mfaError != nil {
		return mfaError
	}
	SettingsObject.MFA = mfa
	accesskey, saccesskey, token, assumeRoleError := awsLogic.GetAsumeRoleCredentials(stsConfig, SettingsObject, accountToConnect, accountRole, sessionTime)
	if assumeRoleError != nil {
		return assumeRoleError
	}
	awsSession.Accesskey = accesskey
	awsSession.SecretAccessKey = saccesskey
	awsSession.Token = token
	return nil
}

func getStSConfig(SettingsObject sharedStructs.FederationAccountSettingsObject) (*sharedStructs.STSConfig, error) {
	stsConfig, stsError := awsLogic.InitializeSTSConfig(SettingsObject)
	if stsError != nil {
		return nil, stsError
	} else {
		return stsConfig, nil
	}
}

func connectAccount(STSConfig *sharedStructs.STSConfig, selectedAccountInfo string, writer interfaces.LocalWriter, sessionTimeOption string) error {
	sessionTime, sessionTimeError := ParseSessiontime(sessionTimeOption)
	if sessionTimeError != nil {
		return sessionTimeError
	}
	var splittedaccountinfo = strings.Split(selectedAccountInfo, "|")
	var accountToConnect = splittedaccountinfo[1]
	var accountRole = splittedaccountinfo[2]
	FetchAndSaveAccountCredentials(STSConfig, accountToConnect, accountRole, region, sessionTime)
	writer.UpdateShortTermKeys(awsSession.Accesskey, awsSession.SecretAccessKey, awsSession.Token)
	return nil
}

func ParseSessiontime(sessionTimeOption string) (int64, error) {
	sessionHoursString := strings.Split(sessionTimeOption, " ")[0]
	return strconv.ParseInt(sessionHoursString, 10, 64) //converts string to int, error if failed
}
