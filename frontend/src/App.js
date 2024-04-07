import "./App.css";
import Map from "./components/Map";
import { useState, useEffect } from "react";
import "leaflet/dist/leaflet.css";
import LoadingSpinner from "./components/LoadingSpinner";

const App = () => {
  const [earthquakeData, setEarthquakeData] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    console.log("RENDER");
    fetchEarthquakeData();

    const interval = setInterval(() => {
      fetchEarthquakeData();
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  const fetchEarthquakeData = async () => {
    setLoading(true);
    try {
      const response = await fetch("http://127.0.0.1:8080/api");
      if (!response.ok) {
        throw new Error("Failed to fetch earthquake data");
      }
      const data = await response.json();
      setLoading(false);
      setEarthquakeData(data);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="container">
      <h1>Earthquake App</h1>
      {loading ? (
        <div className="loading">
          <LoadingSpinner />
        </div>
      ) : (
        <Map center={[0, 0]} zoom={2} data={earthquakeData} />
      )}
    </div>
  );
};

export default App;
