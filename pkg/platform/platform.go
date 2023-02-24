package platform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"

	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/errors"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/tonedefdev/terracreds/pkg/vault"
)

type Platform struct{}

// Creates stores or updates a secret in in a vault
func (platform *Platform) Create(cfg *api.Config, hostname string, token any, user *user.User, vault vault.TerraVault) error {
	var method string
	method = "Updated"

	if token == nil {
		var response api.CredentialResponse
		err := json.NewDecoder(os.Stdin).Decode(&response)
		if err != nil {
			helpers.CheckError(err)
		}

		token = response.Token
	}

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

// Deletes removes or forgets a secret in a vault
func (platform *Platform) Delete(cfg *api.Config, command string, hostname string, user *user.User, vault vault.TerraVault) error {
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
		err = &errors.CustomError{
			Message: "You do not have permission to modify this credential",
			Level:   "ERROR",
		}

		return err
	}

	return nil
}

// Get retrieves a secret from a vault
func (platform *Platform) Get(cfg *api.Config, hostname string, user *user.User, vault vault.TerraVault) ([]byte, error) {
	if cfg.Logging.Enabled {
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

		if cfg.Logging.Enabled && err == nil {
			msg := fmt.Sprintf("- token was retrieved for: %s", hostname)
			helpers.Logging(cfg, msg, "INFO")
		}

		return token, err
	}

	if cfg.Logging.Enabled {
		helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
	}

	err = &errors.CustomError{
		Message: "You do not have permission to view this credential",
		Level:   "ERROR",
	}

	return nil, err
}

// List returns a list of secrets from a vault in a specified format
func (platform *Platform) List(c *cli.Context, cfg *api.Config, secretNames []string, user *user.User, vault vault.TerraVault) ([]string, error) {
	var secretValues []string
	if cfg.Logging.Enabled {
		msg := fmt.Sprintf("- user requesting access: %s", string(user.Username))
		helpers.Logging(cfg, msg, "INFO")
	}

	if vault != nil {
		secrets, err := vault.List(secretNames)
		if err != nil {
			return nil, err
		}

		return secrets, nil
	}

	for _, secret := range secretNames {
		if cfg.Logging.Enabled {
			msg := fmt.Sprintf("- secret name requested: %s", secret)
			helpers.Logging(cfg, msg, "INFO")
		}

		cred, err := keyring.Get(secret, string(user.Username))
		if err != nil {
			return nil, err
		}

		value := string(cred)
		secretValues = append(secretValues, value)
	}

	return secretValues, nil
}
