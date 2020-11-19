# Terracreds
A credential helper for Terraform Cloud/Enterprise that allows secure storage of your API token within the operating system's vault instead of in a plain text configuration file. Storing secrets in plain text can pose major security threats, and Terraform doesn't come pre-packaged with a credential helper, so we decided to create one and to share it with the greater Terraform/DevOps community to help enable stronger security practices.

#### Currently supported Operating Systems:
- [x] Windows (Credential Manager)
- [x] MacOS (Keychain)
- [ ] Linux (ksecretservice or gnome-keyring)

*The Linux version is currently in development. If you'd like to support the project please feel free to submit a PR*

## Windows Install via Chocolatey
The fastest way to install `terracreds` on Windows is via our Chocolatey package:
```shell
choco install terracreds -y
```

Once installed run the following command to verify `terracreds` was installed properly:
```shell
terracreds -v
```

## Manual Install
Before downloading or installing from source set your `$env:GOPATH` to the root directory of where you plan on downloading the source files. Once completed navigate to that directory.

Download the source files by entering the following command:
```go
go get github.com/tonedefev/terracreds 
```

Once the files have been downloaded navigate to the `terracreds` directory and then run:
```go
go install -v
```

Navigate to the `bin` directory and you should see the `terracreds.exe` binary for Windows or `terracreds` for macOS. On Windows, copy the `.exe` to any directory of your choosing. Be sure to add the directory on `$env:PATH` for Windows to make using the application easier. On macOS we recommend you place the binary in `/usr/bin` as this directory should already be on the `$PATH` environment variable.

## Initial Configuration
In order for `terracreds` to act as your credential provider you'll need to generate the binary and the plugin directory in the default location that Terraform looks for plugins. Specifically, for credential helpers, and for Windows, the directory is `%APPDATA%\terraform.d\plugins` and for macOS `$HOME/.terraformrc`

To make things as simple as possible we created a helper command to generate everthing needed to use the app. All you need to do is run the following command in `terracreds` to generate the plugin directory, and the correctly formatted binary that Terraform will use:
```shell
terracreds generate
```

This command will generate the binary as `terraform-credentials-terracreds.exe` for Windows or `terraform-credentials-terracreds` for macOS which is the valid naming convention for Terraform to recognize this plugin as a credential helper.

In addition to the binary and plugin a `terraform.rc` file is required for Windows or `.terraformrc` for macOS with a `credentials_helper` block which instructs Terraform to use the specified credential helper. If you don't already have a `terraform.rc` or a `.terraformrc` file you can pass in `--create-cli-config` to create the file with the credentials helper block already generated for use with the `terracreds` binary for your OS. However, if you already have a `terraform.rc` or `.terraformrc` file you will need to add the following block to your file instead:

```hcl
credentials_helper "terracreds" {
  args = []
}
```

Once you have moved all of your tokens from this file to the `Windows Credential Manager` or `KeyChain` via `terracreds` you can remove the tokens from the file. If you don't remove the tokens, and you add the `credentials_helper` block to this file, Terraform will still use the tokens instead of `terracreds` to retreive the tokens, so be sure to remove your tokens from this file once you have used the `create` command to create the credentials in `terracreds` so you can actually leverage the credential helper.

The last configuration step is specific to Windows. You will need to add a Terraform environment variable that points to the path fo the `terraform.rc` file. Terraform's documentation states that on Windows the default location is `%APPDATA%\terraform.d\` however in our testing this wasn't the case. You can set the environment variable one of two ways:

Add the following to your PowerShell profile (`Microsoft.PowerShell_profile.ps1`) to persist this environment variable each time a PowerShell session is launched:
```powershell
$env:TF_CLI_CONFIG_FILE="$($env:APPDATA)\terraform.d\terraform.rc"
```

Manually add the environment variable as a user variable by navigating to `Control Panel > All Control Panel Items > System > Advanced system settings > Environment variables... > User variables > New...` then enter:

```txt
Variable name: TF_CLI_CONFIG_FILE
Variable value: %APPDATA%\terraform.d\terraform.rc
```

## Storing Credentials
For Terraform to properly use the credentials stored in your credential manager they need to be stored a specific way. The name of the credential object must be the domain name of the Terraform Cloud or Enterprise server. For instance `my.terraform.com`

The value for the password will correspond to the API token associated for that specific Terraform Cloud or Enterprise server.

To store the credentials you'll need to run the following command:
```shell
terracreds create -n my.terraform.com -t yourAPITokenString
```

If all went well you should receive a success message:
```
SUCCESS: Created\updated the credential object 'my.terraform.com'
```

## Verifying Credentials
When Terraform leverages `terracreds` as the credential provider it will run the following command to get the credentials value:
```shell
terraform-credentials-terracreds get my.terraform.com
```

Alternatively, you can run the same command using either binary to return the credentials. The response is formatted as a JSON object as required by Terraform to use the token:
```powershell
terracreds get my.terraform.com
```

Example output:
```json
{"token":"reallybigtokenyoudontevenknow"}
```

## Updating Credentials
To update a credential simply run the same create command and it will update the token instead:
```shell
terracreds create -n my.terraform.com -t reallybignewtoken
```

## Deleting Credentials
You can delete the credential object at any time by running:
```shell
terracreds delete -n my.terraform.com
```

## Protection
In order to add some protection `terracreds` adds a username to the credential object, and checks to ensure that the user requesting access to the token is the same user as the token's creator. This means that only the user account used to create the token can view the token from `terracreds` which ensures that the token can only be read by the account used to create it. Any attempt to access or modify this token from `terracreds` outside of the user that created the credentail will lead to denial messages. Additionally, if the credential name is not found, the same access denied message will be provided in lieu of a generic not found message to help prevent brute force attempts.

## Logging
Wherever either binary is stored `terracreds` or `terraform-credential-terracreds` a `config.yaml` file is generated on first launch of the binary. Currently, this configuration file only enables/disables logging and sets the log path. If logging is enabled you'll find the log named `terracreds.log` at the provided path. 
>It's important to note that you'll have two configuration files due to Terraform requiring that the credential helper have a very specific binary name, so when troubleshooting credential issues with Terraform remember to setup the configuration file in the `%APPDATA%\terraform.d\plugins` directory for Windows and `$HOME/.terraformrc` directory for macOS.

To enable logging for Windows setup the `config.yaml` as follows:
```yaml
logging:
  enabled: true
  path: C:\Temp\
```

To enable logging for macOS:
```yaml
logging:
  enabled: true
  path: /usr/
```

The log will be located at the defined path as `terracreds.log` 