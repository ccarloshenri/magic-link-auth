package memory

import "log"

type LogEmailService struct{}

func NewLogEmailService() *LogEmailService {
	return &LogEmailService{}
}

func (s *LogEmailService) SendMagicLink(email, link string) error {
	log.Printf("[EMAIL] To: %s | Magic Link sent (use /auth/validate to authenticate)", email)
	return nil
}
