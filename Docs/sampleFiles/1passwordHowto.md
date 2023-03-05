# Locally stored credentials usage

## Why
This helped app was mainly designed to use with 1password, and rest are bit of a after though. It was tool for me personally to easily switch account on customer environments with more than 100 aws accounts. After which AWS-MFA and other command line tools are just not enjoyable to use

## Step to start using
//TODO link to local
//TODO link to 1password

### Have account list generated
you should generate similar list of your accounts as example in accountList file.

### Set up AWS account credentials to 1password
Like you would usually add AWS credentials, add them into your 1password by using 1password browser extension to read qr code aws presents. 1Password should then handle token service for you.

Add following text fields:
MFADevice : with arn of your MFA device
alias: name that will be shown when you interact with the role. I recomend something that makes identifying you easier.
Region: Set default region like eu-west-1
Accounts: copy your accountList file here. It will be fetched to tool when you connect via 1password connect.

Then set following as "password" so others wont be able to read it from your screen.
access_key
secret_access_key
and copy your aws accesskey and secretaccesskey values to them.
Now your setup should be completed. I would recomend removing accountlist file from local pc as well as accesskeys. They should be stored and versioned in 1password for now on.

### Connecting with 1password

Open add, select "connect" and "Connect via 1password". Note that you have to have 1password cliv2 installed as it will be used to issue commands towards 1password. You should have logged in with it at least once. If you don't have organizational account your domain is "my". AWS login item name is name of entry in your 1password. 1password Password is your password for 1password. 

## Recomendations
We really recomend that you only test functionality with localcredentials, and shift to 1password for more production type of workloads as it has benefit of not storing credentials on your PC. Especially if you are sharing your PC with others.Especially if you are not only user please consider aws-vault instead https://github.com/99designs/aws-vault .Localcredentials with this tool are not crypted in anyway, so someone can copy your credentials if they have access to your files since MFA is in plaintext MFA device can be duplicate.

1password is kinda cool with this app as you do not have to have configurations on your local device to be able to connect your AWS accounts and only required short time credentials are stored similar way AWS stores them. This means that as long as you have required apps installed, you can switch PC really easily. In addition, accesskey rotation has safeguards of failed rotation as old AWS key wont be deleted until key has successfully stored in 1password. Incase of failure you can always retrive old key that should be usable because of the check before delete. This checking is also main reason why key rotation takes awhile.

### Summary
 * Localcredentials might be ok if you lock your pc when you leave it unattended and have drive encryption eneabled. Otherwise please do not use with production accounts
 * Please use 1password with this tool if you can.