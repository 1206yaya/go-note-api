package repositories

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/1206yaya/go-note-api/internal/models"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, note *models.Note) error
	GetNote(ctx context.Context, id string) (*models.Note, error)
	UpdateNote(ctx context.Context, note *models.Note) error
	DeleteNote(ctx context.Context, id string) error
	ListNotes(ctx context.Context) ([]*models.Note, error)
}

type DynamoDBNoteRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBNoteRepository(client *dynamodb.Client, tableName string) *DynamoDBNoteRepository {
	return &DynamoDBNoteRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBNoteRepository) CreateNote(ctx context.Context, note *models.Note) error {
	item, err := attributevalue.MarshalMap(note)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	return err
}

func (r *DynamoDBNoteRepository) GetNote(ctx context.Context, id string) (*models.Note, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var note models.Note
	err = attributevalue.UnmarshalMap(result.Item, &note)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *DynamoDBNoteRepository) UpdateNote(ctx context.Context, note *models.Note) error {
	item, err := attributevalue.MarshalMap(note)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	return err
}

func (r *DynamoDBNoteRepository) DeleteNote(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})

	return err
}

func (r *DynamoDBNoteRepository) ListNotes(ctx context.Context) ([]*models.Note, error) {
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		return nil, fmt.Errorf("DYNAMODB_TABLE environment variable is not set")
	}

	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	var notes []*models.Note
	err = attributevalue.UnmarshalListOfMaps(result.Items, &notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}
