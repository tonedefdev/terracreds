package platform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/fatih/color"
	"github.com/zalando/go-keyring"

	api "github.com/tonedefdev/terracreds/api"
	helpers "github.com/tonedefdev/terracreds/pkg/helpers"
)

type Mac struct {
	ApiToken api.CredentialResponse
	Token    interface{}
}

func (m Mac) Create(cfg api.Config, hostname string, user *user.User) {
	var method string
	_, err := keyring.Get(hostname, string(user.Username))
	if err != nil {
		method = "Created"
	} else {
		method = "Updated"
	}

	if m.Token == nil {
		err = json.NewDecoder(os.Stdin).Decode(&m.ApiToken)
		if err != nil {
			fmt.Print(err.Error())
		}
		err = keyring.Set(hostname, string(user.Username), m.ApiToken.Token)
	} else {
		str := fmt.Sprintf("%v", m.Token)
		err = keyring.Set(hostname, string(user.Username), str)
	}

	if err == nil {
		msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), hostname)
		helpers.Logging(cfg, msg, "SUCCESS")

		if m.Token != nil {
			fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, hostname)
		}
	} else {
		helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")

		if m.Token != nil {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

func (m Mac) Delete(cfg api.Config, command string, hostname string, user *user.User) {
	err := keyring.Delete(hostname, string(user.Username))
	if err == nil {
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

func (m Mac) Get(cfg api.Config, hostname string, user *user.User) {
	if cfg.Logging.Enabled == true {
		msg := fmt.Sprintf("- terraform server: %s", hostname)
		helpers.Logging(cfg, msg, "INFO")
		msg = fmt.Sprintf("- user requesting access: %s", string(user.Username))
		helpers.Logging(cfg, msg, "INFO")
	}

	secret, err := keyring.Get(hostname, string(user.Username))
	if err != nil {
		if cfg.Logging.Enabled == true {
			helpers.Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
		}
		fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
	} else {
		response := &api.CredentialResponse{
			Token: secret,
		}
		responseA, _ := json.Marshal(response)
		fmt.Println(string(responseA))

		if cfg.Logging.Enabled == true {
			msg := fmt.Sprintf("- token was retrieved for: %s", hostname)
			helpers.Logging(cfg, msg, "INFO")
		}
	}
}
