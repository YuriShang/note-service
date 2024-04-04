package db

import (
	"context"
	"note_service/app/internal/note"
	"note_service/app/pkg/logging"
	"note_service/app/pkg/postgres"

	"github.com/google/uuid"
)

var _ note.Storage = &db{}

type db struct {
	client *postgres.Client
	logger logging.Logger
}

func NewStorage(client *postgres.Client, logger logging.Logger) note.Storage {
	return &db{
		client: client,
		logger: logger,
	}
}

func (s *db) GetNotes(ctx context.Context, userUUID uuid.UUID) (*note.Notes, error) {
	note, err := s.client.GetNotes(ctx, userUUID)
	return note, err
}

func (s *db) Create(ctx context.Context, note note.Note) error {
	err := s.client.CreateNote(ctx, &note)
	return err
}

func (s *db) GetByID(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) (*note.Note, error) {
	note, err := s.client.GetNoteByID(ctx, noteUUID, userUUID)
	return note, err
}

func (s *db) Update(ctx context.Context, note note.Note, userUUID uuid.UUID) error {
	err := s.client.UpdateNote(ctx, &note, userUUID)
	return err
}

func (s *db) Delete(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) error {
	err := s.client.DeleteNote(ctx, noteUUID, userUUID)
	return err
}
