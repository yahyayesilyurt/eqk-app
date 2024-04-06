import "./App.css";
import Map from "./components/Map";
import { useState, useEffect } from "react";

const App = () => {
  const [earthquakeData, setEarthquakeData] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchEarthquakeData();
  }, []);

  const fetchEarthquakeData = async () => {
    setLoading(true);
    const response = await fetch("http://localhost:8080/api");
    const data = await response.json();
    setLoading(false);
    console.log(data);
    setEarthquakeData(data);
  };

  if (loading) {
    return <div>Loading...</div>;
  } else {
    return (
      <div>
        <h1>Earthquake App</h1>
        <Map center={[0, 0]} zoom={2} data={earthquakeData} />
      </div>
    );
  }
};

export default App;
