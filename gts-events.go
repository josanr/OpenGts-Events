package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./models"

	"./repos"
	"os"
	"path"
	"path/filepath"
)

var db *repos.DB
var err error

type Configuration struct {
	Host     string `json:"host"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Db       string `json:"db"`
	Port     string `json:"port"`
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	filename := filepath.Dir(ex)
	filePath := path.Join(path.Dir(filename), "/opengts-events/config.json")
	file, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Panic("error:", err)
	}

	db, err = repos.NewStore(conf.Login, conf.Password, conf.Db, conf.Host, conf.Port)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization")

		w.Write([]byte("OKOK"))
		return
	}
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization")
		w.Write([]byte("OKOK"))
		return
	}
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
