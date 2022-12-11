package accesskeyRotation

import (
	"errors"

	awslogic "github.com/oappi/awsroler/awsLogic"
	"github.com/oappi/awsroler/interfaces"
)

func RotateAccesskeys(settingsInterface interfaces.SettingsInterface) error {
	oldAccesskey := settingsInterface.GetAccesskey()
	oldSecretAccesskey := settingsInterface.GetSecretAccessKey()
	region, regionError := settingsInterface.GetRegion()
	if regionError != nil {
		return regionError
	}

	newAccesskey, newSecretAccesskey, err := awslogic.CreateNewAccesskey(settingsInterface, region)
	if err != nil {
		return err
	}
	sveErrors := settingsInterface.SetLongtermAccessKeys(newAccesskey, newSecretAccesskey)
	if sveErrors != nil {
		return sveErrors
	}
	currentAccesskey := settingsInterface.GetAccesskey()
	currentSecretAccesskey := settingsInterface.GetSecretAccessKey()
	if newAccesskey == currentAccesskey && newSecretAccesskey == currentSecretAccesskey {
		errD := awslogic.DeleteAccesskeyPair(settingsInterface, region, oldAccesskey)
		if errD != nil {
			return errD
		}
		return nil
	} else {
		return errors.New("Accesskey missmatch, please check accesskey from your accesskey storage old keys:" + oldAccesskey + " || " + oldSecretAccesskey)
	}

}
