package platform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/fatih/color"
	"github.com/zalando/go-keyring"

	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/tonedefdev/terracreds/pkg/vault"
)

type Linux struct{}

// Create stores or updates a Terraform API token in Gnome Keyring or an external vault provider
func (l *Linux) Create(cfg api.Config, hostname string, token interface{}, user *user.User, vault vault.TerraVault) error {
	var method string
	method = "Updated"

	if vault != nil {
		_, err := vault.Get()
		if err != nil {
			method = "Created"
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

	_, err := keyring.Get(hostname, string(user.Username))
	if err != nil {
		method = "Created"
	}

	if token == nil {
		err = json.NewDecoder(os.Stdin).Decode(&api.CredentialResponse{})
		if err != nil {
			helpers.CheckError(err)
		}
		err = keyring.Set(hostname, string(user.Username), api.CredentialResponse{}.Token)
		return err
	}

	str := fmt.Sprintf("%v", token)
	err = keyring.Set(hostname, string(user.Username), str)

	if err != nil && token != nil {
		fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		return err
	}

	if err != nil {
		helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
		return nil
	}

	msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), hostname)
	helpers.Logging(cfg, msg, "SUCCESS")

	if token != nil {
		fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, hostname)
		return err
	}

	return err
}

// Delete removes or forgets a Terraform API token in Gnome Keyring or an external vault provider
func (l *Linux) Delete(cfg api.Config, command string, hostname string, user *user.User, vault vault.TerraVault) error {
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

	err := keyring.Delete(hostname, string(user.Username))
	if err == nil {
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

// Get retrieves a Terraform API token in Gnome Keyring or an external vault provider
func (l *Linux) Get(cfg api.Config, hostname string, user *user.User, vault vault.TerraVault) ([]byte, error) {
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

	secret, err := keyring.Get(hostname, string(user.Username))
	if err == nil {
		response := &api.CredentialResponse{
			Token: secret,
		}
		token, err := json.Marshal(response)

		if cfg.Logging.Enabled == true && err == nil {
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
