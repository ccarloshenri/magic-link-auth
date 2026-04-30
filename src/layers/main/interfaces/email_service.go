package interfaces

type EmailService interface {
	SendMagicLink(email, link string) error
}
