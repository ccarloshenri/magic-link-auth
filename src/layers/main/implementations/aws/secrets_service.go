package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerService struct {
	client *secretsmanager.Client
}

func NewSecretsManagerService(client *secretsmanager.Client) *SecretsManagerService {
	return &SecretsManagerService{client: client}
}

func (s *SecretsManagerService) GetSecret(name string) (string, error) {
	result, err := s.client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	})
	if err != nil {
		return "", fmt.Errorf("get secret %q: %w", name, err)
	}
	if result.SecretString == nil {
		return "", fmt.Errorf("secret %q has no string value", name)
	}
	return *result.SecretString, nil
}
