package postgres

import (
	"context"
	"database/sql"
	"fmt"
	e "note_service/app/internal/apperror"
	"note_service/app/internal/note"
	"note_service/app/pkg/logging"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type Client struct {
	logger logging.Logger
	db     *sql.DB
}

func NewClient(ctx context.Context, host, port, username, password, database string, logger logging.Logger) (*Client, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	logger.Info("postgresql db initiated")
	return &Client{
		logger: logger,
		db:     db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return rows, nil
}

func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}
	return result, nil
}

func (c *Client) CreateNote(ctx context.Context, note *note.Note) error {
	ID := uuid.New()
	currentTime := time.Now()
	note.NoteUUID = &ID
	note.CreateTime = &currentTime
	query := `INSERT INTO notes (id, user_id, text, public, create_time)
               VALUES ($1, $2, $3, $4, $5)`
	_, err := c.Exec(ctx, query, note.NoteUUID, note.UserUUID, note.Text, note.Public, note.CreateTime)
	if err != nil {
		return fmt.Errorf("error creating note: %w", err)
	}
	return nil
}

func (c *Client) GetNotes(ctx context.Context, userUUID uuid.UUID) (*note.Notes, error) {
	var notes note.Notes
	query := fmt.Sprintf(`
		SELECT id, user_id, create_time, text, public
		FROM notes
		WHERE (public = true OR (user_id = '%s' AND public = false))
	`, userUUID)
	rows, err := c.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting notes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var note_ note.Note
		if err := rows.Scan(&note_.NoteUUID, &note_.UserUUID, &note_.CreateTime, &note_.Text, &note_.Public); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil // Notes not found
			}
			return nil, fmt.Errorf("error getting note by ID: %w", err)
		}
		notes.Notes = append(notes.Notes, note_)
	}
	if len(notes.Notes) == 0 {
		return &notes, e.ErrNotFound
	}
	return &notes, nil
}

func (c *Client) GetNoteByID(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) (*note.Note, error) {
	var note note.Note
	query := `
		SELECT id, user_id, create_time, text, public FROM notes WHERE id = $1 
	`
	row := c.db.QueryRowContext(ctx, query, noteUUID)
	if err := row.Scan(&note.NoteUUID, &note.UserUUID, &note.CreateTime, &note.Text, &note.Public); err != nil {
		if err == sql.ErrNoRows {
			return nil, e.ErrNotFound // Note not found
		}
		return nil, fmt.Errorf("error getting note by ID: %w", err)
	}
	if *note.UserUUID != userUUID && !*note.Public {
		return &note, e.ErrForbidden
	}
	return &note, nil
}

func (c *Client) UpdateNote(ctx context.Context, note *note.Note, userUUID uuid.UUID) error {
	_, err := c.GetNoteByID(ctx, *note.NoteUUID, userUUID)
	if err != nil {
		return err
	}
	updateQuery := "UPDATE notes SET"

	if note.Text != nil {
		updateQuery += fmt.Sprintf(" text = '%s',", *note.Text)
	}

	if note.Public != nil {
		updateQuery += fmt.Sprintf(" public = '%t',", *note.Public)
	}

	// Remove trailing comma
	updateQuery = updateQuery[:len(updateQuery)-1]

	// Add WHERE clause for noteID
	updateQuery += fmt.Sprintf(" WHERE id = '%s'", *note.NoteUUID)

	stmt, err := c.db.Prepare(updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("error updating note: %w", err)
	}

	return nil
}

func (c *Client) DeleteNote(ctx context.Context, noteUUID uuid.UUID, userUUID uuid.UUID) error {
	_, err := c.GetNoteByID(ctx, noteUUID, userUUID)
	if err != nil {
		return err
	}
	query := `DELETE FROM notes WHERE id = $1`
	_, err = c.Exec(ctx, query, noteUUID)
	if err != nil {
		return fmt.Errorf("error deleting note: %w", err)
	}
	return nil
}
