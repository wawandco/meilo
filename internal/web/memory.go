package web

import (
	"fmt"
	"strings"
	"sync"
)

type MemoryEmailStore struct {
	emails []Email
	mutex  sync.RWMutex
	nextID int64
}

func NewMemoryEmailStore() *MemoryEmailStore {
	return &MemoryEmailStore{
		emails: make([]Email, 0),
		nextID: 1,
	}
}

func (m *MemoryEmailStore) Add(email Email) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.emails = append(m.emails, email)
	m.nextID++
}

// List returns all emails in reverse chronological order (latest first)
func (m *MemoryEmailStore) List(limit int) []Email {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if limit <= 0 || limit > len(m.emails) {
		limit = len(m.emails)
	}

	// Return emails in reverse order (latest first)
	result := make([]Email, limit)
	start := len(m.emails) - limit
	for i := 0; i < limit; i++ {
		result[i] = m.emails[start+limit-1-i]
	}

	return result
}

// GetByID retrieves a specific email by ID
func (m *MemoryEmailStore) GetByID(id int64) (*Email, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, email := range m.emails {
		if email.ID == id {
			email.SetHTMLBody()
			return &email, nil
		}
	}
	return nil, fmt.Errorf("email with ID %d not found", id)
}

// Count returns the total number of emails stored
func (m *MemoryEmailStore) Count() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.emails)
}

// Clear removes all emails from memory
func (m *MemoryEmailStore) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.emails = m.emails[:0]
	m.nextID = 1
}

func (m *MemoryEmailStore) GetLatestEmailByRecipient(recipient string) *Email {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	recipient = strings.TrimSpace(strings.ToLower(recipient))
	for _, e := range m.emails {
		for _, r := range e.Recipients {
			if strings.TrimSpace(strings.ToLower(r)) == recipient {
				e.SetHTMLBody()
				return &e
			}
		}
	}

	return nil
}
