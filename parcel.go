package main

import (
	"database/sql"
	"errors"
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
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return []Parcel{}, err
	}

	res := make([]Parcel, 1)

	for rows.Next() {
		var num, client int
		var status, address, createdAt string

		err := rows.Scan(&num, &client, &status, &address, &createdAt)
		if err != nil {
			return []Parcel{}, err
		}
		defer rows.Close()

		p := Parcel{
			Number:    num,
			Client:    client,
			Address:   address,
			Status:    status,
			CreatedAt: createdAt,
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status != ParcelStatusRegistered {
		return errors.New("status can be changed only for registered parcel")
	}

	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address), sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	return nil
}
