import { useState, useEffect } from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { Icon } from "leaflet";
import MarkerClusterGroup from "react-leaflet-cluster";

const Map = ({ center, zoom, data }) => {
  const [markers, setMarkers] = useState(data);

  const customIcon = new Icon({
    iconUrl: "/mapicon.png",
    iconSize: [38, 38],
  });

  useEffect(() => {
    const timer = setTimeout(() => {
      setMarkers([]);
    }, 28000);

    return () => clearTimeout(timer);
  }, [data]);

  return (
    <MapContainer
      center={center}
      zoom={zoom}
      style={{ height: "%80", width: "100%" }}
    >
      <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
      <MarkerClusterGroup chunkedLoading>
        {markers.map((earthquake, index) => (
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
