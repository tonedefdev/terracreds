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

	api "github.com/tonedefdev/terracreds/api"
	helpers "github.com/tonedefdev/terracreds/pkg/helpers"
)

type Windows struct {
	ApiToken api.CredentialResponse
	Command  string
	Config   api.Config
	Context  *cli.Context
	Hostname string
	Token    interface{}
	User     *user.User
}

func (w Windows) Create() {
	var method string
	_, err := wincred.GetGenericCredential(w.Hostname)
	if err != nil {
		method = "Created"
	} else {
		method = "Updated"
	}

	cred := wincred.NewGenericCredential(w.Hostname)
	if w.Token == nil {
		err = json.NewDecoder(os.Stdin).Decode(&w.ApiToken)
		if err != nil {
			fmt.Print(err.Error())
		}
		cred.CredentialBlob = []byte(w.ApiToken.Token)
	} else {
		str := fmt.Sprintf("%v", w.Token)
		cred.CredentialBlob = []byte(str)
	}

	cred.UserName = string(w.User.Username)
	err = cred.Write()

	if err == nil {
		msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), w.Hostname)
		helpers.Logging(w.Config, msg, "SUCCESS")

		if w.Token != nil {
			fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, w.Hostname)
		}
	} else {
		helpers.Logging(w.Config, fmt.Sprintf("- %s", err), "ERROR")

		if w.Token != nil {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

func (w Windows) Delete() {
	cred, err := wincred.GetGenericCredential(w.Hostname)
	if err == nil && cred.UserName == w.User.Username {
		cred.Delete()

		msg := fmt.Sprintf("- the credential object '%s' has been removed", w.Hostname)
		helpers.Logging(w.Config, msg, "INFO")

		if w.Command == "delete" {
			msg := fmt.Sprintf("The credential object '%s' has been removed", w.Hostname)
			fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
		}
	} else {
		helpers.Logging(w.Config, fmt.Sprintf("- %s", err), "ERROR")

		if w.Command == "delete" {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

func (w Windows) Get() {
	if w.Config.Logging.Enabled == true {
		msg := fmt.Sprintf("- terraform server: %s", w.Hostname)
		helpers.Logging(w.Config, msg, "INFO")
		msg = fmt.Sprintf("- user requesting access: %s", string(w.User.Username))
		helpers.Logging(w.Config, msg, "INFO")
	}

	cred, err := wincred.GetGenericCredential(w.Hostname)
	if err == nil && cred.UserName == w.User.Username {
		response := &api.CredentialResponse{
			Token: string(cred.CredentialBlob),
		}
		responseA, _ := json.Marshal(response)
		fmt.Println(string(responseA))

		if w.Config.Logging.Enabled == true {
			msg := fmt.Sprintf("- token was retrieved for: %s", w.Hostname)
			helpers.Logging(w.Config, msg, "INFO")
		}
	} else {
		if w.Config.Logging.Enabled == true {
			helpers.Logging(w.Config, fmt.Sprintf("- %s", err), "ERROR")
		}
		fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
	}
}
