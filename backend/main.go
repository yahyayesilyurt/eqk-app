package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type EarthquakeData struct {
	Features []struct {
		Properties struct {
			Mag float64 `json:"mag"`
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"properties"`
	} `json:"features"`
}

func main() {
	go startDataCollection()
	go startServer()
	go addDataManually()
	go addDataRandomly()

	select {}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func startServer() {
	http.HandleFunc("/", getDataHandler)
	http.HandleFunc("/add", addDataHandler)
	log.Println("HTTP server started via port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startDataCollection() {
	for {
		getData()
		time.Sleep(2 * time.Minute)
	}
}

func getDataHandler(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are accepted.", http.StatusMethodNotAllowed)
		return
	}

	data, err := getData()
	if err != nil {
		http.Error(w, "Failed to retrieve data.", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}


func getData() (EarthquakeData, error) {
	response, err := http.Get("https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson")
	if err != nil {
		return EarthquakeData{}, err
	}
	defer response.Body.Close()

	var data EarthquakeData
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return EarthquakeData{}, err
	}

	return data, nil
}

func addDataHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are accepted.", http.StatusMethodNotAllowed)
		return
	}

	var data EarthquakeData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error reading incoming data.", http.StatusBadRequest)
		return
	}

	log.Printf("New data added: %v\n", data)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data added successfully."))
}


// Script to add earthquake data manually
func addDataManually() {
	for {
		earthquake := getEarthquakeFromUser()
		err := addData(earthquake)
		if err != nil {
			log.Printf("Error while adding data: %v\n", err)
		} else {
			log.Printf("Manually data added successfully: %v\n", earthquake)
		}

		time.Sleep(2 * time.Minute)
	}
}

// Script to add random earthquake data
func addDataRandomly() {
	for {
		earthquake := generateRandomEarthquake()
		err := addData(earthquake)
		if err != nil {
			log.Printf("Error while adding data: %v\n", err)
		} else {
			log.Printf("Random data added successfully: %v\n", earthquake)
		}

		time.Sleep(2 * time.Minute)
	}
}

func getEarthquakeFromUser() EarthquakeData {
	var mag, lat, lon float64

	fmt.Println("Please enter the magnitude of the earthquake:")
	fmt.Scanln(&mag)

	fmt.Println("Please enter the latitude of the earthquake:")
	fmt.Scanln(&lat)

	fmt.Println("Please enter the longitude value of the earthquake:")
	fmt.Scanln(&lon)

	return EarthquakeData{
		Features: []struct {
			Properties struct {
				Mag float64 `json:"mag"`
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"properties"`
		}{
			{
				Properties: struct {
					Mag float64 `json:"mag"`
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				}{
					Mag: mag,
					Lat: lat,
					Lon: lon,
				},
			},
		},
	}
}

func generateRandomEarthquake() EarthquakeData {
	mag := rand.Float64() * 10
	lat := rand.Float64() * 180 - 90
	lon := rand.Float64() * 360 - 180

	return EarthquakeData{
		Features: []struct {
			Properties struct {
				Mag float64 `json:"mag"`
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"properties"`
		}{
			{
				Properties: struct {
					Mag float64 `json:"mag"`
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				}{
					Mag: mag,
					Lat: lat,
					Lon: lon,
				},
			},
		},
	}
}

func addData(data EarthquakeData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/add", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed response received from server: %s", resp.Status)
	}

	return nil
}
