package nhtsa

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	nhtsaBaseURL = "https://vpic.nhtsa.dot.gov/api/vehicles"
)

// HttpClient untuk dependensi
var httpClient = &http.Client{Timeout: 10 * time.Second}

// Struct untuk parsing JSON response 'Get All Makes'
type NhtsaMake struct {
	MakeID   int    `json:"Make_ID"`
	MakeName string `json:"Make_Name"`
}
type NhtsaMakeResponse struct {
	Results []NhtsaMake `json:"Results"`
}

// Struct untuk parsing JSON response 'Get Models for Make ID'
type NhtsaModel struct {
	ModelID   int    `json:"Model_ID"`
	ModelName string `json:"Model_Name"`
	MakeID    int    `json:"Make_ID"`
	MakeName  string `json:"Make_Name"`
}
type NhtsaModelResponse struct {
	Results []NhtsaModel `json:"Results"`
}

// FetchAllMakes mengambil semua merek mobil dari NHTSA
func FetchAllMakes() ([]NhtsaMake, error) {
	url := fmt.Sprintf("%s/getallmakes?format=json", nhtsaBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NHTSA API returned status: %s", resp.Status)
	}

	var apiResponse NhtsaMakeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	log.Printf("Sukses mengambil %d merek dari NHTSA API", len(apiResponse.Results))
	return apiResponse.Results, nil
}

// FetchModelsForMakeID mengambil model untuk merek tertentu dari NHTSA
func FetchModelsForMakeID(makeID string) ([]NhtsaModel, error) {
	url := fmt.Sprintf("%s/GetModelsForMakeId/%s?format=json", nhtsaBaseURL, makeID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NHTSA API returned status: %s", resp.Status)
	}

	var apiResponse NhtsaModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	log.Printf("Sukses mengambil %d model untuk MakeID %s dari NHTSA API", len(apiResponse.Results), makeID)
	return apiResponse.Results, nil
}

// PENJELASAN FILE nhtsa_client.go:
// File ini menangani komunikasi dengan NHTSA API eksternal
//
// Constant:
// - nhtsaBaseURL: Base URL NHTSA API (https://vpic.nhtsa.dot.gov/api/vehicles)
// - httpClient: HTTP client dengan timeout 10 detik
//
// Struct NhtsaMake & NhtsaModel:
// - Untuk parsing JSON response dari API
// - MakeID/ModelID: ID numerik dari NHTSA
// - MakeName/ModelName: Nama merek/model mobil
//
// Fungsi FetchAllMakes:
// - Request GET ke /getallmakes?format=json
// - Return semua merek mobil yang ada di database NHTSA
// - Parse JSON response ke slice []NhtsaMake
//
// Fungsi FetchModelsForMakeID:
// - Request GET ke /GetModelsForMakeId/{makeID}?format=json
// - Return semua model untuk merek tertentu
// - Parse JSON response ke slice []NhtsaModel
//
// Error handling:
// - Cek HTTP status code (harus 200 OK)
// - Decode JSON response
// - Log jumlah data yang berhasil diambil

// PENJELASAN FILE nhtsa_client.go:
// File ini menangani komunikasi dengan NHTSA API eksternal
//
// Constant:
// - nhtsaBaseURL: Base URL NHTSA API (https://vpic.nhtsa.dot.gov/api/vehicles)
// - httpClient: HTTP client dengan timeout 10 detik
//
// Struct NhtsaMake & NhtsaModel:
// - Untuk parsing JSON response dari API
// - MakeID/ModelID: ID numerik dari NHTSA
// - MakeName/ModelName: Nama merek/model mobil
//
// Fungsi FetchAllMakes:
// - Request GET ke /getallmakes?format=json
// - Return semua merek mobil yang ada di database NHTSA
// - Parse JSON response ke slice []NhtsaMake
//
// Fungsi FetchModelsForMakeID:
// - Request GET ke /GetModelsForMakeId/{makeID}?format=json
// - Return semua model untuk merek tertentu
// - Parse JSON response ke slice []NhtsaModel
//
// Error handling:
// - Cek HTTP status code (harus 200 OK)
// - Decode JSON response
// - Log jumlah data yang berhasil diambil

// PENJELASAN FILE nhtsa_client.go:
// File ini menangani komunikasi dengan NHTSA API eksternal
//
// Constant:
// - nhtsaBaseURL: Base URL NHTSA API (https://vpic.nhtsa.dot.gov/api/vehicles)
// - httpClient: HTTP client dengan timeout 10 detik
//
// Struct NhtsaMake & NhtsaModel:
// - Untuk parsing JSON response dari API
// - MakeID/ModelID: ID numerik dari NHTSA
// - MakeName/ModelName: Nama merek/model mobil
//
// Fungsi FetchAllMakes:
// - Request GET ke /getallmakes?format=json
// - Return semua merek mobil yang ada di database NHTSA
// - Parse JSON response ke slice []NhtsaMake
//
// Fungsi FetchModelsForMakeID:
// - Request GET ke /GetModelsForMakeId/{makeID}?format=json
// - Return semua model untuk merek tertentu
// - Parse JSON response ke slice []NhtsaModel
//
// Error handling:
// - Cek HTTP status code (harus 200 OK)
// - Decode JSON response
// - Log jumlah data yang berhasil diambil
