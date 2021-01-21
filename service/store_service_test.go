package service

import (
	"errors"
	"testing"

	"github.com/rokoga/filas-backend/repository"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	svc := NewStoreMockServiceImpl()

	tests := []struct {
		name       string
		resultURL  string
		resultName string
		err        error
	}{
		{name: "Outback", resultURL: "http://app.filas.com/outback", resultName: "Outback", err: nil},
		{name: "Jeronimo", resultURL: "http://app.filas.com/jeronimo", resultName: "Jeronimo", err: nil},
		{name: "", resultURL: "", resultName: "", err: errors.New(ErrorArgumentNotValidAddStore)},
	}

	for _, test := range tests {
		store, err := svc.Create(test.name)
		if err == nil {
			assert.NotNil(t, store)
			assert.Equal(t, test.resultURL, store.URLName)
			assert.Equal(t, test.resultName, store.Name)
			assert.NotEmpty(t, store.ID)
		} else {
			assert.Equal(t, test.err, err)
			assert.Nil(t, store)
		}
	}

	// fmt.Printf("Store: %v", store)
}

func TestRemoveStore(t *testing.T) {

	svc := NewStoreMockServiceImpl()

	name := "Outback"

	store, err := svc.Create(name)

	assert.Nil(t, err)
	assert.NotNil(t, store)

	tests := []struct {
		id  string
		err error
	}{
		{id: store.ID, err: nil},
		{id: store.ID, err: errors.New(repository.ErrorNotFoundStore)},
		{id: "fakeID", err: errors.New(repository.ErrorNotFoundStore)},
		{id: "", err: errors.New(ErrorArgumentNotValidRemoveStore)},
	}

	for _, test := range tests {
		err := svc.RemoveStore(test.id)
		assert.Equal(t, test.err, err)
	}

}

func TestGetStore(t *testing.T) {

	svc := NewStoreMockServiceImpl()

	name := "Outback"

	store, err := svc.Create(name)

	assert.Nil(t, err)
	assert.NotNil(t, store)

	tests := []struct {
		name       string
		resultURL  string
		resultName string
		err        error
	}{
		{name: "Outback", resultURL: "http://app.filas.com/outback", resultName: "Outback", err: nil},
		{name: "Jeronimo", resultURL: "", resultName: "", err: errors.New(repository.ErrorNotFoundStore)},
		{name: "", resultURL: "", resultName: "", err: errors.New(ErrorArgumentNotValidGetStore)},
	}

	for _, test := range tests {
		store, err := svc.GetStore(test.name)
		if err == nil {
			assert.NotNil(t, store)
			assert.Equal(t, test.resultURL, store.URLName)
			assert.Equal(t, test.resultName, store.Name)
			assert.NotEmpty(t, store.ID)
		} else {
			assert.Equal(t, test.err, err)
			assert.Nil(t, store)
		}
	}

}

func TestAddConsumer(t *testing.T) {

	svc := NewStoreMockServiceImpl()

	name := "Outback"

	store, err := svc.Create(name)

	assert.Nil(t, err)
	assert.NotNil(t, store)

	tests := []struct {
		id     string
		name   string
		phone  string
		status string
		err    error
	}{
		{id: store.ID, name: "Fulano", phone: "011998989898", status: "Na fila", err: nil},
		{id: store.ID, name: "Ciclano", phone: "011922222222", status: "Na fila", err: nil},
		{id: store.ID, name: "", phone: "", status: "Na fila", err: errors.New(ErrorArgumentNotValidAddConsumer)},
		{id: "", name: "Fulaninho", phone: "011888888888", status: "Na fila", err: errors.New(ErrorArgumentNotValidAddConsumer)},
		{id: "FakeID", name: "Fulaninho", phone: "011888888888", status: "Na fila", err: errors.New(repository.ErrorNotFoundStore)},
	}

	for _, test := range tests {
		accessURL, err := svc.AddConsumer(test.id, test.name, test.phone, test.status)
		if err == nil {
			assert.NotNil(t, accessURL)
		} else {
			assert.Equal(t, test.err, err)
			assert.Empty(t, accessURL)
		}
	}

}

