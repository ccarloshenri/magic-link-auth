// Package aws provides AWS-backed implementations (DynamoDB, SES).
// These are stubs — wire up with the AWS SDK when deploying to ECS.
package aws

import (
	"errors"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/models"
)

// DynamoDBMagicLinkRepository is a placeholder for the DynamoDB-backed repository.
type DynamoDBMagicLinkRepository struct{}

func (r *DynamoDBMagicLinkRepository) Save(_ models.MagicLink) error {
	return errors.New("DynamoDBMagicLinkRepository: not implemented")
}

func (r *DynamoDBMagicLinkRepository) FindByToken(_ string) (*models.MagicLink, error) {
	return nil, errors.New("DynamoDBMagicLinkRepository: not implemented")
}

func (r *DynamoDBMagicLinkRepository) MarkAsUsed(_ string) error {
	return errors.New("DynamoDBMagicLinkRepository: not implemented")
}
