package tests

import (
	"fmt"
	"testing"

	"github.com/rokoga/filas-backend/service"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	svc := service.NewStoreServiceImpl()

	URLname := "outback"
	name := "Outback"

	store, err := svc.Create(URLname, name)

	assert.Nil(t, err)
	assert.NotNil(t, store)

	assert.NotEmpty(t, store.ID)
	assert.NotEmpty(t, store.Name)
	assert.NotEmpty(t, store.URLName)
	assert.Equal(t, name, store.Name)
	assert.Equal(t, URLname, store.URLName)

	store2, err := svc.Create("", "")
	assert.Nil(t, store2)
	assert.NotNil(t, err)

	fmt.Printf("Store: %v", store)
}
