# Terracreds
Is a credential helper for Terraform Cloud/Enterprise that allows storing your API token credentials inside of the operating system's vault instead of in a plain text configuration file. We all know that storing sensitive secrets in plain text can pose security threats, and Terraform doesn't come pre-packaged with a credential helper so we decided to create one that leverages your internal operating system's credential manager to store your secrets instead. Terracreds allows for tighter security controls and better management of your Terraform API tokens.

#### Currently supported Operating Systems:
- [x] Windows (Credential Manager)
- [ ] MacOS (Keychain)
- [ ] Linux (ksecretservice or gnome-keyring)

## Windows Install via Chocolatey
The fastest way to install `terracreds` on Windows is via our Chocolatey package:
```powershell
choco install terracreds -y
```

Once installed run the following command to verify `terracreds` was installed properly:
```powerhsell
terracreds -v
```

## Manual Install
To install from source make sure you have the latest version of Go installed, and run the following command:
```go
go get github.com/tonedefev/terracreds 
```