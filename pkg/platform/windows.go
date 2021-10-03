package platform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/danieljoos/wincred"
	"github.com/fatih/color"

	api "github.com/tonedefdev/terracreds/api"
	helpers "github.com/tonedefdev/terracreds/pkg/helpers"
)

type Windows struct{}

// Create stores or updates a credential in the Windows Credential Manager
func (w Windows) Create(cfg api.Config, hostname string, token interface{}, user *user.User) {
	var method string
	_, err := wincred.GetGenericCredential(hostname)
	if err != nil {
		method = "Created"
	} else {
		method = "Updated"
	}

	cred := wincred.NewGenericCredential(hostname)
	if token == nil {
		err = json.NewDecoder(os.Stdin).Decode(&api.CredentialResponse{})
		if err != nil {
			fmt.Print(err.Error())
		}
		cred.CredentialBlob = []byte(api.CredentialResponse{}.Token)
	} else {
		str := fmt.Sprintf("%v", token)
		cred.CredentialBlob = []byte(str)
	}

	cred.UserName = string(user.Username)
	err = cred.Write()

	if err == nil {
		msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), hostname)
		helpers.Logging(cfg, msg, "SUCCESS")

		if token != nil {
			fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, hostname)
		}
	} else {
		helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")

		if token != nil {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

// Delete removes or forgets a Terraform API token from the Windows Credential Manager
func (w Windows) Delete(cfg api.Config, command string, hostname string, user *user.User) {
	cred, err := wincred.GetGenericCredential(hostname)
	if err == nil && cred.UserName == user.Username {
		cred.Delete()

		msg := fmt.Sprintf("- the credential object '%s' has been removed", hostname)
		helpers.Logging(cfg, msg, "INFO")

		if command == "delete" {
			msg := fmt.Sprintf("The credential object '%s' has been removed", hostname)
			fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
		}
	} else {
		helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")

		if command == "delete" {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

// Get retrieves a Terraform API token stored in Windows Credential Manager
func (w Windows) Get(cfg api.Config, hostname string, user *user.User) {
	if cfg.Logging.Enabled == true {
		msg := fmt.Sprintf("- terraform server: %s", hostname)
		helpers.Logging(cfg, msg, "INFO")
		msg = fmt.Sprintf("- user requesting access: %s", string(user.Username))
		helpers.Logging(cfg, msg, "INFO")
	}

	cred, err := wincred.GetGenericCredential(hostname)
	if err == nil && cred.UserName == user.Username {
		response := &api.CredentialResponse{
			Token: string(cred.CredentialBlob),
		}
		responseA, _ := json.Marshal(response)
		fmt.Println(string(responseA))

		if cfg.Logging.Enabled == true {
			msg := fmt.Sprintf("- token was retrieved for: %s", hostname)
			helpers.Logging(cfg, msg, "INFO")
		}
	} else {
		if cfg.Logging.Enabled == true {
			helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
		}
		fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
	}
}
