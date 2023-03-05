// This file basically includes GUI and minimal calls to logic
package main

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	keyrotation "github.com/oappi/awsroler/accesskeyRotation"
	idp "github.com/oappi/awsroler/awsLogic"
	creds "github.com/oappi/awsroler/credentialFileLogic"
	"github.com/oappi/awsroler/interfaces"
	"github.com/oappi/awsroler/sharedStructs"
)

var awsSession sharedStructs.SessionInfo
var lock = &sync.Mutex{}

var localWriter = interfaces.IniLogic{AWSFolderLocation: creds.GetAWSFolderStripError()}
var gregion = ""
var accountsList = []string{"Connect to crendetial service first"}
var gOptionSelection *widget.SelectEntry
var SettingsInterface interfaces.SettingsInterface
var SettingsObject sharedStructs.FederationAccountSettingsObject
var selectedSessionTime = "1 hour session"

// sessionInfo
/*
func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}*/

func main() {
	a := app.NewWithID("io.fyne.oappi.AWSRoleSwitcher")
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("AWS Role Switcher")
	w.Resize(fyne.NewSize(550, 480))

	_, err := creds.GetAWSFolder(runtime.GOOS)
	if err != nil {
		errorPopUp(a, "OS not Supported")
	}

	settingsItem := fyne.NewMenuItem("GUI Settings", func() {
		w := a.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	})

	connect1passwordItem := fyne.NewMenuItem("Connect via 1password", func() {
		go show1PSettings(a)
	})

	connectlocalSettings := fyne.NewMenuItem("Connect via local settings", func() {
		go showLocalSettings(a)
	})

	advancedMenu := fyne.NewMenu("Advanced",
		fyne.NewMenuItem("Rotate Accesskey", func() {
			go showKeyRotation(a)
		}))

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Info", func() {
			go showAuthor(a)
		}))

	connectMenu := fyne.NewMenu("Connect", connect1passwordItem)
	file := fyne.NewMenu("File", settingsItem)
	connectMenu.Items = append(connectMenu.Items, fyne.NewMenuItemSeparator(), connectlocalSettings)
	//file.Items = append(file.Items, fyne.NewMenuItemSeparator(), guiSettingsItem)

	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		file,
		connectMenu,
		advancedMenu,
		helpMenu,
	)
	w.SetMainMenu(mainMenu)
	w.SetMaster()

	accountName := widget.NewLabel("not set")
	accountName.TextStyle.Bold = true
	accountName.TextStyle.Italic = true
	accountName.Alignment = fyne.TextAlignLeading
	reconnectButton := widget.NewButton("Reconnect", func() {
		stsSettings, stsError := getStSConfig(SettingsObject)
		if stsError != nil {
			popError(a, stsError)
		} else {
			connectAccount(stsSettings, accountName.Text, localWriter, selectedSessionTime)
		}
	})
	openBrowserButton := widget.NewButton("Open in Browser", func() {
		idp.LoginBrowser(accountName.Text, awsSession, SettingsInterface)
	})

	reconnectButton.Importance = 0
	//reconnectButton.Importance = 1
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	allElements, _ := filteredListForSelect(nil)
	accountSelectEntry := widget.NewSelectEntry(allElements)
	timerOptions := []string{"1 hour session", "2 hour session", "4 hour session", "8 hour session", "12 hour session"}
	timerSelectEntry := widget.NewSelectEntry(timerOptions)
	timerSelectEntry.SetPlaceHolder("1 hour session")
	gOptionSelection = accountSelectEntry
	accountSelectEntry.PlaceHolder = "Type or select an account"
	accountSelectEntry.OnChanged = func(input string) {

		filteredList, match := filteredListForSelect(&input)
		accountSelectEntry.SetOptions(filteredList)
		if match {
			stsSettings, stsError := getStSConfig(SettingsObject)
			if stsError != nil {
				popError(a, stsError)
			} else {
				accountName.SetText(accountSelectEntry.Text)
				accountName.Alignment = fyne.TextAlignLeading
				connectAccount(stsSettings, accountName.Text, localWriter, selectedSessionTime)
				accountName.Alignment = fyne.TextAlignCenter
			}

		}
	}
	timerSelectEntry.OnChanged = func(input string) {

		_, match := filteredCustomListForSelect(&input, timerOptions)
		if match {
			selectedSessionTime = input
		}
	}
	//openBrowserButton
	acclabelOpenBrowser := container.NewVSplit(accountName, openBrowserButton)
	bottomComponents := container.NewVSplit(acclabelOpenBrowser, reconnectButton)
	searchselect := container.NewVSplit(accountSelectEntry, bottomComponents)
	timeselector := container.NewVSplit(timerSelectEntry, searchselect)
	timeselector.Offset = 0.1
	searchselect.Offset = 0.1

	//searchselect := container.NewAdaptiveGrid(1, optionSelectEntry, accountName)
	w.SetContent(timeselector)
	w.Resize(fyne.NewSize(240, 260))

	w.ShowAndRun()
}

