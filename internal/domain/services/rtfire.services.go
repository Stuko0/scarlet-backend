package services

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	repository "scarlet_backend/config"
	"scarlet_backend/internal/domain/entities"
)

var (
	rtfireRepo repository.RTFireRepository = repository.NewRTFireRepository()
)

// GetFireData obtiene los datos del incendio
func GetRTFireData() ([]*entities.RTFire, error) {
	url := "https://api.ambeedata.com/fire/latest/by-lat-lng?lat=-17.3895&lng=-66.1568"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-api-key", "c82c0cd946c6a4648f9fb19a32032f4952def0810665e4f6fcd372335c3a312f")
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

func SaveRTFireData() ([]*entities.RTFire, error) {
	// Obtiene los datos del incendio
	fires, err := GetRTFireData()
	if err != nil {
		return nil, err
	}

	// Obtiene los incendios existentes de la base de datos
	existingFires, err := rtfireRepo.GetRTFireFromDB()
	if err != nil {
		return nil, err
	}

	var savedRTFires []*entities.RTFire
	for _, fire := range fires {

		// Verifica si el incendio ya existe en la base de datos
		existing := false
		for _, existingFire := range existingFires {
			if fire.Latitude == existingFire.Latitude && fire.Longitude == existingFire.Longitude {
				existing = true
				fire.DocID = existingFire.DocID
				fire.Id = existingFire.Id
				break
			}
		}

		var savedRTFire *entities.RTFire
		if existing {
			savedRTFire, err = rtfireRepo.UpdateFire(fire)
		} else {
			fire.Id = rand.Int63()
			savedRTFire, err = rtfireRepo.SaveRTFire(fire)
		}

		if err != nil {
			return nil, err
		}
		savedRTFires = append(savedRTFires, savedRTFire)
	}

	return savedRTFires, nil
}

func SaveRTFire(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-type", "application/json")

	savedRTFires, err := SaveRTFireData()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error guardando los datos del incendio"}`))
		return
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(savedRTFires)
}

func UpdateRTFire(fire *entities.RTFire) (*entities.RTFire, error) {
	updatedFire, err := rtfireRepo.UpdateFire(fire)
	if err != nil {
		return nil, err
	}
	return updatedFire, nil
}

func GetSavedRTFires(resp http.ResponseWriter, req *http.Request) {
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
