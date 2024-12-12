package repo

import (
	"context"
	"sync"

	"github.com/cutlery47/email-service/internal/models"
)

// мок-репозиторий для ранней стадии тестирования
type MockRepository struct {
	data []models.UserData
	mu   *sync.Mutex
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		data: make([]models.UserData, 0),
		mu:   &sync.Mutex{},
	}
}

func (ms *MockRepository) Create(ctx context.Context, user models.UserData) error {
	ms.mu.Lock()
	ms.data = append(ms.data, user)
	ms.mu.Unlock()

	return nil
}
