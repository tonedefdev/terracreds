package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/vault"
	"github.com/urfave/cli/v2"
)

type recordingTerraCreds struct {
	created []string
	deleted []string
	got     []string
	listed  []string
}

func (mock *recordingTerraCreds) Create(_ *api.Config, hostname string, token any, _ *user.User, _ vault.TerraVault) error {
	mock.created = append(mock.created, hostname+":"+fmt.Sprint(token))
	return nil
}

func (mock *recordingTerraCreds) Delete(_ *api.Config, command string, hostname string, _ *user.User, _ vault.TerraVault) error {
	mock.deleted = append(mock.deleted, command+":"+hostname)
	return nil
}

func (mock *recordingTerraCreds) Get(_ *api.Config, hostname string, _ *user.User, _ vault.TerraVault) ([]byte, error) {
	mock.got = append(mock.got, hostname)
	return []byte(`{"token":"value"}`), nil
}

func (mock *recordingTerraCreds) List(_ *cli.Context, _ *api.Config, names []string, _ *user.User, _ vault.TerraVault) ([]string, error) {
	mock.listed = append([]string(nil), names...)
	values := make([]string, len(names))
	for i, name := range names {
		values[i] = "value-" + name
	}
	return values, nil
}

func testConfig(t *testing.T) *Config {
	t.Helper()
	return &Config{
		Cfg:        &api.Config{},
		ConfigFile: ConfigFile{Path: filepath.Join(t.TempDir(), "config.yaml")},
		TerraCreds: &recordingTerraCreds{},
	}
}

func commandContext(t *testing.T, args ...string) *cli.Context {
	t.Helper()
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	for _, name := range []string{"name", "secret", "secret-names", "override-replace-string", "description", "region", "secret-name", "subscription-id", "vault-uri", "project-id", "secret-id", "environment-token-name", "key-vault-path", "secret-path", "path", "secret-list"} {
		set.String(name, "", "")
	}
	for _, name := range []string{"as-json", "as-tfvars", "from-config", "enabled", "force", "use-local-vault-only"} {
		set.Bool(name, false, "")
	}
	if err := set.Parse(args); err != nil {
		t.Fatal(err)
	}
	return cli.NewContext(cli.NewApp(), set, nil)
}

func setArgs(t *testing.T, args ...string) {
	t.Helper()
	original := os.Args
	os.Args = append([]string(nil), args...)
	t.Cleanup(func() { os.Args = original })
}

func readConfig(t *testing.T, cmd *Config) api.Config {
	t.Helper()
	if err := cmd.LoadConfig(cmd.ConfigFile.Path); err != nil {
		t.Fatal(err)
	}
	return *cmd.Cfg
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	fn()
	writer.Close()
	os.Stdout = original
	var output bytes.Buffer
	_, _ = io.Copy(&output, reader)
	reader.Close()
	return output.String()
}

func TestGetSecretNameAndNewTerraVault(t *testing.T) {
	cases := []struct {
		name string
		cfg  api.Config
		want string
	}{
		{"aws", api.Config{Aws: api.Aws{SecretName: "aws-secret"}}, "aws-secret"},
		{"azure", api.Config{Azure: api.Azure{SecretName: "azure-secret"}}, "azure-secret"},
		{"hashicorp", api.Config{HashiVault: api.HCVault{SecretName: "vault-secret"}}, "vault-secret"},
		{"hostname", api.Config{}, "host.example"},
	}
	for _, test := range cases {
		if got := GetSecretName(&test.cfg, "host.example"); got != test.want {
			t.Errorf("%s: got %q, want %q", test.name, got, test.want)
		}
	}

	providers := []struct {
		name string
		cfg  api.Config
		want any
	}{
		{"aws", api.Config{Aws: api.Aws{Region: "us-east-1", Description: "desc"}}, &vault.AwsSecretsManager{}},
		{"azure", api.Config{Azure: api.Azure{VaultUri: "https://vault", SubscriptionId: "sub"}}, &vault.AzureKeyVault{}},
		{"gcp", api.Config{GCP: api.GCP{ProjectId: "project"}}, &vault.GCPSecretManager{}},
		{"hashicorp", api.Config{HashiVault: api.HCVault{VaultUri: "https://vault"}}, &vault.HashiVault{}},
	}
	for _, test := range providers {
		cmd := Config{Cfg: &test.cfg}
		if got := reflect.TypeOf(cmd.NewTerraVault("host")); got != reflect.TypeOf(test.want) {
			t.Errorf("%s: got %v, want %v", test.name, got, reflect.TypeOf(test.want))
		}
	}
	emptyConfig := Config{Cfg: &api.Config{}}
	if got := emptyConfig.NewTerraVault("host"); got != nil {
		t.Fatalf("empty config returned %T vault", got)
	}
}

