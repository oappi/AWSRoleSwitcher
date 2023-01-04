package awsLogic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/oappi/awsroler/interfaces"
	"github.com/oappi/awsroler/sharedStructs"
)

func LoginBrowser(selectedAccountInfo string, sessionInfo sharedStructs.SessionInfo, SettingsInterface interfaces.SettingsInterface) error {
	region, errRegion := SettingsInterface.GetRegion()
	if errRegion != nil {
		return errRegion
	}
	loginURLPrefix, destination := GenerateLoginURL(region, "")
	req, err := http.NewRequest("GET", loginURLPrefix, nil)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(map[string]string{
		"sessionId":    sessionInfo.Accesskey,
		"sessionKey":   sessionInfo.SecretAccessKey,
		"sessionToken": sessionInfo.Token,
	})
	q := req.URL.Query()
	q.Add("Action", "getSigninToken")
	q.Add("Session", string(jsonBytes))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Response body was %s", body)
		fmt.Errorf("Call to getSigninToken failed with %v", resp.Status)
	}
	var respParsed map[string]string
	err = json.Unmarshal([]byte(body), &respParsed)
	if err != nil {
		return err
	}

	signinToken, ok := respParsed["SigninToken"]
	if !ok {
		fmt.Errorf("Expected a response with SigninToken")
	}

	loginURL := fmt.Sprintf("%s?Action=login&Issuer=AWS-Roller&Destination=%s&SigninToken=%s",
		loginURLPrefix, url.QueryEscape(destination), url.QueryEscape(signinToken))
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", loginURL).Start()
		if err != nil {
			return err
		}
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", loginURL).Start()
		if err != nil {
			return err
		}
	case "darwin":
		err = exec.Command("open", loginURL).Start()
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("unsupported platform")
		if err != nil {
			return err
		}
	}

	return nil
}

// shout out to https://github.com/99designs/aws-vault/blob/master/cli/login.go
func GenerateLoginURL(region string, path string) (string, string) {
	loginURLPrefix := "https://signin.aws.amazon.com/federation"
	destination := "https://console.aws.amazon.com/"

	if region != "" {
		destinationDomain := "console.aws.amazon.com"
		switch {
		case strings.HasPrefix(region, "cn-"):
			loginURLPrefix = "https://signin.amazonaws.cn/federation"
			destinationDomain = "console.amazonaws.cn"
		case strings.HasPrefix(region, "us-gov-"):
			loginURLPrefix = "https://signin.amazonaws-us-gov.com/federation"
			destinationDomain = "console.amazonaws-us-gov.com"
		}
		if path != "" {
			destination = fmt.Sprintf("https://%s.%s/%s?region=%s",
				region, destinationDomain, path, region)
		} else {
			destination = fmt.Sprintf("https://%s.%s/console/home?region=%s",
				region, destinationDomain, region)
		}
	}
	return loginURLPrefix, destination
}
