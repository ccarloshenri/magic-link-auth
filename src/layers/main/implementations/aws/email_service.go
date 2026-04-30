package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SESEmailService struct {
	client      *sesv2.Client
	senderEmail string
}

func NewSESEmailService(client *sesv2.Client, senderEmail string) *SESEmailService {
	return &SESEmailService{client: client, senderEmail: senderEmail}
}

func (s *SESEmailService) SendMagicLink(email, link string) error {
	body := fmt.Sprintf(
		"Click the link below to authenticate:\n\n%s\n\nThis link expires in 15 minutes and can only be used once.",
		link,
	)
	_, err := s.client.SendEmail(context.Background(), &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.senderEmail),
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String("Your magic link"),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(body),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("send email to %s: %w", email, err)
	}
	return nil
}
