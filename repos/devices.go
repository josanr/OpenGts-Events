package repos

import (
	"database/sql"

	"../models"
)

type DeviceRepository struct {
	store *DB
}

func NewDeviceRepository(db *DB) *DeviceRepository {
	var drepo = &DeviceRepository{
		store: db,
	}

	return drepo
}

func (dr DeviceRepository) processRows(rows *Rows) (list []models.Device, err error) {

	for rows.Next() {
		var tempDev = models.Device{}

		err := rows.Scan(&tempDev.Id, &tempDev.Imei, &tempDev.Name)
		if err != nil {
			return list, err
		}

		list = append(list, tempDev)
	}
	err = rows.Err()
	if err != nil {
		return list, err
	}
	if len(list) == 0 {
		return list, NewEmptyListError("Empty list.")
	}
	return list, nil
}

func (dr *DeviceRepository) GetAll() (list []models.Device, err error) {
	request := `select
		deviceID,
		imeiNumber,
		description
	FROM Device
	WHERE
	1;`

	rows, err := dr.store.Query(request)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	list, err = dr.processRows(&Rows{rows})
	if err != nil {
		return list, err
	}

	return list, err
}

func (dr *DeviceRepository) GetByName(name string) (list []models.Device, err error) {
	request := `select
		deviceID,
		imeiNumber,
		description
	FROM Device
	WHERE
	description LIKE ?;`

	rows, err := dr.store.Query(request, name)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	list, err = dr.processRows(&Rows{rows})
	if err != nil {
		return list, err
	}

	return list, err
}

func (dr *DeviceRepository) GetById(id string) (device models.Device, err error) {
	request := `select
		deviceID,
		imeiNumber,
		description
	FROM Device
	WHERE
	deviceID = ?;`

	rows := dr.store.QueryRow(request, id)

	device = models.Device{}
	err = rows.Scan(&device.Id, &device.Imei, &device.Name)
	if err != nil && err != sql.ErrNoRows {
		return device, NewEmptyListError("Empty list.")
	}

	return device, nil
}
