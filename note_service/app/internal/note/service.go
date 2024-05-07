package note

import (
	"context"
	"errors"
	"fmt"
	"note_service/app/internal/apperror"
	"note_service/app/pkg/logging"

	"github.com/google/uuid"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(noteStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: noteStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto CreateNoteDTO) (string, error)
	GetMany(ctx context.Context, userUUID uuid.UUID) (*Notes, error)
	GetOne(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) (*Note, error)
	Update(ctx context.Context, dto UpdateNoteDTO, userUUID uuid.UUID) error
	Delete(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) error
}

func (s service) Create(ctx context.Context, dto CreateNoteDTO) (noteUUID string, err error) {
	note := CreateNote(dto)
	err = s.storage.Create(ctx, note)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return noteUUID, err
		}
		return noteUUID, fmt.Errorf("failed to create note. error: %w", err)
	}

	return noteUUID, nil
}

func (s service) GetOne(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) (n *Note, err error) {
	n, err = s.storage.GetByID(ctx, noteUUID, userUUID)

	if err != nil {
		if errors.Is(err, apperror.ErrForbidden) || errors.Is(err, apperror.ErrNotFound) {
			return n, err
		}
		return n, fmt.Errorf("failed to find note by uuid. error: %w", err)
	}
	return n, nil
}

func (s service) GetMany(ctx context.Context, userUUID uuid.UUID) (n *Notes, err error) {
	n, err = s.storage.GetNotes(ctx, userUUID)

	if err != nil {
		return n, err
	}
	return n, nil
}

func (s service) Update(ctx context.Context, dto UpdateNoteDTO, userUUID uuid.UUID) error {
	note := UpdatedNote(dto)
	err := s.storage.Update(ctx, note, userUUID)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update note. error: %w", err)
	}
	return nil
}

func (s service) Delete(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) error {
	err := s.storage.Delete(ctx, noteUUID, userUUID)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete note. error: %w", err)
	}
	return err
}
