package note

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	NoteUUID   *uuid.UUID `json:"id,omitempty"`
	UserUUID   *uuid.UUID `json:"user_id,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty"`
	Text       *string    `json:"text,omitempty"`
	Public     *bool      `json:"public,omitempty"`
}

type Notes struct {
	Notes []Note `json:"notes" bson:"notes,omitempty"`
}

func CreateNote(dto CreateNoteDTO) Note {
	return Note{
		UserUUID: dto.UserUUID,
		Text:     dto.Text,
		Public:   dto.Public,
	}
}

func UpdatedNote(dto UpdateNoteDTO) Note {
	return Note{
		NoteUUID: dto.NoteUUID,
		Text:     dto.Text,
		Public:   dto.Public,
	}
}

type CreateNoteDTO struct {
	UserUUID *uuid.UUID `json:"id"`
	Text     *string    `json:"text"`
	Public   *bool      `json:"public"`
}

type UpdateNoteDTO struct {
	NoteUUID *uuid.UUID `json:"id"`
	Text     *string    `json:"text,omitempty"`
	Public   *bool      `json:"public,omitempty"`
}
