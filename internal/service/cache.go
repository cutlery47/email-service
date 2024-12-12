package service

import (
	"sync"

	"github.com/cutlery47/email-service/internal/models"
)

// Кэш для хранения данных пользователей до момента подтверждения почты
type Cache interface {
	Put(user models.CachedUserData) error
	Get(mail string) (models.CachedUserData, error)
}

type MapCache struct {
	data map[string]models.CachedUserData
	mu   *sync.RWMutex
}

func NewMapCache() (*MapCache, error) {
	return &MapCache{
		data: make(map[string]models.CachedUserData),
		mu:   &sync.RWMutex{},
	}, nil
}

func (mc *MapCache) Put(user models.CachedUserData) error {
	mc.mu.Lock()
	mc.data[user.Mail] = user
	mc.mu.Unlock()

	return nil
}

func (mc *MapCache) Get(mail string) (models.CachedUserData, error) {
	mc.mu.RLock()
	v, ok := mc.data[mail]
	mc.mu.RUnlock()

	if !ok {
		return models.CachedUserData{}, ErrCacheNotFound
	}

	return v, nil
}
