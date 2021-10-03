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

type Linux struct {
	ApiToken api.CredentialResponse
	Command  string
	Config   api.Config
	Context  *cli.Context
	Hostname string
	Token    interface{}
	User     *user.User
}

func (l Linux) Create() {
	var method string
	_, err := keyring.Get(l.Hostname, string(l.User.Username))
	if err != nil {
		method = "Created"
	} else {
		method = "Updated"
	}

	if l.Token == nil {
		err = json.NewDecoder(os.Stdin).Decode(&l.ApiToken)
		if err != nil {
			fmt.Print(err.Error())
		}
		err = keyring.Set(l.Hostname, string(l.User.Username), l.ApiToken.Token)
	} else {
		str := fmt.Sprintf("%v", l.Token)
		err = keyring.Set(l.Hostname, string(l.User.Username), str)
	}

	if err == nil {
		msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), l.Hostname)
		helpers.Logging(l.Config, msg, "SUCCESS")

		if l.Token != nil {
			fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, l.Hostname)
		}
	} else {
		helpers.Logging(l.Config, fmt.Sprintf("- %s", err), "ERROR")

		if l.Token != nil {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

func (l Linux) Delete() {
	err := keyring.Delete(l.Hostname, string(l.User.Username))
	if err == nil {
		msg := fmt.Sprintf("- the credential object '%s' has been removed", l.Hostname)
		helpers.Logging(l.Config, msg, "INFO")

		if l.Command == "delete" {
			msg := fmt.Sprintf("The credential object '%s' has been removed", l.Hostname)
			fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
		}
	} else {
		helpers.Logging(l.Config, fmt.Sprintf("- %s", err), "ERROR")

		if l.Command == "delete" {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

func (l Linux) Get() {
	if l.Config.Logging.Enabled == true {
		msg := fmt.Sprintf("- terraform server: %s", l.Hostname)
		helpers.Logging(l.Config, msg, "INFO")
		msg = fmt.Sprintf("- user requesting access: %s", string(l.User.Username))
		helpers.Logging(l.Config, msg, "INFO")
	}

	secret, err := keyring.Get(l.Hostname, string(l.User.Username))
	if err != nil {
		if l.Config.Logging.Enabled == true {
			helpers.Logging(l.Config, fmt.Sprintf("- %s", err), "ERROR")
		}
		fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
	} else {
		response := &api.CredentialResponse{
			Token: secret,
		}
		responseA, _ := json.Marshal(response)
		fmt.Println(string(responseA))

		if l.Config.Logging.Enabled == true {
			msg := fmt.Sprintf("- token was retrieved for: %s", l.Hostname)
			helpers.Logging(l.Config, msg, "INFO")
		}
	}
}
