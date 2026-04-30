package interfaces

type AuthTokenService interface {
	GenerateJWT(email string) (string, error)
}
