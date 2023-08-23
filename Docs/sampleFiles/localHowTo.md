# Locally stored credentials usage

## Why
This is mainly intended for evaluating tool, but you could also use it if you are not allowed to store passwords to something like 1password. Note though that It is higly recomended that you have laptop storage medium encrypted, OS locked when you are not using it and are sole user of the said machine.

## Step to start using

### Have account list generated
There is example in this folder that contains accountList file where name in square backets will be name of the account, id aws account id and role, the name of the role that can be assumed in said account. This file should be stored in your aws folder where aws cli saves credentials file.

###  Set up MFA
We expect people to use MFA with AWS when using this tool, so you have to get secret seed that can be used to generate MFA tokens. QR Code Reader chrome extension is one option, but feel free to check other alternatives. Follow GUI instructions to point when you have to read QR code with your authenticator device. Instead of using mobile app, use QR code reader and copy code between secret=  and end it with & charracter. Now go to app. select "connect" and  connect via local settings. Add seed you read with QR code. Select Show MFA code. This should show you your MFA code that you need to input on AWS console. Wait awhile to get second token as two are needed. You can use copy button function, so you don't have to manually write token down. After this is done. Please also copy your MFA device id from same AWS consol page.

### Rest of the settings

You should then copy your accesskey and secretaccess keys and default region you want to use. Also add meaningful alias, as it helps other people to see who is actually performing actions with out going too deep into cloudtrail logs. I personally use my email address. You have to press connect after that. Then you should be done and be able to list accounts in your accountlist file with the tool. 


