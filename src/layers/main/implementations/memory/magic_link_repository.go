package memory

import (
	"errors"
	"sync"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/enums"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/models"
)

type InMemoryMagicLinkRepository struct {
	mu    sync.RWMutex
	store map[string]models.MagicLink
}

func NewInMemoryMagicLinkRepository() *InMemoryMagicLinkRepository {
	return &InMemoryMagicLinkRepository{
		store: make(map[string]models.MagicLink),
	}
}

func (r *InMemoryMagicLinkRepository) Save(link models.MagicLink) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[link.Token] = link
	return nil
}

func (r *InMemoryMagicLinkRepository) FindByToken(token string) (*models.MagicLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	link, ok := r.store[token]
	if !ok {
		return nil, errors.New("token not found")
	}
	return &link, nil
}

func (r *InMemoryMagicLinkRepository) MarkAsUsed(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	link, ok := r.store[token]
	if !ok {
		return errors.New("token not found")
	}
	link.Status = enums.Used
	r.store[token] = link
	return nil
}
