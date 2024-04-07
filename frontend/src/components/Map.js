import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { Icon } from "leaflet";
import MarkerClusterGroup from "react-leaflet-cluster";

const Map = ({ center, zoom, data }) => {
  const customIcon = new Icon({
    iconUrl: "/mapicon.png",
    iconSize: [38, 38],
  });

  return (
    <MapContainer
      center={center}
      zoom={zoom}
      style={{ height: "500px", width: "100%" }}
    >
      <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
      <MarkerClusterGroup chunkedLoading>
        {data.map((earthquake, index) => (
          <Marker
            key={index}
            position={[earthquake.latitude, earthquake.longitude]}
            icon={customIcon}
          >
            <Popup>
              <div>
                <h3>Magnitude: {earthquake.magnitude}</h3>
                <p>Latitude: {earthquake.latitude}</p>
                <p>Longitude: {earthquake.longitude}</p>
              </div>
            </Popup>
          </Marker>
        ))}
      </MarkerClusterGroup>
    </MapContainer>
  );
};

export default Map;
