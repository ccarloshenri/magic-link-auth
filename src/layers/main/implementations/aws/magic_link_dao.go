// Package aws provides AWS-backed implementations (DynamoDB, SES).
// These are stubs — wire up with the AWS SDK when deploying to ECS.
package aws

import (
	"errors"

	"magic-link-auth/src/layers/main/models"
)

// DynamoDBMagicLinkDAO is a placeholder for the DynamoDB-backed DAO.
type DynamoDBMagicLinkDAO struct{}

func (d *DynamoDBMagicLinkDAO) Save(_ models.MagicLink) error {
	return errors.New("DynamoDBMagicLinkDAO: not implemented")
}

func (d *DynamoDBMagicLinkDAO) FindByToken(_ string) (*models.MagicLink, error) {
	return nil, errors.New("DynamoDBMagicLinkDAO: not implemented")
}

func (d *DynamoDBMagicLinkDAO) MarkAsUsed(_ string) error {
	return errors.New("DynamoDBMagicLinkDAO: not implemented")
}
