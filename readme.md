# AWSRoleSwitcher
Note that AWSRoleSwitcher has no assosication with Amazon Web Services or Amazon in general. Name contains AWS to define which cloud provider this tool works.

Tool is meant to make it easier to switch roles when you have to use iam user instead of something like AWS SSO to login to your accounts. This tool provides minimalistic, but intuitive GUI tool for this kind of setup. If you use external credential provider such as 1Password, you will only be storing temporary credentials locally while MFA and longtime credentials are read from external credential provider when needed.




## What type of AWS setup you need
Set up this tool expects is AWS account containing all IAM users and this account has been set up to assume role to different accounts. We also expect users to use MFA and calls will most likely fail with out it
This wont work with AWS SSO as we expect IAM credentials to be stored somewhere tool can read them. 


# Connection setup instructions
Currently tool supports both 1Password and local storing of credentials. I recommend you use 1Password as storing credentials to unencrypted files is not that smart thing to do and it is meant only to teste this tool. If you don't want to use external credential providers that are supported I would suggest you try something  like aws-vault instead

[1Password Credential Setup](other_file.md)
[LocalCredential Setup](other_file.md)


## How it works
You connect to credential provider like 1password or localcredentials. For 1password you save your credentials in 1password and this app will use 1password cli tool to fetch them, including MFA. You can also save credentials locally for testing purposes. Link to HowTo's is in Connection setup instructions. 

After you have "connected" to your credential provider you can pick session lenght. Default is one hour after which cli and browser session will expire. You can create new session with reconnect button. Note that application does not create new session when you select session duration. This is mainly because unlike with selecting account, there are some rules that would have to be in place for this to work as expected that has not been implemented, thus workflow for now is that you select duration and then from account dropdown account or if you want to continue using same account press reconnect

When you want to switch to account from second dropdown pick what ever account you want to assume. Note that you can write on this field to filter accounts by name, id or assumed role name. example: account name is testin-account-3. you can just write account and it will match to testing-account-3. Imediatly when you click account name it will sign into account.

###  AWS CLI usage
We overwrite default entry in users aws credential file so you don't have to define profile in your commands. You can use AWS CLI as you would normally

### Open AWS Console in browser
When you want to use browser open browser button opens AWS console in your last used browser. This also allows you to sign in with multiple chrome profiles to different accounts or with totally different browsers. This feature works in a way that we create token with AWS credentials that we submit to AWS when we log into console

### Rotate credentials
It is recomended to rotate accesskeys every now and then, so this is under advance settings. You should only use it if you have 1 accesskey enabled as of writing this you could only have 2. Process needs 2 as we are removing old one after we have checked that accesskey is actually stored. As this tool does not currently save old version of credentials file rotation has been disabled when storing longtime AWS credentials locally


