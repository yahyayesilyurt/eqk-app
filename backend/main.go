package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EarthquakeData struct {
	Magnitude float64   `json:"magnitude"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Coords    []float64 `json:"coordinates"`
}

var (
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
)

func main() {
	connectToMongoDB()
	go startDataCollection()
	go startServer()
	go addDataManually()
	go addDataRandomly()

	select {}
}

func connectToMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("earthquakes")
	collection = db.Collection("data")
	log.Println("Connected to MongoDB")
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
	getData()
}

func getDataHandler(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are accepted.", http.StatusMethodNotAllowed)
		return
	}

	data, err := getDataFromMongoDB()
	if err != nil {
		http.Error(w, "Failed to retrieve data.", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getDataFromMongoDB() ([]EarthquakeData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var data []EarthquakeData
	for cursor.Next(ctx) {
		var eqData EarthquakeData
		if err := cursor.Decode(&eqData); err != nil {
			return nil, err
		}
		data = append(data, eqData)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return data, nil
}


func getData() {
	response, err := http.Get("https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var responseData struct {
		Features []struct {
			Properties struct {
				Mag float64 `json:"mag"`
			} `json:"properties"`
			Geometry struct {
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}

	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		log.Fatal(err)
	}

	for _, feature := range responseData.Features {
		earthquake := EarthquakeData{
			Magnitude: feature.Properties.Mag,
			Coords:    feature.Geometry.Coordinates,
			Longitude: feature.Geometry.Coordinates[0],
			Latitude:  feature.Geometry.Coordinates[1],
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := collection.InsertOne(ctx, earthquake)
		if err != nil {
			log.Fatal(err)
		}
	}
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		http.Error(w, "Error while adding data to MongoDB.", http.StatusInternalServerError)
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
		Magnitude: mag,
		Latitude:  lat,
		Longitude: lon,
		Coords:    []float64{lon, lat},
	}
}

func generateRandomEarthquake() EarthquakeData {
	mag := rand.Float64() * 10
	lat := rand.Float64() * 180 - 90
	lon := rand.Float64() * 360 - 180

	return EarthquakeData{
		Magnitude: mag,
		Latitude:  lat,
		Longitude: lon,
		Coords:    []float64{lon, lat},
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
