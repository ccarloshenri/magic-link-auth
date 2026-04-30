package memory

import (
	"errors"
	"sync"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/enums"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/models"
)

type InMemoryMagicLinkDAO struct {
	mu    sync.RWMutex
	store map[string]models.MagicLink
}

func NewInMemoryMagicLinkDAO() *InMemoryMagicLinkDAO {
	return &InMemoryMagicLinkDAO{
		store: make(map[string]models.MagicLink),
	}
}

func (d *InMemoryMagicLinkDAO) Save(link models.MagicLink) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.store[link.Token] = link
	return nil
}

func (d *InMemoryMagicLinkDAO) FindByToken(token string) (*models.MagicLink, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	link, ok := d.store[token]
	if !ok {
		return nil, errors.New("token not found")
	}
	return &link, nil
}

func (d *InMemoryMagicLinkDAO) MarkAsUsed(token string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	link, ok := d.store[token]
	if !ok {
		return errors.New("token not found")
	}
	link.Status = enums.Used
	d.store[token] = link
	return nil
}