func showLocalSettings(a fyne.App) {
	win := a.NewWindow("Local Connect Settings")
	MFASeedLabel := widget.NewLabel("MFA seed")
	MFASeedText := widget.NewEntry()
	regionLabel := widget.NewLabel("Region")
	regionListText := widget.NewEntry()
	MFADeviceLabel := widget.NewLabel("MFA Device")
	MFADeviceText := widget.NewEntry()
	MFACodeButtonLabel := widget.NewLabel("MFA code")
	MFACodeButton := widget.NewButton("Show MFA code", func() {
		go showLocalMFA(a, OverRideSavedIfUserGivesInput(MFASeedText.Text, ""))
	})
	AccessKeyLabel := widget.NewLabel("AccessKey")
	AccessKeyText := widget.NewPasswordEntry()
	SecretAccessKeyLabel := widget.NewLabel("SecretAccessKey")
	SecretAccessKeyText := widget.NewPasswordEntry()
	aliasLabel := widget.NewLabel("MyAssumeAlias")
	aliasText := widget.NewEntry()
	accountListLocation := creds.GetAWSFolderStripError() + "accountlist"

	mfaDevice, mfaSeed, accesskey, secretaccesskey, alias, region, fetcherror := localWriter.GetLocalSettings()
	if fetcherror != nil {
		MFADeviceText.SetPlaceHolder(mfaDevice)
		MFASeedText.SetPlaceHolder(mfaSeed)
		AccessKeyText.SetPlaceHolder(accesskey)
		SecretAccessKeyText.SetPlaceHolder(secretaccesskey)
		aliasText.SetPlaceHolder(alias)
	}
	labels := container.NewGridWithColumns(1, MFASeedLabel, MFACodeButtonLabel, MFADeviceLabel, regionLabel, AccessKeyLabel, SecretAccessKeyLabel, aliasLabel)
	textFields := container.NewGridWithColumns(1, MFASeedText, MFACodeButton, MFADeviceText, regionListText, AccessKeyText, SecretAccessKeyText, aliasText)
	settingscontainer := container.NewGridWithColumns(2, labels, textFields)

	applySettingsButton := widget.NewButton("Connect", func() {
		MFASeedOption := OverRideSavedIfUserGivesInput(MFASeedText.Text, mfaSeed)
		MFADeviceOption := OverRideSavedIfUserGivesInput(MFADeviceText.Text, mfaDevice)
		RegionOption := OverRideSavedIfUserGivesInput(regionListText.Text, region)
		AccessKeyOption := OverRideSavedIfUserGivesInput(AccessKeyText.Text, accesskey)
		SecretAccessKeyOption := OverRideSavedIfUserGivesInput(SecretAccessKeyText.Text, secretaccesskey)
		UserAliasOption := OverRideSavedIfUserGivesInput(aliasText.Text, alias)
		AccountListLocationOption := OverRideSavedIfUserGivesInput(accountListLocation, accountListLocation)
		SettingsInterface = interfaces.LocalSettings{Lock: lock, MFASeed: MFASeedOption, MFADevice: MFADeviceOption, Region: RegionOption, AccessKey: AccessKeyOption, SecretAccessKey: SecretAccessKeyOption, UserAlias: UserAliasOption, AccountListLocation: AccountListLocationOption, AWSFolderLocation: creds.GetAWSFolderStripError(), LocalWriter: localWriter}
		err := updateSettings(SettingsInterface)
		if err != nil {
			popError(a, err)

		} else {
			localWriter.SetLocalSettings(MFADeviceOption, MFASeedOption, AccessKeyOption, SecretAccessKeyOption, UserAliasOption, RegionOption)
			win.Close()
		}
	})

	settingsplit := container.NewVSplit(settingscontainer, applySettingsButton)
	settingsplit.Offset = 0.9
	win.SetContent(settingsplit)
	win.Show()
	win.Close()
}

func show1PSettings(a fyne.App) {
	win := a.NewWindow("1Password Settings")
	savedDomain, savedEntity := localWriter.Get1PasswordSettings()
	domainLabel := widget.NewLabel("1Password domain")
	domainText := widget.NewEntry()
	domainText.SetPlaceHolder(savedDomain)
	entityNameLabel := widget.NewLabel("AWS login item name")
	entityNameText := widget.NewEntry()
	entityNameText.SetPlaceHolder(savedEntity)
	passwordLabel := widget.NewLabel("1password Password")
	passwordText := widget.NewPasswordEntry()

	labels := container.NewGridWithColumns(1, domainLabel, entityNameLabel, passwordLabel)
	textFields := container.NewGridWithColumns(1, domainText, entityNameText, passwordText)
	settingscontainer := container.NewGridWithColumns(2, labels, textFields)

	applySettingsButton := widget.NewButton("Apply", func() {
		var OPEntity = ""
		var OPDomain = ""
		// here we use saved value if user does not overwrite it
		if entityNameText.Text == "" {
			OPEntity = savedEntity
		} else {
			OPEntity = entityNameText.Text
		}
		if domainText.Text == "" {
			OPDomain = savedDomain
		} else {
			OPDomain = domainText.Text
		}
		SettingsInterface = interfaces.Onepassword{Lock: lock, Uuid: OPEntity, OPDomain: OPDomain, Password: passwordText.Text}
		err := updateSettings(SettingsInterface)
		if err != nil {
			popError(a, err)
		} else {
			localWriter.Set1PasswordSettings(domainText.Text, entityNameText.Text)
			win.Close()
		}
	})

	settingsplit := container.NewVSplit(settingscontainer, applySettingsButton)

	settingsplit.Offset = 0.9
	win.SetContent(settingsplit)
	win.Show()
	win.Close()
}

