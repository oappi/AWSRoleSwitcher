package credentialFileLogic

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

func GetOTP(secret string, interval int64) (string, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", errors.New("Failed to convert mfa-seed to base32")
	}
	byteArray := make([]byte, 8)
	binary.BigEndian.PutUint64(byteArray, uint64(interval))
	hash := hmac.New(sha1.New, key)
	hash.Write(byteArray)
	hs := hash.Sum(nil)
	ob := (hs[19] & 15)
	var header uint32
	br := bytes.NewReader(hs[ob : ob+4])
	errb := binary.Read(br, binary.BigEndian, &header)
	if errb != nil {
		return "", errors.New("Failed generate hash from MFA seed")
	}
	h12 := (int(header) & 0x7fffffff) % 1000000
	otp := strconv.Itoa(int(h12))

	return AddPadding(otp), nil

}

func AddPadding(otp string) string {
	if len(otp) == 6 {
		return otp
	}
	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}
	return otp
}
