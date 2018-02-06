package repos

import (
	"log"

	"github.com/josanr/OpenGts-Events/models"
)

type EventsRepository struct {
	store *DB
}

func NewEventsRepository(db *DB) *EventsRepository {
	var drepo = &EventsRepository{
		store: db,
	}

	return drepo
}

func (dr *EventsRepository) GetByDevice(device models.Device, timeStart, timeEnd string) (ev models.Events, err error) {
	request := `
	SELECT
	  latitude,
	  longitude
	FROM
		EventData
	WHERE
		deviceID = ?
		AND timestamp BETWEEN ? AND ?
	`

	rows, err := dr.store.Query(request, device.Id, timeStart, timeEnd)

	if err != nil {
		return ev, err
	}
	defer rows.Close()
	var tempEv = models.Events{
		"sysadmin",
		"Numina",
		"GMT+02:00",
		models.DeviceList{
			Device:      device.Id,
			Device_desc: device.Name,
			EventData:   []models.EventData{},
		},
	}

	for rows.Next() {

		var evData = models.EventData{}
		err := rows.Scan(&evData.GPSPoint_lat, &evData.GPSPoint_lon)
		if err != nil {
			return ev, err
		}

		tempEv.DeviceList.EventData = append(tempEv.DeviceList.EventData, evData)
	}
	err = rows.Err()
	if err != nil {
		return ev, err
	}
	log.Println(tempEv.DeviceList)
	if len(tempEv.DeviceList.EventData) == 0 {
		return ev, NewEmptyListError("Empty list.")
	}

	return tempEv, nil
}
