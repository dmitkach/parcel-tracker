package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)

	parc, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, id, parc.Number)
	assert.Equal(t, parcel.Client, parc.Client)
	assert.Equal(t, parcel.Status, parc.Status)
	assert.Equal(t, parcel.Address, parc.Address)
	assert.Equal(t, parcel.CreatedAt, parc.CreatedAt)

	err = store.Delete(id)
	assert.NoError(t, err)
	_, err = store.Get(id)
	assert.ErrorIs(t, sql.ErrNoRows, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	assert.NoError(t, err)

	parc, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newAddress, parc.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)

	status := ParcelStatusSent
	err = store.SetStatus(id, status)
	assert.NoError(t, err)

	parc, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, status, parc.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		assert.NoError(t, err)
		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	assert.NoError(t, err)
	assert.Len(t, parcels, len(storedParcels))

	for _, parcel := range storedParcels {
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
