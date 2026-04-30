package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"magic-link-auth/src/layers/main/enums"
	"magic-link-auth/src/layers/main/models"
)

type DynamoDBMagicLinkDAO struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBMagicLinkDAO(client *dynamodb.Client, tableName string) *DynamoDBMagicLinkDAO {
	return &DynamoDBMagicLinkDAO{client: client, tableName: tableName}
}

func (d *DynamoDBMagicLinkDAO) Save(link models.MagicLink) error {
	item, err := attributevalue.MarshalMap(link)
	if err != nil {
		return fmt.Errorf("marshal magic link: %w", err)
	}
	_, err = d.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put item: %w", err)
	}
	return nil
}

func (d *DynamoDBMagicLinkDAO) FindByToken(token string) (*models.MagicLink, error) {
	result, err := d.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"Token": &types.AttributeValueMemberS{Value: token},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}
	if result.Item == nil {
		return nil, fmt.Errorf("token not found")
	}
	var link models.MagicLink
	if err := attributevalue.UnmarshalMap(result.Item, &link); err != nil {
		return nil, fmt.Errorf("unmarshal magic link: %w", err)
	}
	return &link, nil
}

func (d *DynamoDBMagicLinkDAO) MarkAsUsed(token string) error {
	_, err := d.client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"Token": &types.AttributeValueMemberS{Value: token},
		},
		UpdateExpression: aws.String("SET #s = :used"),
		ExpressionAttributeNames: map[string]string{
			"#s": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":used": &types.AttributeValueMemberS{Value: string(enums.Used)},
		},
	})
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}
	return nil
}
