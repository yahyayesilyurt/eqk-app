import "./App.css";
import Map from "./components/Map";
import { useState, useEffect } from "react";

const App = () => {
  const [earthquakeData, setEarthquakeData] = useState([]);

  useEffect(() => {
    fetchEarthquakeData();
  }, []);

  const fetchEarthquakeData = async () => {
    const response = await fetch("http://localhost:8080");
    const data = await response.json();
    // console.log(data.features);
    setEarthquakeData(data.features);
  };
  return (
    <div>
      <h1>Earthquake Tracker</h1>
      <Map center={[0, 0]} zoom={2} data={earthquakeData} />
    </div>
  );
};

export default App;
