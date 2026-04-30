package interfaces

type TokenService interface {
	Generate() (string, error)
}