func TestRemoveConsumer(t *testing.T) {

	svc := NewStoreMockServiceImpl()

	name := "Outback"

	store, err := svc.Create(name)

	assert.Nil(t, err)
	assert.NotNil(t, store)

	consumerName := "Fulano"
	consumerPhone := "011998989898"
	consumerFakePhone := "011988888888"
	status := "Na fila"

	accessConsumerURL, err2 := svc.AddConsumer(store.ID, consumerName, consumerPhone, status)

	assert.Nil(t, err2)
	assert.NotNil(t, accessConsumerURL)

	tests := []struct {
		id    string
		phone string
		err   error
	}{
		{id: store.ID, phone: consumerPhone, err: nil},
		{id: store.ID, phone: consumerFakePhone, err: errors.New(repository.ErrorNotFoundConsumer)},
		{id: "fakeID", phone: consumerPhone, err: errors.New(repository.ErrorNotFoundConsumer)},
		{id: "", phone: consumerPhone, err: errors.New(ErrorArgumentNotValidRemoveConsumer)},
	}

	for _, test := range tests {
		err := svc.RemoveConsumer(test.id, test.phone)
		assert.Equal(t, test.err, err)
	}

}

func TestGetConsumer(t *testing.T) {

	svc := NewStoreMockServiceImpl()

	name := "Outback"

	store, err := svc.Create(name)

	assert.Nil(t, err)
	assert.NotNil(t, store)

	consumerName := "Fulano"
	consumerPhone := "011998989898"
	consumerFakePhone := "011988888888"
	status := "Na fila"

	accessConsumerURL, err2 := svc.AddConsumer(store.ID, consumerName, consumerPhone, status)

	assert.Nil(t, err2)
	assert.NotNil(t, accessConsumerURL)

	tests := []struct {
		id     string
		phone  string
		status string
		err    error
	}{
		{id: store.ID, phone: consumerPhone, status: "Na fila", err: nil},
		{id: store.ID, phone: consumerFakePhone, status: "Na fila", err: errors.New(repository.ErrorNotFoundConsumer)},
		{id: "fakeID", phone: consumerPhone, status: "Na fila", err: errors.New(repository.ErrorNotFoundConsumer)},
		{id: "", phone: consumerPhone, status: "Na fila", err: errors.New(ErrorArgumentNotValidGetConsumer)},
	}

	for _, test := range tests {
		position, consumer, err := svc.GetConsumer(test.id, test.phone)
		if err == nil {
			assert.NotNil(t, consumer)
			assert.NotEmpty(t, consumer.Name)
			assert.NotEmpty(t, consumer.Phone)
			assert.NotEmpty(t, consumer.Accesskey)
			assert.NotEqual(t, position, -1)
		} else {
			assert.Equal(t, test.err, err)
			assert.Nil(t, consumer)
		}
	}

}

func TestGetAllConsumers(t *testing.T) {

	svc := NewStoreMockServiceImpl()

	name := "Outback"

	store, err := svc.Create(name)

	assert.Nil(t, err)
	assert.NotNil(t, store)

	consumers := []struct {
		name   string
		phone  string
		status string
	}{
		{name: "Fulano Um", phone: "011998989899", status: "Na fila"},
		{name: "Fulano Dois", phone: "011976767676", status: "Na fila"},
		{name: "Fulano Tres", phone: "011954545454", status: "Na fila"},
	}

	for _, c := range consumers {
		accessConsumerURL, err := svc.AddConsumer(store.ID, c.name, c.phone, c.status)
		assert.Nil(t, err)
		assert.NotNil(t, accessConsumerURL)
	}

	result, err := svc.GetAllConsumers(store.ID)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	result, err2 := svc.GetAllConsumers("")

	assert.NotNil(t, err2)
	assert.Equal(t, err2, errors.New(ErrorArgumentNotValidGetConsumer))
	assert.Nil(t, result)

}
