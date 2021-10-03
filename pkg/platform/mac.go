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

	api "github.com/tonedefdev/terracreds/api"
	helpers "github.com/tonedefdev/terracreds/pkg/helpers"
)

type Mac struct {
	ApiToken api.CredentialResponse
	Command  string
	Config   api.Config
	Context  *cli.Context
	Hostname string
	Token    interface{}
	User     *user.User
}

func (m Mac) Create() {
	var method string
	_, err := keyring.Get(m.Hostname, string(m.User.Username))
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
		err = keyring.Set(m.Hostname, string(m.User.Username), m.ApiToken.Token)
	} else {
		str := fmt.Sprintf("%v", m.Token)
		err = keyring.Set(m.Hostname, string(m.User.Username), str)
	}

	if err == nil {
		msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), m.Hostname)
		helpers.Logging(m.Config, msg, "SUCCESS")

		if m.Token != nil {
			fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, m.Hostname)
		}
	} else {
		helpers.Logging(m.Config, fmt.Sprintf("- %s", err), "ERROR")

		if m.Token != nil {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

func (m Mac) Delete() {
	err := keyring.Delete(m.Hostname, string(m.User.Username))
	if err == nil {
		msg := fmt.Sprintf("- the credential object '%s' has been removed", m.Hostname)
		helpers.Logging(m.Config, msg, "INFO")

		if m.Command == "delete" {
			msg := fmt.Sprintf("The credential object '%s' has been removed", m.Hostname)
			fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
		}
	} else {
		helpers.Logging(m.Config, fmt.Sprintf("- %s", err), "ERROR")

		if m.Command == "delete" {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

func (m Mac) Get() {
	if m.Config.Logging.Enabled == true {
		msg := fmt.Sprintf("- terraform server: %s", m.Hostname)
		helpers.Logging(m.Config, msg, "INFO")
		msg = fmt.Sprintf("- user requesting access: %s", string(m.User.Username))
		helpers.Logging(m.Config, msg, "INFO")
	}

	secret, err := keyring.Get(m.Hostname, string(m.User.Username))
	if err != nil {
		if m.Config.Logging.Enabled == true {
			helpers.Logging(m.Config, fmt.Sprintf("- %s", err), "ERROR")
		}
		fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
	} else {
		response := &api.CredentialResponse{
			Token: secret,
		}
		responseA, _ := json.Marshal(response)
		fmt.Println(string(responseA))

		if m.Config.Logging.Enabled == true {
			msg := fmt.Sprintf("- token was retrieved for: %s", m.Hostname)
			helpers.Logging(m.Config, msg, "INFO")
		}
	}
}
