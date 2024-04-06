import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import "leaflet/dist/leaflet.css";

const Map = ({ center, zoom, data }) => {
  console.log(data);
  return (
    <MapContainer
      center={center}
      zoom={zoom}
      style={{ height: "500px", width: "100%" }}
    >
      <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
      {data.map((earthquake, index) => (
        <Marker key={index} position={[earthquake.lat, earthquake.lon]}>
          <Popup>
            <div>
              <h3>Magnitude: {earthquake.mag}</h3>
              <p>Latitude: {earthquake.lat}</p>
              <p>Longitude: {earthquake.lon}</p>
            </div>
          </Popup>
        </Marker>
      ))}
    </MapContainer>
  );
};

export default Map;
