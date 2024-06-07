package services

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	repository "scarlet_backend/config"
	"scarlet_backend/internal/domain/entities"
)

var(rtfireRepo repository.RTFireRepository = repository.NewRTFireRepository())

// GetFireData obtiene los datos del incendio
func GetRTFireData() ([]*entities.RTFire, error) {
	url := "https://api.ambeedata.com/fire/latest/by-lat-lng?lat=-17.3895&lng=-66.1568"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-api-key", "d73ba263d4f19b2d7c2293e4e55b1e5f408792b0482f51d89e32fccec9f69d3a")
	req.Header.Add("Content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	data := result["data"].([]interface{})
	var fires []*entities.RTFire
	for _, item := range data {
		fireData := item.(map[string]interface{})
		fire := &entities.RTFire{
			Latitude:   fireData["lat"].(float64),
			Longitude:  fireData["lng"].(float64),
			DetectedAt: fireData["detectedAt"].(string),
			Confidence: fireData["confidence"].(string),
			FRP:        fireData["frp"].(float64),
			FWI:        fireData["fwi"].(float64),
			FireType:   fireData["fireType"].(string),
		}
		fires = append(fires, fire)
	}

	return fires, nil
}


func GetRTFires(resp http.ResponseWriter, req *http.Request) {
    fire, err := GetRTFireData()
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


func SaveRTFire(resp http.ResponseWriter, req *http.Request){
    resp.Header().Set("Content-type", "application/json")

    // Obtiene los datos del incendio
    fires, err := GetRTFireData()
    if err != nil {
        resp.WriteHeader(http.StatusInternalServerError)
        resp.Write([]byte(`{"error:" "Error obteniendo los datos del incendio"}`))
        return
    }

    var savedRTFires []*entities.RTFire
    for _, fire := range fires {
        fire.Id=rand.Int63()
        savedRTFire, err := rtfireRepo.SaveRTFire(fire)
        if err != nil {
            resp.WriteHeader(http.StatusInternalServerError)
            resp.Write([]byte(`{"error:" "Error guardando los datos del incendio"}`))
            return
        }
        savedRTFires = append(savedRTFires, savedRTFire)
    }

    resp.WriteHeader(http.StatusOK)
    json.NewEncoder(resp).Encode(savedRTFires)
}

func GetSavedRTFires(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	fires, err := rtfireRepo.GetRTFireFromDB()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error obteniendo los datos del incendio de la base de datos"}`))
		return
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(fires)
}

func DeleteAllRTFires(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-type", "application/json")

	err := rtfireRepo.DeleteAllRTFires()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error eliminando los datos del incendio de la base de datos"}`))
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`{"message:" "Todos los incendios han sido eliminados exitosamente"}`))
}
