package awsLogic

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	iam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/oappi/awsroleswitcher/interfaces"
	"github.com/oappi/awsroleswitcher/sharedStructs"
)

func CreateSTSSession(settingsFile sharedStructs.FederationAccountSettingsObject) (*sharedStructs.STSConfig, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(settingsFile.Region),
		Credentials: credentials.NewStaticCredentials(settingsFile.AccessKey, settingsFile.SecretAccessKey, ""),
	})
	if err != nil {
		return &sharedStructs.STSConfig{}, fmt.Errorf("unable to create a session to aws with error: %v", err)
	}
	return &sharedStructs.STSConfig{
		STS: sts.New(sess),
	}, nil
}

/*
CreateIAMSession creates iam client with MFA enabled
*/
func CreateIAMSession(settingsFile sharedStructs.FederationAccountSettingsObject, SettingsInterface interfaces.SettingsInterface) (*iam.IAM, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(settingsFile.Region),
		Credentials: credentials.NewStaticCredentials(settingsFile.AccessKey, settingsFile.SecretAccessKey, ""),
	})
	_sts := sts.New(sess)

	MFASerial, deviceError := SettingsInterface.GetMFADevice()
	if deviceError != nil {
		return nil, deviceError
	}

	mfa, mfaError := SettingsInterface.GetMFA() //required for password rotation on accounts that use aws recomended policy
	if mfaError != nil {
		return nil, mfaError
	}

	resSessionToken, stsErr := _sts.GetSessionToken(&sts.GetSessionTokenInput{
		TokenCode:    &mfa,
		SerialNumber: &MFASerial,
	})
	if stsErr != nil {
		return nil, stsErr
	}

	MFASession, errMfa := session.NewSession(&aws.Config{
		Region:      aws.String(settingsFile.Region),
		Credentials: credentials.NewStaticCredentials(*resSessionToken.Credentials.AccessKeyId, *resSessionToken.Credentials.SecretAccessKey, *resSessionToken.Credentials.SessionToken),
	})

	if errMfa != nil {
		return nil, errMfa
	}

	svc := iam.New(MFASession)
	if err != nil {
		return nil, err
	}
	return svc, err
}

func GetAsumeRoleCredentials(stsConfig *sharedStructs.STSConfig, settingsFile sharedStructs.FederationAccountSettingsObject, accountnumber string, switchrole string, sessionTime int64) (string, string, string, error) {
	roleToAssumeArn := "arn:aws:iam::" + accountnumber + ":role/" + switchrole
	sessionName := settingsFile.Alias
	var duration int64 = 3600 * sessionTime
	result, err := stsConfig.STS.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &roleToAssumeArn,
		TokenCode:       &settingsFile.MFA,
		RoleSessionName: &sessionName,
		DurationSeconds: &duration,
		SerialNumber:    &settingsFile.MFADevice,
	})
	if err != nil {
		return "", "", "", err
	}
	return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, nil
}

func CreateNewAccesskey(iamSession *iam.IAM) (string, string, error) {
	iamInput := &iam.CreateAccessKeyInput{}
	result, err := iamSession.CreateAccessKey(iamInput)
	return *result.AccessKey.AccessKeyId, *result.AccessKey.SecretAccessKey, err
}

/*
GetNumberOfAccessKeys returns number of keys found
*/
func GetNumberOfAccessKeys(iamSession *iam.IAM) (int, error) {
	input := &iam.ListAccessKeysInput{}
	result, err := iamSession.ListAccessKeys(input)
	return len(result.AccessKeyMetadata), err
}

func DeleteAccesskeyPair(iamSession *iam.IAM, accesskey string) error {
	input := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accesskey),
	}
	_, error := iamSession.DeleteAccessKey(input)
	return error
}