func TestConfigActionsWriteExpectedValues(t *testing.T) {
	tests := []struct {
		name string
		run  func(*Config) error
		want func(api.Config) bool
	}{
		{"aws", func(cmd *Config) error {
			return cmd.newCommandActionAws(commandContext(t, "--description=d", "--region=r", "--secret-name=s"))
		}, func(cfg api.Config) bool {
			return cfg.Aws.Description == "d" && cfg.Aws.Region == "r" && cfg.Aws.SecretName == "s"
		}},
		{"azure", func(cmd *Config) error {
			return cmd.newCommandActionAzure(commandContext(t, "--subscription-id=i", "--vault-uri=u", "--secret-name=s"))
		}, func(cfg api.Config) bool {
			return cfg.Azure.SubscriptionId == "i" && cfg.Azure.VaultUri == "u" && cfg.Azure.SecretName == "s"
		}},
		{"gcp", func(cmd *Config) error {
			return cmd.newCommandActionGcp(commandContext(t, "--project-id=p", "--secret-id=s"))
		}, func(cfg api.Config) bool { return cfg.GCP.ProjectId == "p" && cfg.GCP.SecretId == "s" }},
		{"hashicorp", func(cmd *Config) error {
			return cmd.newCommandActionHashi(commandContext(t, "--environment-token-name=e", "--key-vault-path=k", "--secret-name=s", "--secret-path=p", "--vault-uri=u"))
		}, func(cfg api.Config) bool {
			return cfg.HashiVault.EnvironmentTokenName == "e" && cfg.HashiVault.KeyVaultPath == "k" && cfg.HashiVault.SecretName == "s" && cfg.HashiVault.SecretPath == "p" && cfg.HashiVault.VaultUri == "u"
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := testConfig(t)
			if err := test.run(cmd); err != nil {
				t.Fatal(err)
			}
			if !test.want(readConfig(t, cmd)) {
				t.Fatalf("config action wrote unexpected values: %#v", readConfig(t, cmd))
			}
		})
	}
}

func TestConfigLoggingSecretsViewAndLoadErrors(t *testing.T) {
	cmd := testConfig(t)
	cmd.Cfg.Logging = api.Logging{Enabled: false}
	if err := cmd.newCommandActionLogging(commandContext(t, "--enabled", "--path=/tmp")); err != nil {
		t.Fatal(err)
	}
	if err := cmd.newCommandActionSecrets(commandContext(t, "--secret-list=one,two")); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cmd.Cfg.Secrets, []string{"one", "two"}) {
		t.Fatalf("secrets = %#v", cmd.Cfg.Secrets)
	}
	if err := cmd.newCommandActionView(commandContext(t)); err != nil {
		t.Fatal(err)
	}
	if err := cmd.LoadConfig(filepath.Join(t.TempDir(), "missing.yaml")); err == nil {
		t.Fatal("LoadConfig returned nil for missing file")
	}
}

