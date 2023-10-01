package accesskeyRotation

import (
	"errors"

	iam "github.com/aws/aws-sdk-go/service/iam"
	awslogic "github.com/oappi/awsroleswitcher/awsLogic"
	"github.com/oappi/awsroleswitcher/interfaces"
	"github.com/oappi/awsroleswitcher/sharedStructs"
)

func RotateAccesskeys(settingsInterface interfaces.SettingsInterface, SettingsObject sharedStructs.FederationAccountSettingsObject) (error, string, string) {
	if settingsInterface.AdvancedFeaturesEnabled() == false {
		return errors.New("option is not enabled to your password service provider (such as locally stored passwords)"), "", ""
	}
	iamSession, err := awslogic.CreateIAMSession(SettingsObject)
	if err != nil {
		return err, "", ""
	}

	oneAccesskeyFound, keyerr := hasExactlyOneAccesskey(iamSession)
	if keyerr != nil {
		return keyerr, "", ""
	}
	if !oneAccesskeyFound {
		return errors.New("check that you have only one active access key"), "", ""
	}
	oldAccesskey := settingsInterface.GetAccesskey()
	oldSecretAccesskey := settingsInterface.GetSecretAccessKey()
	newAccesskey, newSecretAccesskey, err := awslogic.CreateNewAccesskey(iamSession)
	if err != nil {
		return err, "", ""
	}
	sveErrors := settingsInterface.SetLongtermAccessKeys(newAccesskey, newSecretAccesskey)
	if sveErrors != nil {
		return sveErrors, "", ""
	}
	currentAccesskey := settingsInterface.GetAccesskey()
	currentSecretAccesskey := settingsInterface.GetSecretAccessKey()
	if newAccesskey == currentAccesskey && newSecretAccesskey == currentSecretAccesskey {
		errD := awslogic.DeleteAccesskeyPair(iamSession, oldAccesskey)
		if errD != nil {
			return errD, "", ""
		}
		return nil, currentAccesskey, newSecretAccesskey
	} else {
		return errors.New("Accesskey missmatch, please check accesskey from your accesskey storage old keys:" + oldAccesskey + " || " + oldSecretAccesskey), "", ""
	}
}

func hasExactlyOneAccesskey(iamSession *iam.IAM) (bool, error) {
	numberOfKeys, err := awslogic.GetNumberOfAccessKeys(iamSession)
	if err != nil {
		return false, err
	}
	if numberOfKeys == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
