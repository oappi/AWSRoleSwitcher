package awsLogic

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

var packsize int64 = 3
var sourceIdentity string = "localhost"
var arn string = "arn:aws:iam::123456789012:role/service-role/QuickSightAction"
var roleid string = "roleid"
var AccessKeyId = "12345678910"
var Expiration = time.Now()
var SecretAccessKey = "secretaccesskey"
var SessionToken = "sessionToken"

type STSConfigMock struct {
	STS stsiface.STSAPI
}

func (m *STSConfigMock) AssumeRole(*sts.AssumeRoleInput) (
	*sts.AssumeRoleOutput,
	error,
) {
	return &sts.AssumeRoleOutput{
		AssumedRoleUser:  &sts.AssumedRoleUser{Arn: &arn, AssumedRoleId: &roleid},
		Credentials:      &sts.Credentials{AccessKeyId: &AccessKeyId, Expiration: &Expiration, SecretAccessKey: &SecretAccessKey, SessionToken: &SessionToken},
		PackedPolicySize: &packsize,
		SourceIdentity:   &sourceIdentity,
	}, nil
}

func TestGetAsumeRoleCredentials(t *testing.T) {

}
