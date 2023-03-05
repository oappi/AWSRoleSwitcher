# AWSRoleSwitcher
Note that AWSRoleSwitcher has absolutely no assosication with AWS or amazon in general. Only reason it is in the name is that AWS cloud is only cloud it works with.

Tool is meant to make it easier to switch roles when you have to use iam user instead of SSO to login to your accounts. Tool is best used when you have AWS accounts just for iam users from which you switch to "project" accounts. Techincally you can add role to your IAM account there you federate users, and build project there, but at least when writing this it is not AWS best practise.


## What type of AWS setup you need
IAM user based access control. This wont work with AWS SSO.


# connection setup instructions
[1Password Credential Setup](other_file.md)
[LocalCredential Setup](other_file.md)


## How it works
You connect to credential provider like 1password or localcredentials. For 1password you save your credentials in 1password and this app will use 1opv2 cli tool to fetch them, including MFA. You can also not use 1password and save credentials locally. Link to HowTo's is inconnector setup chapter. 

After you have "connected" you can pick session lenght. Default is one hour after which cli and browser session will expire. You can create new session with reconnect button. Note that application does not create new session when you select session duration. This is mainly because unlike with selecting account, there are some rules that would have to be in place for this to work as expected that has not been implemented, thus workflow for now is that you select duration and then from account dropdown account or if you want to continue using same account press reconnect

When you want to switch to account from second dropdown pick what ever account you want to assume. Note that you can write on this field to filter accounts by name, id or assumed role name. example: account name is testin-account-3. you can just write account and it will match to testing-account-3. Imediatly when you click account name it will sign into account.

### CLI usage
We overwrite default entry in users aws credential file. This is what AWS cli uses as default if nothing else has been defined in command options. This is done each time user picks account from dropdown or presses reconnect button

### Open AWS Console
When you want to use browser open browser button opens AWS console in your last used browser. This also allows you to sign in with multiple chrome profiles to different accounts or with totally different browsers.

### Rotate credentials
It is recomended to rotate accesskeys every now and then, so this is under advance settings. You should only use it if you have 1 accesskey enabled as of writing this you could only have 2. Process needs 2 as we are removing old one after we have checked that accesskey is actually stored.


