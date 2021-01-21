package repository

import (
	"github.com/rokoga/filas-backend/domain"
)

// StoreRepository - Repository for persisting a Store
type StoreRepository interface {
	Create(store *domain.Store) (*domain.Store, error)
	RemoveStore(id string) error
	GetAllStores() ([]string, error)
	GetStoreByID(id string) (*domain.Store, error)
	GetStore(name string) (*domain.Store, error)
	AddConsumer(id string, consumer *domain.Consumer) error
	RemoveConsumer(id string, phone string) error
	GetConsumer(id string, phone string) (int, *domain.Consumer, error)
	GetAllConsumers(id string) ([]*domain.Consumer, error)
	ValidateConsumer(storeName, accessKey string) (int, *domain.Consumer, error)
}

// app.filas/outback/token?=24238971alkajrealm
