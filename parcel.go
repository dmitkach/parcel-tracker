package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(
		"INSERT INTO parcel (number, client, status, address, created_at) VALUES (:number, :client, :status, :address, :createdAt)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	resInt, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(resInt), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	var num, client int
	var status, address, createdAt string
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&num, &client, &status, &address, &createdAt)
	if err != nil {
		return Parcel{}, err
	}

	p := Parcel{
		Number:    num,
		Client:    client,
		Address:   address,
		Status:    status,
		CreatedAt: createdAt,
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	return nil
}
