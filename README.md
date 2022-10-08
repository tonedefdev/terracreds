![Terracreds](https://github.com/tonedefdev/terracreds/workflows/Terracreds/badge.svg?branch=v1.0.2)

<img src="https://github.com/tonedefdev/terracreds/blob/main/img/terracreds.png?raw=true" align="right" width="350" height="350">

# Terracreds
A credential helper for Terraform Automation and Collaboration Software, or to store any other secrets, securely in the operating system's credential vault or through a third party vault provider. No longer keep secrets in a plain text configuration file!

We all know storing secrets in plain text can pose major security threats, and Terraform doesn't come pre-packaged with a credential helper, so we decided to create one and to share it with the greater Terraform/DevOps community to help enable stronger security practices.

#### Currently supported Operating Systems:
- [x] Windows (Credential Manager)
- [x] MacOS (Keychain)
- [x] Linux (gnome-keyring) *Tested on Ubuntu 20.04*

#### Currently supported Vault providers:
- [x] AWS Secrets Manager
- [x] Azure Key Vault
- [x] Google Secret Manager 
- [x] HashiCorp Vault

#### Currently Supported Terraform Automation and Collaboration Software:
- [x] env0
- [x] Scalr
- [x] Spacelift
- [x] Terraform Cloud
- [x] Terraform Enterprise

## Quick Links
- Install & Configure
  - [Windows](https://github.com/tonedefdev/terracreds#windows-install-via-chocolatey)
  - [macOS](https://github.com/tonedefdev/terracreds#macos-install)
  - [Linux](https://github.com/tonedefdev/terracreds#linux-install)
  - [From Source](https://github.com/tonedefdev/terracreds#install-from-source)
  - [Upgrading](https://github.com/tonedefdev/terracreds#upgrading)
  - [Initial Configuration](https://github.com/tonedefdev/terracreds#initial-configuration)
- Usage
  - [Storing](https://github.com/tonedefdev/terracreds#storing-credentials)
  - [Verifying](https://github.com/tonedefdev/terracreds#storing-credentials)
  - [Updating](https://github.com/tonedefdev/terracreds#updating-credentials)
  - [Forgetting](https://github.com/tonedefdev/terracreds#forgetting-credentials)
  - [Listing](https://github.com/tonedefdev/terracreds#list-credentials)
- Vault Providers
  - [General Setup](https://github.com/tonedefdev/terracreds#setting-up-a-vault-provider)
  - [AWS Secrets Manager](https://github.com/tonedefdev/terracreds#aws-secrets-manager)
  - [Azure Key Vault](https://github.com/tonedefdev/terracreds#azure-key-vault)
  - [Google Secret Manager](https://github.com/tonedefdev/terracreds#google-secret-manager)
  - [HashiCorp Vault](https://github.com/tonedefdev/terracreds#hashicorp-vault)
- Miscellaneous
  - [Protection](https://github.com/tonedefdev/terracreds#protection)
  - [Logging](https://github.com/tonedefdev/terracreds#logging)
- Troubleshooting
  - [Known Issues](https://github.com/tonedefdev/terracreds#known-issues)
  - [Linux](https://github.com/tonedefdev/terracreds#linux)

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
We are currently working on a `homebrew` package, however, you can leverage this shell script to install `terracreds`

```sh
curl https://github.com/tonedefdev/terracreds/releases/download/v2.1.2/terracreds_2.1.2_darwin_amd64.tar.gz -o terracreds_2.1.2_darwin_amd64.tar.gz  && \
tar -xvf terracreds_2.1.2_darwin_amd64.tar.gz && \
sudo mv -f terracreds /usr/bin/terracreds && \
rm -f terracreds_2.1.2_darwin_amd64.tar.gz README.md
```

## Linux Install
You'll need to download the latest binary from our release page and place it anywhere on `$PATH` of your system. You can also copy and run the following commands:

```bash
wget https://github.com/tonedefdev/terracreds/releases/download/v2.1.2/terracreds_2.1.2_linux_amd64.tar.gz && \
tar -xvf terracreds_2.1.2_linux_amd64.tar.gz && \
sudo mv -f terracreds /usr/bin/terracreds && \
rm -f terracreds_2.1.2_linux_amd64.tar.gz README.md
```

The `terracreds` Linux implementation uses `gnome-keyring` in conjunction with `gnome-keyring-daemon` 
to utilize the credential storage engine.

In order to leverage `terracreds` to have access to the default `Login` collection you'll need to unlock 
the collection with `gnome-keyring-daemon` using an empty password:

```bash
echo "" | gnome-keyring-daemon --unlock
```
> You do have the option of setting a password by passing it in with `echo` but every call to `terracreds get` will require the 
unlock password.

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

Navigate to the root of the project directory and you should see the `terracreds.exe` binary for Windows or `terracreds` for macOS and Linux. On Windows, copy the `.exe` to any directory of your choosing. Be sure to add the directory on `$env:PATH` for Windows to make using the application easier. On macOS and Linux we recommend you place the binary in `/usr/bin` as this directory should already be on the `$PATH` environment variable.

## Upgrading
If you're upgrading to the latest version of `terracreds` from a previous version use one of the methods above to install the latest binary. Once successfully installed on your system you just need to run `terracreds generate` to copy the latest version to the correct `plugins` directory for your operating system.

## Initial Configuration
In order for `terracreds` to act as your credential provider you'll need to generate the binary and the plugin directory in the default location that Terraform looks for plugins. Specifically, for credential helpers, and for Windows, the directory is `%APPDATA%\terraform.d\plugins` and for macOS and Linux `$HOME/.terraform.d/.terraformrc`.

To make things as simple as possible we created a helper command to generate everthing needed to use the app. All you need to do is run the following command in `terracreds` to generate the plugin directory, and the correctly formatted binary that Terraform will use:
```bash
terracreds generate
```

This command will generate the binary as `terraform-credentials-terracreds.exe` for Windows or `terraform-credentials-terracreds` for macOS and Linux which is the valid naming convention for Terraform to recognize this plugin as a credential helper.

In addition to the binary and plugin a `terraform.rc` file is required for Windows or `.terraformrc` for macOS and Linux with a `credentials_helper` block which instructs Terraform to use the specified credential helper. If you don't already have a `terraform.rc` or a `.terraformrc` file you can pass in `--create-cli-config` to create the file with the credentials helper block already generated for use with the `terracreds` binary for your OS.

However, if you already have a `terraform.rc` or `.terraformrc` file you will need to add the following block to your file instead:

```hcl
credentials_helper "terracreds" {
  args = []
}
```

Once you have moved all of your tokens from this file to your preferred vault provider via `terracreds` you can remove the tokens from the file. If you don't remove them, but you add the `credentials_helper` block to this file, Terraform will still use the token from this file instead of from the vault configured with `terracreds`.

## Storing Credentials
For Terraform to properly use the credentials stored in your credential manager they need to be stored a specific way. The name of the credential object must be the domain name of the Terraform Cloud or Enterprise server. For instance `app.terraform.io` which is the default name `terraform login` will leverage.

The value for the password will correspond to the API token associated for that specific Terraform Cloud or Enterprise server.

The entire process is kicked off directly from the Terraform CLI. Run `terraform login` to start the login process with Terraform Cloud. If you're using Terraform Enterprise you'll need to pass the hostname of the server as an additional argument `terraform login my.tfe.com`.

You'll be sent to your Terraform Cloud or Enterprise Software instance where you'll be requested to sign-in with your account, and then sent to create an API token. Create the API token with any name you'd like for this example we'll use `terracreds`.

Once completed, copy the generated token, paste it into your terminal, and then hit enter. Terraform will then leverage `terracreds` to store the credentials in the operating system's credential manager. If all went well you should receive the following success message:

```bash
Success! Terraform has obtained and saved an API token.
```

In the background `terraform` calls `terracreds` as its credential helper, `terraform` passes in a JSON token credential object, and then `terracreds` decodes that object from STDIN for storage in the operating system's credential manager. The following command is what is called by `terraform` during this process:

```bash
terraform-credentials-terracreds store app.terraform.io
```

If you prefer, you can also create credentials manually by running:
```bash
terracreds create -n app.terraform.io -v <TACOS_API_TOKEN>
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

You can also run `terracreds update -n my-secret -v my-secret-value` to update a secret value.

Additionally, you can check the `terracreds.log` if logging is enabled for more information.

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

You can also run `terracreds delete -n app.terraform.io` if you want to manually remove the credential.

Additionally, you can check the `terracreds.log` if logging is enabled for more information.

## List Credentials
> New in version `2.1.0`

You can pass in a comma separated list of secrets to print out the secret values to the screen:
```bash
terracreds list -l mysecret,mysecret2
```

You can also setup a list of secrets in the configuration file by using:
 ```bash
 terracreds config secrets -l mysecret,mysecret2
 ``` 

 To print out the secrets from the names stored in the configuration file:
 ```bash
 terracreds list --from-config
 ```

There's a helper flag `--as-tfvars` which will return the secret values formatted for use with `terraform`. Depending on the shell calling this command will determine how you can readily use these values.

For instance on Linux/macOS you can simply call `eval` to evaluate the output to then convert the returned values into variables in your current shell.

Also, by default, `terracreds` will convert any dashes `[-]` in a secret name with underscores `[_]` since this is the typical variable naming style convention in Terraform. However, you can override that behavior by passing in an override flag with any string value you'd prefer to use:
```bash
terracreds list --as-tfvars --override-replace-string -
```

The above example would maintain the dash `[-]` in the outuput of the formatted TF_VARS instead of replacing it by the default underscore `[_]`

Additionally, you can use `--as-json` to return the secret names and values as a JSON string. This is printed to standard output so you can make use of shell pipes and other commands to ingest the data.

## Setting Up a Vault Provider
> We have example [terraform](https://github.com/tonedefdev/terracreds/tree/main/terraform) code you can reference in order to setup your `AWS` or `Azure` VMs to use `terracreds` for a CI/CD piepline agent or a development workstation.

> New in version `2.1.0`

All of the external vault providers now make use of the provider's default credential authentication schemes. Please, refer to the documentation for each provider's default authentication mechanisms for more information on what options are availabile, and what is required to set up authentication for each method.

### Configure from Terracreds
> New in version `2.1.0`

You can create and view the configuration for any vault provider by running `terracreds config` and then using the subcommand for the specific vault provider. The commands to generate the config from `terracreds` will be shown for each provider listed below.

### AWS Secrets Manager
In order to leverage `terracreds` to manage secrets in `AWS Secrets Manager` the following block needs to be provided in the configuration file:
```yaml
aws:
  description: my_terraform_api_token
  region: us-west-2
  secretName: my-secret-name
```

This can be generated via `terracreds` by running:
```bash
terracreds config aws --description 'my super secret' --region 'us-west-2' --secret-name 'my-secret-name'
```

| Value | Description | Required |
| ----- | ----------- | -------- |
| `description` | A brief description to provide for the secret object viewable in `Secrets Manager` | `yes` |
| `region` | The `Secrets Manager` instance's region where the secret will be stored | `yes` | 
| `secretName` | A name for the secret. If omitted and using `terraform login` the hostname of the TACOS server will be used for the name instead | `no` |

The following role permissions are required in order for the `EC2 Instance Role` to levearge `terracreds` with `AWS Secrets Manager`:
```hcl
Action = [
  "secretsmanager:CreateSecret",
  "secretsmanager:DeleteSecret",
  "secretsmanager:GetSecretValue",
  "secretsmanager:PutSecretValue"
]
```
### Azure Key Vault
In order to leverage `terracreds` to manage secrets in `Azure Key Vault` the following block needs to be provided in the configuration file:
```yaml
azure:
  secretName: my-secret-name
  subscriptionId: 5df41dfe-4310-46e5-800a-c5bc71ac7ac0
  vaultUri: https://mykeyvault.vault.azure.net
```

The configuration can be generated via `terracreds` by running:
```bash
terracreds config azure --subscription-id '5df41dfe-4310-46e5-800a-c5bc71ac7ac0' --vault-uri 'https://mykeyvault.vault.azure.net' --secret-name 'my-secret-name'
```

| Value | Description | Required |
| ----- | ----------- | -------- |
| `secretName` | A name for the secret. If omitted and using `terraform login` the hostname of the TACOS server will be used for the name instead | `no` |
| `subscriptionId` | The Azure subscription ID where the `Azure Key Vault` has been created | `yes` | 
| `vaultUri` | The URI for the `Azure Key Vault` where you want to store or retrieve your credentials | `yes` |

The following `Azure Key Vault Access Policies` are required to be given to the `Managed Service Identity` for it to leverage `terracreds`:
```hcl
secret_permissions = [
  "Get",
  "List",
  "Set",
  "Delete"
]
```
> Since `Azure Key Vault` doesn't support the period character in a secret name a helper function will replace any periods with dashes so they can be successfully stored. This means a `terraform` API token name that would usually be `app.terraform.io` will become `app-terraform-io`

### Google Secret Manager
> New in version `2.1.0`

In order to leverage `terracreds` to manage secrets in `Google Secret Manager` the following block needs to be provided in the configuration file:
```yaml
gcp:
  projectId: my-gcp-project
  secretId: my-secret-name
```

The configuration can be generated via `terracreds` by running:
```bash
terracreds config gcp --project-id 'my-gcp-project' --secret-id 'my-secret-name'
```

| Value | Description | Required |
| ----- | ----------- | -------- |
| `projectId` | The name of the `GCP` project ID where the `Secret Manager` API has been enabled | `yes` |
| `secretId` | The name of the secret ID | `no` |

The `Google IAM` role `secretmanager.admin` is suggested in order to fully manage the secrets with `terracreds`

### HashiCorp Vault
In order to leverage `terracreds` to manage secrets in `HashiCorp Vault` the following block needs to be provided in the configuration file:
```yaml
hcvault:
  environmentTokenName: HASHI_TOKEN
  keyVaultPath: kv
  secretName: my-secret-name
  secretPath: tfe
  vaultUri: http://localhost:8200
```

The configuration can be generated via `terracreds` by running:
```bash
terracreds config hashicorp \ 
  --environment-token-name 'HASHI_TOKEN' \
  --key-vault-path 'kv' \
  --secret-name 'my-secret-name' \
  --secret-path 'tfe' \
  --vault-uri 'http://localhost:8200"
```

| Value | Description | Required |
| ----- | ----------- | -------- |
| `environmentTokenName` | The name of the environment variable that contains the token value to authenticate with `HashiCorp Vault` | `yes` |
| `keyVaultPath` | The path to the `Key Vault` object within the vault | `yes` |
| `secretName` | A name for the secret. If omitted and using `terraform login` the hostname of the TACOS server will be used for the name instead | `no` |
| `secretPath` | The path of the secret within `HashiCorp Vault` | `yes` |
| `vaultUri` | The URI for the `HashiCorp Vault` instance | `yes` |

## Protection
In order to add some protection `terracreds` adds a username to the credential object stored in the local operating system, and checks to ensure that the user requesting access to the secret is the same user as the secret's creator.  

Any attempt to access or modify this secret from `terracreds` outside of the user that created the credentail will lead to denial messages. Additionally, if the credential name is not found, the same access denied message will be provided in lieu of a generic not found message to help prevent brute force attempts

## Logging
> New in version `2.1.0`

By default `terracreds` will generate a configuration file in the same location where the `terracreds` binary was first run. This can now be overridden by setting an environment variable that sets the path to the desired location of the configuration file:

For Linux/macOS:
```bash
export TC_CONFIG_PATH=/home/username/
```

For Windows:
```powershell
$env:TC_CONFIG_PATH="C:\Temp"
```

To persist this change you can set this variables either in `.bashrc` for Linux/macOS or setup a PowerShell profile for Windows.

To enable logging for Windows setup the `config.yaml` as follows:
```yaml
logging:
  enabled: true
  path: C:\Temp
```

To enable logging for macOS and Linux to a directory called `.terracreds` in the user's home profile:
```yaml
logging:
  enabled: true
  path: ~/.terracreds
```

You can also use `terracreds` to configure logging:
```bash
terracreds config logging --path '~/.terracreds' --enabled
```

The log is helpful in understanding if an object was found, deleted, updated or added, and will be found at the path defined in the configuration file as `terracreds.log`.

In addition all error messages returned by the underlying libraries will be logged when logging is enabled and an error is encountered.

## Troubleshooting

### Known Issues
When you enable `terracreds` as a credential helper Terraform will begin using it for all authentication regardless of the destination server. This means that when you try to install/download providers or modules from the public Terraform registry `https://registry.terraform.io/`, or any other public registry, Terraform will try to authenticate against the server using `terracreds`. If there's no credential in the vault found for that server it will error out.

To work around this issue you'll need to set a dummy value for any public registries. Run this command for each public repo that Terraform will need to access. In this example we're using `registry.terraform.io` so be sure to replace it with the correct server value if the one you require is different:
```bash
terracreds create -n registry.terraform.io -v dummy_token
```

### Linux
If you are having trouble viewing, deleting, or saving credentials on Linux systems using `gnome-keyring` you must ensure that you have unlocked the collection using `gnome-keyring-daemon --unlock` otherwise you will see the following error message in the logs:

```txt
ERROR: <TIMESTAMP> - failed to unlock correct collection '/org/freedesktop/secrets/collection/login'
```

If the daemon has unlocked the collection but you're still getting prompted for credentials check to make sure that only a single instance of the daemon is running:

```bash
ps -ef | grep gnome-keyring
```

If more than one daemon is running, take note of the pid, and use `kill` to terminate the additional daemon. Try your previous command again
and it should now be working.
