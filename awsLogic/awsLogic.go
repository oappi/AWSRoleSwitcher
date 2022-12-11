package awsLogic

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	iam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/oappi/awsroler/interfaces"
	"github.com/oappi/awsroler/sharedStructs"
)

func InitializeSTSConfig(settingsFile sharedStructs.FederationAccountSettingsObject) (*sharedStructs.STSConfig, error) {
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

func CreateNewAccesskey(settingsInterface interfaces.SettingsInterface, region string) (string, string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(settingsInterface.GetAccesskey(), settingsInterface.GetSecretAccessKey(), ""),
	})
	svc := iam.New(sess)
	numberOfKeys, err := GetNumberOfAccessKeys(svc, region)
	if err != nil {
		return "", "", err

	} else if numberOfKeys == 1 { //we want to support only creating accesskeys when we try to rotate, since max is 2 we want to have only 1 avaible
		iamInput := &iam.CreateAccessKeyInput{}
		result, err := svc.CreateAccessKey(iamInput)
		return *result.AccessKey.AccessKeyId, *result.AccessKey.SecretAccessKey, err
	} else {
		return "", "", errors.New("we only support rotation to when there are only 1 accesskey avaible")
	}

}

func GetNumberOfAccessKeys(svc *iam.IAM, region string) (int, error) {

	input := &iam.ListAccessKeysInput{}

	result, err := svc.ListAccessKeys(input)

	var lister = result.AccessKeyMetadata

	return len(lister), err
}

func DeleteAccesskeyPair(settingsInterface interfaces.SettingsInterface, region, accesskey string) error {
	input := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accesskey),
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(settingsInterface.GetAccesskey(), settingsInterface.GetSecretAccessKey(), ""),
	})
	if err != nil {
		return err
	}
	svc := iam.New(sess)
	_, error := svc.DeleteAccessKey(input)
	return error
}
