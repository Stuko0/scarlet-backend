package services

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	repository "scarlet_backend/config"
	"scarlet_backend/internal/domain/entities"
	"strconv"
	"strings"
)

var(fireRepo repository.FireRepository = repository.NewFireRepository())

// GetFireData obtiene los datos del incendio
func GetFireData() ([]*entities.Fire, error) {
    url := "https://firms.modaps.eosdis.nasa.gov/api/country/csv/529a9508e24d53be1007c992621172cb/MODIS_SP/BOL/10/2024-01-10"
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error haciendo la solicitud: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error leyendo la respuesta: %v", err)
        return nil, err
    }

    r := csv.NewReader(strings.NewReader(string(body)))
    records, err := r.ReadAll()
    if err != nil {
        log.Printf("Error leyendo el CSV: %v", err)
        return nil, err
    }

    var fires []*entities.Fire
    for i, record := range records {
        // Ignora el primer registro porque contiene los encabezados del CSV
        if i == 0 {
            continue
        }

        latitude, err := strconv.ParseFloat(record[1], 64)
        if err != nil {
            log.Printf("Error convirtiendo la latitud a float64: %v", err)
            return nil, err
        }
        longitude, err := strconv.ParseFloat(record[2], 64)
        if err != nil {
            log.Printf("Error convirtiendo la longitud a float64: %v", err)
            return nil, err
        }

        fire := &entities.Fire{
            Latitude:  latitude,
            Longitude: longitude,
            AcqDate:   record[6],
            AcqTime:   record[7],
            DayNight:  record[14],
        }
        fires = append(fires, fire)
    }

    return fires, nil
}

func GetFires(resp http.ResponseWriter, req *http.Request) {
    fire, err := GetFireData()
    if err != nil {
        http.Error(resp, err.Error(), http.StatusInternalServerError)
        return
    }

    fireJson, err := json.Marshal(fire)
    if err != nil {
        http.Error(resp, err.Error(), http.StatusInternalServerError)
        return
    }

    resp.Write(fireJson)
}


func SaveFire(resp http.ResponseWriter, req *http.Request){
    resp.Header().Set("Content-type", "application/json")

    // Obtiene los datos del incendio
    fires, err := GetFireData()
    if err != nil {
        resp.WriteHeader(http.StatusInternalServerError)
        resp.Write([]byte(`{"error:" "Error obteniendo los datos del incendio"}`))
        return
    }

    var savedFires []*entities.Fire
    for _, fire := range fires {
        fire.Id=rand.Int63()
        savedFire, err := fireRepo.SaveFire(fire)
        if err != nil {
            resp.WriteHeader(http.StatusInternalServerError)
            resp.Write([]byte(`{"error:" "Error guardando los datos del incendio"}`))
            return
        }
        savedFires = append(savedFires, savedFire)
    }

    resp.WriteHeader(http.StatusOK)
    json.NewEncoder(resp).Encode(savedFires)
}

func GetSavedFires(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")

	fires, err := fireRepo.GetFireFromDB()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error obteniendo los datos del incendio de la base de datos"}`))
		return
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(fires)
}
