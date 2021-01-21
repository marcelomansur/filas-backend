package service

import "github.com/rokoga/filas-backend/domain"

// StoreService - Provides a Store services layer
type StoreService interface {
	Create(name string) (*domain.Store, error)
	RemoveStore(id string) error
	GetAllStores() ([]string, error)
	GetStore(name string) (*domain.Store, error)
	GetStoreByID(id string) (*domain.Store, error)
	AddConsumer(id, name, phone, status string) (string, error)
	RemoveConsumer(id string, phone string) error
	GetConsumer(id string, phone string) (int, *domain.Consumer, error)
	GetAllConsumers(id string) ([]*domain.Consumer, error)
	ValidateConsumer(storeName, accessKey string) (int, *domain.Consumer, error)
}
