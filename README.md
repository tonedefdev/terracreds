![Terracreds](https://github.com/tonedefdev/terracreds/workflows/Terracreds/badge.svg?branch=v1.0.2)

<img src="https://github.com/tonedefdev/terracreds/blob/main/img/terracreds.png?raw=true" align="right" width="350" height="350">

# Terracreds
A credential helper for Terraform Cloud/Enterprise that allows secure storage of your API token within the operating system's vault instead of in a plain text configuration file

We all know storing secrets in plain text can pose major security threats, and Terraform doesn't come pre-packaged with a credential helper, so we decided to create one and to share it with the greater Terraform/DevOps community to help enable stronger security practices

#### Currently supported Operating Systems:
- [x] Windows (Credential Manager)
- [x] MacOS (Keychain)
- [x] Linux (gnome-keyring) *Tested on Ubuntu 20.04*

## Windows Install via Chocolatey
The fastest way to install `terracreds` on Windows is via our Chocolatey package:
```powershell
choco install terracreds -y
```

Once installed run the following command to verify `terracreds` was installed properly:
```powershell
terracreds -v
```

To upgrade `terracreds` to the latest version with Chocolatey run the the following command:
```powershell
choco upgrade terracreds -y
```

## macOS Install
We are currently working on a `homebrew` package, however, to install the package simply download our latest release from this repository, 
extract the package, and then place it in a directory available on `$HOME`

## Linux Install
You'll need to download the latest binary from our release page and place it anywhere on `$PATH` of your system. You can also copy and run the following commands:

```bash
wget https://github.com/tonedefdev/terracreds/releases/download/v1.1.1/terracreds_1.1.1_linux_amd64.tar.gz && \
tar -xvf terracreds_1.1.1_linux_amd64.tar.gz && \
sudo mv -f terracreds /usr/bin/terracreds && \
rm -f terracreds_1.1.1_linux_amd64.tar.gz README.md
```

The `terracreds` Linux implementation uses `gnome-keyring` in conjunction with `gnome-keyring-daemon` 
to utilize the credential storage engine

In order to leverage `terracreds` to have access to the default `Login` collection you'll need to unlock 
the collection with `gnome-keyring-daemon` using an empty password:

```bash
echo "" | gnome-keyring-daemon --unlock
```
> You do have the option of setting a password by passing it in with `echo` but every call to `terracreds get` will require the 
unlock password

The command, if successful, should return the following:
```txt
SSH_AUTH_SOCK=/run/user/1000/keyring/ssh
```

You can verify that it's running properly with:
```bash
ps -ef | grep 'gnome-keyring-daemon'
``` 

## Install From Source
Download the source files by entering the following command:
```go
go get github.com/tonedefev/terracreds 
```

Ensure you have the environment variable `GO111MODULE` enabled since this project leverages `go.mod`

For Windows:
```powershell
$env:GO111MODULE='on'
```

For macOS and Linux:
```bash
export GO111MODULE='on'
```

Once the files have been downloaded navigate to the `terracreds` directory in the and then run:
```go
go install -v
```

Navigate to the root of the project directory and you should see the `terracreds.exe` binary for Windows or `terracreds` for macOS and Linux. On Windows, copy the `.exe` to any directory of your choosing. Be sure to add the directory on `$env:PATH` for Windows to make using the application easier. On macOS and Linux we recommend you place the binary in `/usr/bin` as this directory should already be on the `$PATH` environment variable

## Upgrading
If you're upgrading to the latest version of `terracreds` from a previous version use one of the methods above to install the latest binary. Once successfully installed on your system you just need to run `terracreds generate` to copy the latest version to the correct `plugins` directory for your operating system

## Initial Configuration
In order for `terracreds` to act as your credential provider you'll need to generate the binary and the plugin directory in the default location that Terraform looks for plugins. Specifically, for credential helpers, and for Windows, the directory is `%APPDATA%\terraform.d\plugins` and for macOS and Linux `$HOME/.terraform.d/.terraformrc`

To make things as simple as possible we created a helper command to generate everthing needed to use the app. All you need to do is run the following command in `terracreds` to generate the plugin directory, and the correctly formatted binary that Terraform will use:
```bash
terracreds generate
```

This command will generate the binary as `terraform-credentials-terracreds.exe` for Windows or `terraform-credentials-terracreds` for macOS and Linux which is the valid naming convention for Terraform to recognize this plugin as a credential helper

In addition to the binary and plugin a `terraform.rc` file is required for Windows or `.terraformrc` for macOS and Linux with a `credentials_helper` block which instructs Terraform to use the specified credential helper. If you don't already have a `terraform.rc` or a `.terraformrc` file you can pass in `--create-cli-config` to create the file with the credentials helper block already generated for use with the `terracreds` binary for your OS

However, if you already have a `terraform.rc` or `.terraformrc` file you will need to add the following block to your file instead:

```hcl
credentials_helper "terracreds" {
  args = []
}
```

Once you have moved all of your tokens from this file to the `Windows Credential Manager` or `KeyChain` via `terracreds` you can remove the tokens from the file. If you don't remove the tokens, and you add the `credentials_helper` block to this file, Terraform will still use the tokens instead of `terracreds` to retreive the tokens, so be sure to remove your tokens from this file once you have used the `create` or `terraform login` command to create the credentials in `terracreds` so you can actually leverage the credential helper

## Storing Credentials
For Terraform to properly use the credentials stored in your credential manager they need to be stored a specific way. The name of the credential object must be the domain name of the Terraform Cloud or Enterprise server. For instance `app.terraform.io` which is the default name `terraform login` will use

The value for the password will correspond to the API token associated for that specific Terraform Cloud or Enterprise server

The entire process is kicked off directly from the Terraform CLI. Run `terraform login` to start the login process with Terraform Cloud. If you're using Terraform Enterprise you'll need to pass the hostname of the server as an additional argument `terraform login my.tfe.com`

You'll be sent to your Terraform Cloud instance where you'll be requested to sign-in with your account, and then sent to create an API token. Create the API token with any name you'd like for this example we'll use `terracreds`

Once completed, copy the generated token, paste it into your terminal, and then hit enter. Terraform will then leverage `terracreds` to store the credentials in the operating system's credential manager. If all went well you should receive the following success message:

```bash
Success! Terraform has obtained and saved an API token.
```

In the background `terraform` calls `terracreds` as its credential helper, `terraform` passes in a JSON token credential object, and then `terracreds` decodes that object from STDIN for storage in the operating system's credential manager. The following command is what is called by `terraform` during this process:

```bash
terraform-credentials-terracreds store app.terraform.io
```

## Verifying Credentials
When Terraform leverages `terracreds` as the credential provider it will run the following command to get the credentials value:
```bash
terraform-credentials-terracreds get app.terraform.io
```

Alternatively, you can run the same command using either binary to return the credentials. The response is formatted as a JSON object as required by Terraform to use the token:
```powershell
terracreds get app.terraform.io
```

Example output:
```json
{"token":"reallybigtokenyoudontevenknow"}
```

## Updating Credentials
To update a credential in your credential manager simply go through the same `terraform login` process and it will generate a new token and save it for you!

If the token was updated successfully the following message is returned:
```bash
Success! Terraform has obtained and saved an API token.
```

Additionally, you can check the `terracreds.log` if logging is enabled for more information

## Forgetting Credentials
You can delete the credential object at any time by running:
```bash
terraform logout
```

In the background `terraform` calls `terracreds` to perform:
```bash
terracreds forget app.terraform.io
```

If the credential was successfully deleted `terraform` will return:
```text
Success! Terraform has removed the stored API token for app.terraform.io.
```

Additionally, you can check the `terracreds.log` if logging is enabled for more information

## Protection
In order to add some protection `terracreds` adds a username to the credential object, and checks to ensure that the user requesting access to the token is the same user as the token's creator. This means that only the user account used to create the token can view the token from `terracreds` which ensures that the token can only be read by the account used to create it. Any attempt to access or modify this token from `terracreds` outside of the user that created the credentail will lead to denial messages. Additionally, if the credential name is not found, the same access denied message will be provided in lieu of a generic not found message to help prevent brute force attempts

## Logging
Wherever either binary is stored `terracreds` or `terraform-credential-terracreds` a `config.yaml` file is generated on first launch of the binary. Currently, this configuration file only enables/disables logging and sets the log path. If logging is enabled you'll find the log named `terracreds.log` at the provided path 
>It's important to note that you'll have two configuration files due to Terraform requiring that the credential helper have a very specific binary name, so when troubleshooting credential issues with Terraform remember to setup the configuration file in the `%APPDATA%\terraform.d\plugins` directory for Windows and `$HOME/.terraform.d/plugins` directory for macOS and Linux

To enable logging for Windows setup the `config.yaml` as follows:
```yaml
logging:
  enabled: true
  path: C:\Temp\
```

To enable logging for macOS and Linux:
```yaml
logging:
  enabled: true
  path: /home/username/
```

The log is helpful in understanding if an object was found, deleted, updated or added, and will be found at the path defined in the configuration file as `terracreds.log`

In addition all error messages returned by the underlying libraries will be logged when logging is enabled and an error is encountered

## Troubleshooting Linux
If you are having trouble viewing, deleting, or saving credentials on Linux systems using `gnome-keyring` you must ensure that you have unlocked the collection using `gnome-keyring-daemon --unlock` otherwise you will see the following error message in the logs:

```txt
ERROR: <TIMESTAMP> - failed to unlock correct collection '/org/freedesktop/secrets/collection/login'
```

If the daemon has unlocked the collection but you're still getting prompted for credentials -- 
check to make sure that only a single instance of the daemon is running:

```bash
ps -ef | grep gnome-keyring
```

If more than one daemon is running, take note of the pid, and use `kill` to terminate the additional daemon. Try your previous command again
and it should now be working
