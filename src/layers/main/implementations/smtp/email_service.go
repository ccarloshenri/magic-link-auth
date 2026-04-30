package smtp

import (
	"fmt"
	"net/smtp"
)

type SMTPEmailService struct {
	host string
	port string
	from string
}

func NewSMTPEmailService(host, port, from string) *SMTPEmailService {
	return &SMTPEmailService{host: host, port: port, from: from}
}

func (s *SMTPEmailService) SendMagicLink(email, link string) error {
	addr := s.host + ":" + s.port
	body := fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: Your magic link\r\n\r\n"+
			"Click the link below to authenticate:\r\n\r\n%s\r\n\r\n"+
			"This link expires in 15 minutes and can only be used once.",
		email, s.from, link,
	)
	if err := smtp.SendMail(addr, nil, s.from, []string{email}, []byte(body)); err != nil {
		return fmt.Errorf("smtp send to %s: %w", email, err)
	}
	return nil
}
