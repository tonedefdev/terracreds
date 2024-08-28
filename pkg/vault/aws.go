package vault

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/tonedefdev/terracreds/pkg/helpers"
)

type AwsSecretsManager struct {
	Description string
	Region      string
	SecretName  string
}

func (asm *AwsSecretsManager) getAwsSecetsManager() *secretsmanager.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		helpers.CheckError(err)
	}

	svc := secretsmanager.NewFromConfig(cfg)
	return svc
}

func (asm *AwsSecretsManager) Create(secretValue string, method string) error {
	svc := asm.getAwsSecetsManager()

	if method == "Updated" {
		input := &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(asm.SecretName),
			SecretString: aws.String(secretValue),
		}

		_, err := svc.PutSecretValue(ctx, input)
		if err != nil {
			return err
		}

		return err
	}

	input := &secretsmanager.CreateSecretInput{
		Description:  aws.String(asm.Description),
		Name:         aws.String(asm.SecretName),
		SecretString: aws.String(secretValue),
	}

	_, err := svc.CreateSecret(ctx, input)
	if err != nil {
		return err
	}

	return err
}

func (asm *AwsSecretsManager) Delete() error {
	svc := asm.getAwsSecetsManager()

	input := &secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(7),
		SecretId:             aws.String(asm.SecretName),
	}

	_, err := svc.DeleteSecret(ctx, input)
	if err != nil {
		return err
	}

	return err
}

func (asm *AwsSecretsManager) Get() ([]byte, error) {
	svc := asm.getAwsSecetsManager()

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(asm.SecretName),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}

	return []byte(*result.SecretString), err
}

func (asm *AwsSecretsManager) List(secretNames []string) ([]string, error) {
	var secretValues []string
	svc := asm.getAwsSecetsManager()

	for _, secret := range secretNames {
		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secret),
		}

		result, err := svc.GetSecretValue(ctx, input)
		if err != nil {
			return nil, err
		}

		secretValues = append(secretValues, *result.SecretString)
	}

	return secretValues, nil
}
