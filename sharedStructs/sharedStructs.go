package sharedStructs

import (
	//	"sync"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type STSConfig struct {
	STS stsiface.STSAPI
}

type SessionInfo struct {
	Accesskey       string
	SecretAccessKey string
	Token           string
}

/*
	type MFAInfo struct {
		profilename   string
		opdomain      string
		opuuid        string
		password      string
		accountnumber string
		switchrole    string
		lock          *sync.Mutex
		awsSession    SessionInfo
		region        string
	}
*/
func CredentialFileSplitter(r rune) bool {
	return r == '[' || r == ']'
}

func ConfigsSplitter(r rune) bool {
	return r == '\n' || r == '\r'
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