func showAuthor(a fyne.App) {
	win := a.NewWindow("Info")
	win.SetContent(widget.NewLabel("\n Copyright 2021 Ossi Ala-Peijari (MIT lincese) \n\n\n" +
		"Permission is hereby granted, free of charge, to any person obtaining a copy \n of this software and associated documentation files  (the \"Software\"),\n to deal in the Software without restriction, including without\n limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,\n and/or sell copies of the Software, and to permit persons to whom \n the Software is furnished to do so, subject to the following conditions: \n\n\n" +

		"The above copyright notice and this permission notice shall be included in all\n copies or substantial portions of the Software. \n\n\n" +

		"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, \n INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, \n FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT \n HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, \n WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, \n OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS \n IN THE SOFTWARE. \n\n\n" +
		"Additional licenses: https://github.com/fyne-io/fyne/blob/master/LICENSE"))
	win.Resize(fyne.NewSize(200, 200))
	win.Show()
	win.Close()
}

func showLocalMFA(a fyne.App, MFAsecret string) {
	win := a.NewWindow("MFA")
	MFACode := widget.NewLabel("")
	Timer := widget.NewLabel("Timer")
	TimerLabel := widget.NewLabel("Time left")
	MFACodeLabel := widget.NewLabel("MFA")
	copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		win.Clipboard().SetContent(MFACode.Text)
	})
	Empty := widget.NewLabel("")
	mfaContainer := container.NewGridWithColumns(1, MFACode)
	timercontainer := container.NewGridWithColumns(1, Timer)

	labels := container.NewGridWithColumns(1, container.NewCenter(MFACodeLabel), container.NewCenter(TimerLabel))
	infos := container.NewGridWithColumns(1, container.NewCenter(mfaContainer), container.NewCenter(timercontainer))
	actions := container.NewGridWithColumns(1, container.NewCenter(copyButton), Empty, Empty)
	settingscontainer := container.NewGridWithColumns(3, labels, infos, actions)

	win.SetContent(settingscontainer)
	win.Show()
	win.Close()

	forever := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				forever <- struct{}{}
				return
			default:
				var unixtime = time.Now().Unix()
				var timeLeft = 30 - unixtime%30
				var mfaTimePast = 30 - timeLeft
				var unixtimeOfGettingMFA = unixtime - mfaTimePast //we only fetch every 30s
				var stringTimeLeft = strconv.FormatInt(timeLeft, 10)
				Timer.SetText(stringTimeLeft)
				if timeLeft == 30 || MFACode.Text == "" {
					otp, error := creds.GetOTP(MFAsecret, unixtimeOfGettingMFA/30)
					if error != nil {
						MFACode.SetText("incorrect mfa secret")
					} else {
						MFACode.SetText(otp)
					}
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}(ctx)

	<-forever

	win.SetOnClosed(func() {
		cancel()
	})
}

func showKeyRotation(a fyne.App) {
	win := a.NewWindow("Key Rotation")
	infoLabel := widget.NewLabel("This feature creates new accesskeys," +
		"\nsaves it to your storage medium (1password locally etc)" +
		"\nand then removes old key. Note that this function only " +
		"\nworks if you have 1 accesskey on your iam account as aws limits " +
		"\nkeys to two")

	textField := container.NewGridWithColumns(1, infoLabel)
	infocontainer := container.NewGridWithColumns(1, textField)

	applySettingsButton := widget.NewButton("Rotate", func() {
		err, newAccesskey, newSecretAccesskey := keyrotation.RotateAccesskeys(SettingsInterface)
		if err != nil {
			var errormessage = err.Error()
			go errorPopUp(a, errormessage)
		}
		SettingsObject.AccessKey = newAccesskey
		SettingsObject.SecretAccessKey = newSecretAccesskey
		win.Close()
	})

	settingsplit := container.NewVSplit(infocontainer, applySettingsButton)

	settingsplit.Offset = 0.9
	win.SetContent(settingsplit)
	win.Show()
	win.Close()
}

func popError(a fyne.App, err error) {
	var errormessage = err.Error()
	go errorPopUp(a, errormessage)
}

func errorPopUp(a fyne.App, message string) {

	win := a.NewWindow("Error")
	infoLabel := widget.NewLabel(message)

	textField := container.NewGridWithColumns(1, infoLabel)
	infocontainer := container.NewGridWithColumns(1, textField)

	applySettingsButton := widget.NewButton("Close", func() {
		win.Close()
	})

	settingsplit := container.NewVSplit(infocontainer, applySettingsButton)

	settingsplit.Offset = 0.9
	win.SetContent(settingsplit)
	win.Show()
	win.Close()
}
