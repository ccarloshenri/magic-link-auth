package interfaces

type SecretsService interface {
	GetSecret(name string) (string, error)
}
