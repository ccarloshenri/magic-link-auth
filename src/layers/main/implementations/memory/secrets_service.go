package memory

import "fmt"

type EnvSecretsService struct {
	secrets map[string]string
}

func NewEnvSecretsService(secrets map[string]string) *EnvSecretsService {
	return &EnvSecretsService{secrets: secrets}
}

func (s *EnvSecretsService) GetSecret(name string) (string, error) {
	if v, ok := s.secrets[name]; ok {
		return v, nil
	}
	return "", fmt.Errorf("secret %q not found", name)
}