func TestCommandActionsUseRecordingMock(t *testing.T) {
	cmd := testConfig(t)
	mock := cmd.TerraCreds.(*recordingTerraCreds)
	setArgs(t, "terracreds", "create", "--name=host", "--secret=value")
	if err := cmd.newCommandActionCreate(commandContext(t, "--name=host", "--secret=value")); err != nil {
		t.Fatal(err)
	}
	setArgs(t, "terracreds", "store", "host")
	if err := cmd.newCommandActionStore(commandContext(t)); err != nil {
		t.Fatal(err)
	}
	setArgs(t, "terracreds", "forget", "host")
	if err := cmd.newCommandActionForget(commandContext(t)); err != nil {
		t.Fatal(err)
	}
	setArgs(t, "terracreds", "delete", "--name=host")
	if err := cmd.newCommandActionDelete(commandContext(t, "--name=host")); err != nil {
		t.Fatal(err)
	}
	setArgs(t, "terracreds", "get", "host")
	if output := captureStdout(t, func() {
		if err := cmd.newCommandActionGet(); err != nil {
			t.Fatal(err)
		}
	}); !strings.Contains(output, `{"token":"value"}`) {
		t.Fatalf("get output = %q", output)
	}
	if !reflect.DeepEqual(mock.created, []string{"host:value", "host:<nil>"}) {
		t.Fatalf("created calls = %#v", mock.created)
	}
	if !reflect.DeepEqual(mock.deleted, []string{"delete:host", "delete:host"}) {
		t.Fatalf("deleted calls = %#v", mock.deleted)
	}
	if !reflect.DeepEqual(mock.got, []string{"host"}) {
		t.Fatalf("get calls = %#v", mock.got)
	}
}

func TestListOutputModesAndConfiguredNames(t *testing.T) {
	for _, test := range []struct {
		name string
		args []string
		want string
	}{
		{"json", []string{"--secret-names=one,two", "--as-json"}, "{\"one\":\"value-one\",\"two\":\"value-two\"}\n"},
		{"tfvars", []string{"--secret-names=one-name,two", "--as-tfvars", "--override-replace-string=_"}, "TF_VAR_one_name=value-one-name\nTF_VAR_two=value-two\n"},
		{"plain", []string{"--from-config"}, "value-one\nvalue-two\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			cmd := testConfig(t)
			cmd.Cfg.Secrets = []string{"one", "two"}
			setArgs(t, "terracreds", "list", "host")
			output := captureStdout(t, func() {
				if err := cmd.newCommandActionList(commandContext(t, test.args...)); err != nil {
					t.Fatal(err)
				}
			})
			if output != test.want {
				t.Fatalf("output = %q, want %q", output, test.want)
			}
		})
	}
}

func TestCommandValidationErrors(t *testing.T) {
	cases := []struct {
		name string
		run  func(*Config) error
		args []string
	}{
		{"create", func(cmd *Config) error { return cmd.newCommandActionCreate(commandContext(t)) }, []string{"terracreds", "create"}},
		{"store", func(cmd *Config) error { return cmd.newCommandActionStore(commandContext(t)) }, []string{"terracreds", "store"}},
		{"forget", func(cmd *Config) error { return cmd.newCommandActionForget(commandContext(t)) }, []string{"terracreds", "forget"}},
		{"get", func(cmd *Config) error { return cmd.newCommandActionGet() }, []string{"terracreds", "get"}},
		{"delete", func(cmd *Config) error { return cmd.newCommandActionDelete(commandContext(t)) }, []string{"terracreds", "delete"}},
		{"list", func(cmd *Config) error { return cmd.newCommandActionList(commandContext(t)) }, []string{"terracreds", "list"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			cmd := testConfig(t)
			setArgs(t, test.args...)
			if err := test.run(cmd); err == nil {
				t.Fatal("validation returned nil")
			}
		})
	}

	cmd := testConfig(t)
	setArgs(t, "terracreds", "delete", "host")
	if err := cmd.newCommandActionDelete(commandContext(t)); err != nil {
		t.Fatalf("unexpected delete warning error: %v", err)
	}
	setArgs(t, "terracreds", "list", "host")
	if err := cmd.newCommandActionList(commandContext(t)); err == nil {
		t.Fatal("list without names returned nil")
	}
}
