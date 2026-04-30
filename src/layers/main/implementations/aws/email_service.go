package aws

import "errors"

// SESEmailService is a placeholder for the Amazon SES-backed email service.
type SESEmailService struct{}

func (s *SESEmailService) SendMagicLink(_, _ string) error {
	return errors.New("SESEmailService: not implemented")
}
