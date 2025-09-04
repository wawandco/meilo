package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/wawandco/meilo/internal/database"
	"github.com/wawandco/meilo/internal/models"
)

type EmailService struct {
	conn database.Connection
}

func NewEmailService(conn database.Connection) *EmailService {
	return &EmailService{conn: conn}
}

func (s *EmailService) Save(e models.Email) error {
	db, err := s.conn()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("meilo: failed to start transaction %w", err)
	}

	// Convert arrays to JSON
	recipientsJSON, err := json.Marshal(e.Recipients)
	if err != nil {
		return fmt.Errorf("meilo: failed to marshal recipients: %w", err)
	}

	ccJSON, err := json.Marshal(e.CC)
	if err != nil {
		return fmt.Errorf("meilo: failed to marshal CC: %w", err)
	}

	bccJSON, err := json.Marshal(e.BCC)
	if err != nil {
		return fmt.Errorf("meilo: failed to marshal BCC: %w", err)
	}

	// 1. Insert into emails table with JSON data
	result, err := tx.Exec("INSERT INTO emails (subject, sender, recipients, cc, bcc) VALUES ($1, $2, $3, $4, $5)",
		e.Subject,
		e.Sender,
		string(recipientsJSON),
		string(ccJSON),
		string(bccJSON),
	)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("meilo: failed to rollback transaction inserting email: %w", errRollback)
		}
		return fmt.Errorf("meilo: failed to save email: %w", err)
	}

	emailID, err := result.LastInsertId()
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("meilo: failed to rollback transaction retrieving last id: %w", errRollback)
		}
		return fmt.Errorf("meilo: failed to retrieve the last ID %w", err)
	}

	// 2. Insert each body part
	for _, b := range e.Bodies {
		if _, err = tx.Exec("INSERT INTO bodies (email_id, content_type, content) VALUES ($1, $2, $3)",
			emailID,
			b.ContentType,
			b.Content,
		); err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				return fmt.Errorf("meilo: failed to rollback transaction inserting email bodies: %w", errRollback)
			}
			return fmt.Errorf("meilo: failed to inster email bodies %w", err)
		}
	}

	// 3. Insert each attachment
	for _, a := range e.Attachments {
		if _, err = tx.Exec("INSERT INTO attachments (email_id, name, content_type, data) VALUES ($1, $2, $3, $4)",
			emailID,
			a.Name,
			a.ContentType,
			a.Data,
		); err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				return fmt.Errorf("meilo: failed to rollback transaction inserting email attachments: %w", errRollback)
			}
			return fmt.Errorf("saveEmail: insert attachments: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("meilo: failed to commit transaction inserting email information: %w", err)
	}

	return nil
}

func (s *EmailService) List() ([]models.Email, error) {
	db, err := s.conn()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT e.id, e.subject, e.sender, e.recipients, b.content, e.received_at
		FROM emails e
		LEFT JOIN bodies b ON e.id = b.email_id AND b.content_type = 'text/html'
		ORDER BY e.received_at DESC
		LIMIT 5;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("meilo: failed to list emails: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Error("meilo: failed to close rows", "error", err)
		}
	}(rows)

	var emails []models.Email
	for rows.Next() {
		var (
			e              models.Email
			recipientsJSON string
		)
		if err := rows.Scan(
			&e.ID,
			&e.Subject,
			&e.Sender,
			&recipientsJSON,
			&e.Body,
			&e.ReceivedAt,
		); err != nil {
			return nil, fmt.Errorf("meilo: failed to scan email rows: %w", err)
		}

		// Unmarshal JSON arrays back to slices
		if err := json.Unmarshal([]byte(recipientsJSON), &e.Recipients); err != nil {
			return nil, fmt.Errorf("meilo: failed to unmarshal recipients: %w", err)
		}

		emails = append(emails, e)
	}

	return emails, nil
}

func (s *EmailService) FindById(id int64) (models.Email, error) {
	db, err := s.conn()
	if err != nil {
		return models.Email{}, err
	}

	var (
		e             models.Email
		recipientsCSV string
	)

	err = db.QueryRow(`
		SELECT e.id, e.subject, e.sender, e.recipients, b.content, e.received_at
		FROM emails e
		LEFT JOIN bodies b ON e.id = b.email_id AND b.content_type = 'text/html'
		WHERE e.id = $1
		LIMIT 1;
	`, id).Scan(&e.ID, &e.Subject, &e.Sender, &recipientsCSV, &e.Body, &e.ReceivedAt)
	if err != nil {
		return models.Email{}, err
	}

	e.Recipients = strings.Split(recipientsCSV, ",")
	return e, nil
}

func (s *EmailService) DeleteAll() error {
	db, err := s.conn()
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM emails")
	if err != nil {
		return fmt.Errorf("meilo: failed to delete all emails: %w", err)
	}

	return nil
}

func (s *EmailService) GetLatestEmailByRecipient(email string) (models.Email, error) {
	db, err := s.conn()
	if err != nil {
		return models.Email{}, err
	}

	// Use SQLite's JSON functions to search within the JSON arrays
	query := `
        SELECT e.id, e.subject, e.sender, e.recipients, b.content, e.received_at
        FROM emails e
        LEFT JOIN bodies b ON e.id = b.email_id AND b.content_type = 'text/html'
        WHERE EXISTS (
            SELECT 1 FROM json_each(e.recipients) WHERE value = $1
        )
        ORDER BY e.received_at DESC
        LIMIT 1
    `

	var (
		e              models.Email
		recipientsJSON string
	)

	err = db.QueryRow(query, email).Scan(
		&e.ID,
		&e.Subject,
		&e.Sender,
		&recipientsJSON,
		&e.Body,
		&e.ReceivedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Email{}, fmt.Errorf("meilo: no email found for recipient %s", email)
		}
		return models.Email{}, fmt.Errorf("meilo: failed to get latest email for recipient: %w", err)
	}

	// Unmarshal JSON arrays back to slices
	if err := json.Unmarshal([]byte(recipientsJSON), &e.Recipients); err != nil {
		return models.Email{}, fmt.Errorf("meilo: failed to unmarshal recipients: %w", err)
	}

	return e, nil
}
