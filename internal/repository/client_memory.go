package repository

import (
	"sync"
	"translator/internal/model"
)

type ClientRepositoryMemory struct {
	storage map[int64]*model.Client
	mu      sync.RWMutex
}

func NewClientRepositoryMemory() *ClientRepositoryMemory {
	return &ClientRepositoryMemory{storage: make(map[int64]*model.Client)}
}

func (c *ClientRepositoryMemory) SaveIfAbsent(id int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.storage[id] != nil {
		return false
	}

	c.storage[id] = model.NewClient(id)
	return true
}

func (c *ClientRepositoryMemory) GetLang(id int64) (model.Language, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	client := c.storage[id]
	if client == nil {
		return "", ErrClientNotFound
	}

	return client.Language, nil
}

func (c *ClientRepositoryMemory) SetLang(id int64, lang model.Language) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	client := c.storage[id]
	if client == nil {
		return ErrClientNotFound
	}

	client.Language = lang
	return nil
}
