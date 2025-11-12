package nhtsa

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	nhtsaBaseURL   = "https://vpic.nhtsa.dot.gov/api/vehicles"
	maxRetries     = 3
	initialTimeout = 30 * time.Second
)

// HttpClient dengan timeout lebih panjang
var httpClient = &http.Client{
	Timeout: initialTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 10,
	},
}

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

// FetchModelsForMakeID mengambil model untuk merek tertentu dari NHTSA dengan retry logic
func FetchModelsForMakeID(makeID string) ([]NhtsaModel, error) {
	url := fmt.Sprintf("%s/GetModelsForMakeId/%s?format=json", nhtsaBaseURL, makeID)

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Attempt %d/%d: Fetching models for MakeID %s...", attempt, maxRetries, makeID)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		// Set headers untuk menghindari rate limiting
		req.Header.Set("User-Agent", "CarApp/1.0")
		req.Header.Set("Accept", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			log.Printf("Attempt %d failed: %v", attempt, err)

			// Jangan retry jika ini attempt terakhir
			if attempt < maxRetries {
				// Exponential backoff: 2s, 4s, 8s
				backoff := time.Duration(attempt*2) * time.Second
				log.Printf("Retrying in %v...", backoff)
				time.Sleep(backoff)
				continue
			}
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("NHTSA API returned status: %s", resp.Status)
			log.Printf("Attempt %d: HTTP error %s", attempt, resp.Status)

			if attempt < maxRetries && resp.StatusCode >= 500 {
				// Retry hanya untuk server error (5xx)
				backoff := time.Duration(attempt*2) * time.Second
				time.Sleep(backoff)
				continue
			}
			break
		}

		var apiResponse NhtsaModelResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			lastErr = err
			log.Printf("Attempt %d: JSON decode error: %v", attempt, err)
			break
		}

		log.Printf("✓ Sukses mengambil %d model untuk MakeID %s dari NHTSA API", len(apiResponse.Results), makeID)
		return apiResponse.Results, nil
	}

	return nil, fmt.Errorf("failed after %d attempts: %v", maxRetries, lastErr)
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
