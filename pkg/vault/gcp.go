package vault

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/tonedefdev/terracreds/pkg/helpers"
)

var ctx = context.Background()

type GCPSecretsManager struct {
	ProjectId  string
	SecretId   string
	SecretList []string
}

func (gcp *GCPSecretsManager) getClient() *secretmanager.Client {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		helpers.CheckError(err)
	}
	defer client.Close()
	return client
}

// formatSecretName replaces the periods from the hostname with dashes
// since GCP can't store secrets that contain periods
func formatGcpSecretName(secretName string) string {
	hostname := strings.Replace(secretName, ".", "-", -1)
	return hostname
}

func (gcp *GCPSecretsManager) Create(secretValue string, method string) error {
	client := gcp.getClient()
	secretId := formatGcpSecretName(gcp.SecretId)

	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", gcp.ProjectId),
		SecretId: secretId,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	secret, err := client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		return err
	}

	payload := []byte(secretValue)

	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	print(secret.Name)

	_, err = client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return err
	}

	return err
}

func (gcp *GCPSecretsManager) Delete() error {
	client := gcp.getClient()
	secretId := formatGcpSecretName(gcp.SecretId)

	destroySecretVersionReq := secretmanagerpb.DestroySecretVersionRequest{
		Name: fmt.Sprintf("/projects/%s/secrets/%s/versions/latest", gcp.ProjectId, secretId),
	}

	_, err := client.DestroySecretVersion(ctx, &destroySecretVersionReq)
	if err != nil {
		return err
	}

	return err
}

func (gcp *GCPSecretsManager) Get() ([]byte, error) {
	client := gcp.getClient()
	secretId := formatGcpSecretName(gcp.SecretId)

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("/projects/%s/secrets/%s/versions/latest", gcp.ProjectId, secretId),
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, err
	}

	return result.Payload.Data, err
}

func (gcp *GCPSecretsManager) List(secretNames []string) ([]string, error) {
	return nil, nil
}
