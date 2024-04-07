package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EarthquakeData struct {
	Magnitude float64 `json:"magnitude"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
	uri := os.Getenv("MONGODB_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
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
	http.HandleFunc("/api", getDataHandler)
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
    var data []EarthquakeData

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"magnitude": bson.M{"$gte":4}})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var earthquake EarthquakeData
        if err := cursor.Decode(&earthquake); err != nil {
            return nil, err
        }
        data = append(data, earthquake)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, feature := range responseData.Features {
		if feature.Properties.Mag >= 4{
			earthquake := EarthquakeData{
				Magnitude: feature.Properties.Mag,
				Longitude: feature.Geometry.Coordinates[0],
				Latitude:  feature.Geometry.Coordinates[1],
			}
			_, err := collection.InsertOne(ctx, earthquake)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}


// Script to add earthquake data manually
func addDataManually() {
	for {
		earthquake := getEarthquakeFromUser()
		err := addData(earthquake)
		if err != nil {
			log.Printf("Error while adding data: %v\n", err)
		}

		time.Sleep(1 * time.Minute)

	}
}

// Script to add random earthquake data
func addDataRandomly() {
	for {
		earthquake := generateRandomEarthquake()
		err := addData(earthquake)
		if err != nil {
			log.Printf("Error while adding data: %v\n", err)
		}

		time.Sleep(1 * time.Minute)
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
	}
}

func addData(data EarthquakeData) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if data.Magnitude >= 4 {
        _, err := collection.InsertOne(ctx, data)
        if err != nil {
            return err
        }
        log.Printf("New data added: %v\n", data)
    } else {
        log.Printf("Earthquake magnitude is not greater than 4, skipping insertion: %v\n", data)
    }

    return nil
}
