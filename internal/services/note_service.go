package services

import (
	"context"
	"time"

	"github.com/1206yaya/go-note-api/internal/models"
	"github.com/1206yaya/go-note-api/internal/repositories"
	"github.com/google/uuid"
)

type NoteService interface {
	CreateNote(ctx context.Context, title, body string) (*models.Note, error)
	GetNote(ctx context.Context, id string) (*models.Note, error)
	UpdateNote(ctx context.Context, id, title, body string) (*models.Note, error)
	DeleteNote(ctx context.Context, id string) error
	ListNotes(ctx context.Context) ([]*models.Note, error)
}

type noteService struct {
	repo repositories.NoteRepository
}

func NewNoteService(repo repositories.NoteRepository) NoteService {
	return &noteService{
		repo: repo,
	}
}

func (s *noteService) CreateNote(ctx context.Context, title, body string) (*models.Note, error) {
	now := time.Now()
	note := &models.Note{
		ID:        uuid.New().String(),
		Title:     title,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.repo.CreateNote(ctx, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *noteService) GetNote(ctx context.Context, id string) (*models.Note, error) {
	return s.repo.GetNote(ctx, id)
}

func (s *noteService) UpdateNote(ctx context.Context, id, title, body string) (*models.Note, error) {
	note, err := s.repo.GetNote(ctx, id)
	if err != nil {
		return nil, err
	}

	note.Title = title
	note.Body = body
	note.UpdatedAt = time.Now()

	err = s.repo.UpdateNote(ctx, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *noteService) DeleteNote(ctx context.Context, id string) error {
	return s.repo.DeleteNote(ctx, id)
}

func (s *noteService) ListNotes(ctx context.Context) ([]*models.Note, error) {
	return s.repo.ListNotes(ctx)
}
