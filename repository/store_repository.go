package repository

import (
	"context"
	"errors"
	"time"

	"github.com/rokoga/filas-backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// ErrorNotFoundStore for store not found
	ErrorNotFoundStore = "Não foi encontrado o estabelecimento"
	// ErrorNotFoundAllStores for store not found
	ErrorNotFoundAllStores = "Erro ao pesquisar nomes dos estabelecimentos"
	// ErrorNotFoundConsumer for consumer not found
	ErrorNotFoundConsumer = "Não foi possível encontrar consumidor"
	// ErrorConsumerExists for consumer already exists
	ErrorConsumerExists = "Consumidor já cadastrado na fila"
	// ErrorParserID for error parsing ID string
	ErrorParserID = "Erro ao fazer parser do ID"
)

// StoreRepositoryImpl implements
type StoreRepositoryImpl struct {
	collection *mongo.Collection
}

// NewStoreRepository implements
func NewStoreRepository(db *mongo.Collection) StoreRepository {
	return &StoreRepositoryImpl{
		collection: db,
	}
}

// Create implements
func (repo *StoreRepositoryImpl) Create(store *domain.Store) (*domain.Store, error) {

	// fmt.Printf("store %v \n", store)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := repo.collection.InsertOne(ctx, store)
	if err != nil {
		return nil, err
	}

	strID := result.InsertedID.(primitive.ObjectID).Hex()

	storeCreated, err := repo.GetStoreByID(strID)
	if err != nil {
		return nil, errors.New(ErrorNotFoundStore)
	}

	// fmt.Printf("storeCreated GetStore %v \n", storeCreated)

	return storeCreated, nil
}

// RemoveStore implements
func (repo *StoreRepositoryImpl) RemoveStore(id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(ErrorParserID)
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	result, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		return errors.New(ErrorNotFoundStore)
	}

	return nil
}

// GetAllStores implements
// TODO: usar find com setProjection
func (repo *StoreRepositoryImpl) GetAllStores() ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New(ErrorNotFoundAllStores)
	}

	var stores []*domain.Store

	err = cursor.All(ctx, &stores)
	if err != nil {
		return nil, err
	}

	var result []string

	for _, value := range stores {
		result = append(result, value.Name)
	}

	return result, nil
}

// GetStoreByID implements
func (repo *StoreRepositoryImpl) GetStoreByID(id string) (*domain.Store, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var storeGotID domain.Store

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(ErrorParserID)
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	err = repo.collection.FindOne(ctx, filter).Decode(&storeGotID)
	if err != nil {
		return nil, errors.New(ErrorNotFoundStore)
	}

	return &storeGotID, nil
}

// GetStore implements
func (repo *StoreRepositoryImpl) GetStore(name string) (*domain.Store, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var storeGotName domain.Store

	// fmt.Printf("Name %v \n", name)

	filter := bson.D{{Key: "name", Value: name}}

	err := repo.collection.FindOne(ctx, filter).Decode(&storeGotName)
	if err != nil {
		return nil, errors.New(ErrorNotFoundStore)
	}

	// fmt.Printf("storeGotName GetStore %v \n", storeGotName)

	return &storeGotName, nil
}

// AddConsumer implements
func (repo *StoreRepositoryImpl) AddConsumer(id string, consumer *domain.Consumer) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store, err := repo.GetStoreByID(id)
	if err != nil {
		return errors.New(ErrorNotFoundStore)
	}

	for _, value := range store.Queue {
		if value.Phone == consumer.Phone {
			return errors.New(ErrorConsumerExists)
		}
	}

	store.Queue = append(store.Queue, consumer)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(ErrorParserID)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "queue", Value: store.Queue}}},
	}

	_, err = repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// RemoveConsumer implements
func (repo *StoreRepositoryImpl) RemoveConsumer(id string, phone string) error {
	store, err := repo.GetStoreByID(id)
	if err != nil {
		return errors.New(ErrorNotFoundStore)
	}

	for i, consumer := range store.Queue {
		if consumer.Phone == phone {
			copy(store.Queue[i:], store.Queue[i+1:])
			store.Queue[len(store.Queue)-1] = nil
			store.Queue = store.Queue[:len(store.Queue)-1]

			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return errors.New(ErrorParserID)
			}

			filter := bson.D{{Key: "_id", Value: oid}}
			update := bson.D{
				{Key: "$set", Value: bson.D{{Key: "queue", Value: store.Queue}}},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err = repo.collection.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New(ErrorNotFoundConsumer)
}

// GetConsumer implements
func (repo *StoreRepositoryImpl) GetConsumer(id string, phone string) (int, *domain.Consumer, error) {

	store, err := repo.GetStoreByID(id)
	if err != nil {
		return 0, nil, errors.New(ErrorNotFoundStore)
	}

	for i, consumer := range store.Queue {
		if consumer.Phone == phone {
			return i, consumer, nil
		}
	}

	return 0, nil, errors.New(ErrorNotFoundConsumer)
}

// GetAllConsumers implements
func (repo *StoreRepositoryImpl) GetAllConsumers(id string) ([]*domain.Consumer, error) {

	store, err := repo.GetStoreByID(id)
	if err != nil {
		return nil, errors.New(ErrorNotFoundStore)
	}

	return store.Queue, nil
}
