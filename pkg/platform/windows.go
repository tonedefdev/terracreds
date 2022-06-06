package platform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/danieljoos/wincred"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/tonedefdev/terracreds/pkg/vault"
)

type Windows struct{}

// Create stores or updates a Terafform API token in Windows Credential Manager or a specified Cloud Vault
func (w *Windows) Create(cfg api.Config, hostname string, token interface{}, user *user.User, vault vault.TerraVault) error {
	var method string
	method = "Updated"

	if vault != nil {
		_, err := vault.Get()
		if err != nil {
			method = "Created"
		}

		if token == nil {
			var response api.CredentialResponse
			err = json.NewDecoder(os.Stdin).Decode(&response)
			if err != nil {
				helpers.CheckError(err)
			}

			token = response.Token
		}

		secretValue := fmt.Sprintf("%v", token)
		err = vault.Create(secretValue, method)
		if err != nil {
			helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
			return err
		}

		fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, hostname)
		return err
	}

	_, err := wincred.GetGenericCredential(hostname)
	if err != nil {
		method = "Created"
	}

	cred := wincred.NewGenericCredential(hostname)
	if token == nil {
		var response api.CredentialResponse
		err = json.NewDecoder(os.Stdin).Decode(&response)
		if err != nil {
			helpers.CheckError(err)
		}

		cred.CredentialBlob = []byte(response.Token)
		cred.UserName = string(user.Username)
		err = cred.Write()
		return err
	}

	str := fmt.Sprintf("%v", token)
	cred.CredentialBlob = []byte(str)
	cred.UserName = string(user.Username)
	err = cred.Write()

	if err != nil && token != nil {
		fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		return err
	}

	if err != nil {
		helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
		return err
	}

	msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), hostname)
	helpers.Logging(cfg, msg, "SUCCESS")

	if token != nil {
		fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, hostname)
		return err
	}

	return err
}

// Delete removes or forgets a Terraform API token from the Windows Credential Manager
func (w *Windows) Delete(cfg api.Config, command string, hostname string, user *user.User, vault vault.TerraVault) error {
	if vault != nil {
		err := vault.Delete()
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("- the credential object '%s' has been removed", hostname)
		helpers.Logging(cfg, msg, "INFO")

		if command == "delete" {
			msg := fmt.Sprintf("The credential object '%s' has been removed", hostname)
			fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
		}

		return err
	}

	cred, err := wincred.GetGenericCredential(hostname)
	if err == nil && cred.UserName == user.Username {
		cred.Delete()

		msg := fmt.Sprintf("- the credential object '%s' has been removed", hostname)
		helpers.Logging(cfg, msg, "INFO")

		if command == "delete" {
			msg := fmt.Sprintf("The credential object '%s' has been removed", hostname)
			fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
		}

		return err
	}

	helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
	if command == "delete" {
		fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
	}

	return nil
}

// Get retrieves a Terraform API token in Windows Credential Manager
func (w *Windows) Get(cfg api.Config, hostname string, user *user.User, vault vault.TerraVault) ([]byte, error) {
	if cfg.Logging.Enabled == true {
		msg := fmt.Sprintf("- terraform server: %s", hostname)
		helpers.Logging(cfg, msg, "INFO")
		msg = fmt.Sprintf("- user requesting access: %s", string(user.Username))
		helpers.Logging(cfg, msg, "INFO")
	}

	if vault != nil {
		token, err := vault.Get()
		if err != nil {
			helpers.CheckError(err)
		}

		response := &api.CredentialResponse{
			Token: string(token),
		}

		token, err = json.Marshal(response)
		return token, err
	}

	cred, err := wincred.GetGenericCredential(hostname)
	if err == nil && cred.UserName == user.Username {
		response := &api.CredentialResponse{
			Token: string(cred.CredentialBlob),
		}

		token, err := json.Marshal(response)

		if cfg.Logging.Enabled == true {
			msg := fmt.Sprintf("- token was retrieved for: %s", hostname)
			helpers.Logging(cfg, msg, "INFO")
		}

		return token, err
	}

	if cfg.Logging.Enabled == true {
		helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
	}

	fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
	return nil, err
}

func (w *Windows) List(c *cli.Context, cfg api.Config, secretNames []string, user *user.User, vault vault.TerraVault) ([]string, error) {
	var secretValues []string

	if vault != nil {
		secrets, err := vault.List(secretNames)
		if err != nil {
			return nil, err
		}

		return secrets, nil
	}

	for _, secret := range secretNames {
		cred, err := wincred.GetGenericCredential(secret)
		if err == nil && cred.UserName == user.Username {
			value := string(cred.CredentialBlob)
			secretValues = append(secretValues, value)
		} else {
			return nil, err
		}
	}

	return secretValues, nil
}
