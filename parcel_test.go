package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEqual(t, 0, id)

	parc, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, id, parc.Number)
	require.Equal(t, parcel.Client, parc.Client)
	require.Equal(t, parcel.Status, parc.Status)
	require.Equal(t, parcel.Address, parc.Address)
	require.Equal(t, parcel.CreatedAt, parc.CreatedAt)

	err = store.Delete(id)
	require.NoError(t, err)
	_, err = store.Get(id)
	require.Equal(t, sql.ErrNoRows, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEqual(t, 0, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	parc, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, parc.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEqual(t, 0, id)

	status := ParcelStatusSent
	err = store.SetStatus(id, status)
	require.NoError(t, err)

	parc, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, status, parc.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
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
		require.NoError(t, err)
		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {
		require.Equal(t, parcelMap[parcel.Number], parcel)
		require.Equal(t, parcelMap[parcel.Number].Client, parcel.Client)
		require.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		require.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
		require.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)
	}
}
