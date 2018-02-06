package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/josanr/OpenGts-Events/models"

	"github.com/josanr/OpenGts-Events/repos"
)

var db *repos.DB
var err error

func main() {
	db, err = repos.NewStore()
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/events/", queryEvents)
	http.HandleFunc("/devices/", queryDevices)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	err = http.ListenAndServe(":1603", nil)
	if err != nil {
		log.Panic(err)
	}
}

func queryEvents(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var deviceStr = r.Form.Get("d")
	var timeStart = r.Form.Get("rf")
	var timeEnd = r.Form.Get("rt")
	var evRepo = repos.NewEventsRepository(db)
	var devRepo = repos.NewDeviceRepository(db)

	dev, err := devRepo.GetById(deviceStr)
	if err != nil {
		log.Println("device not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	events, err := evRepo.GetByDevice(dev, timeStart, timeEnd)
	if err != nil {
		log.Println("events for device not found")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func queryDevices(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var queryString = r.Form.Get("query")
	var list []models.Device

	devs := repos.NewDeviceRepository(db)

	if queryString == "" {
		list, err = devs.GetAll()
	} else {
		list, err = devs.GetByName(queryString)
	}

	if err != nil {
		switch err.(type) {
		case repos.EmptyListError:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	res, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
