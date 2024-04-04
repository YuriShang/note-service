package note

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	Create(ctx context.Context, note Note) error
	GetByID(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) (*Note, error)
	GetNotes(ctx context.Context, userUUID uuid.UUID) (*Notes, error)
	Update(ctx context.Context, note Note, userUUID uuid.UUID) error
	Delete(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) error
}
