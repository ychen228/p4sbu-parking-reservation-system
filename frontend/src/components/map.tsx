import React, { useState, useEffect } from "react";
import { APIProvider, Map, Marker } from "@vis.gl/react-google-maps";
import Search from "./search";
import { ParkingLot } from "../types/models";
import { User } from "../api/userApi";
import NavBar from "./top_nav";
import "./css_files/map.css";

const apiKey: string = process.env.REACT_APP_API_MAP_KEY!;
const defaultCenter = {
  lat: 40.91323221202016,
  lng: -73.12363213209412,
}; //Center of campus
const map_style = [
  {
    featureType: "administrative",
    elementType: "geometry",
    stylers: [
      {
        visibility: "off",
      },
    ],
  },
  {
    featureType: "poi",
    stylers: [
      {
        visibility: "off",
      },
    ],
  },
  {
    featureType: "road",
    elementType: "labels.icon",
    stylers: [
      {
        visibility: "off",
      },
    ],
  },
  {
    featureType: "transit",
    stylers: [
      {
        visibility: "off",
      },
    ],
  },
];
interface MapProps {
  user: User;
}

const Sbu_Map: React.FC<MapProps> = ({ user }) => {
  console.log(user);
  const [building_marker, set_building_marker] = useState<{
    lat: number;
    lng: number;
  } | null>(null);
  const [lotMarkers, set_lot_markers] = useState<ParkingLot[] | null>(null);
  const [mapCenter, setMapCenter] = useState(defaultCenter);
  const [current_lot_marker, set_current_lot_marker] = useState<{
    lat: number;
    lng: number;
  } | null>(null);
  // Add the new highlighted lot state
  const [highlighted_lot, set_highlighted_lot] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  return (
    <div className="page-container">
      <NavBar user={user} />
      <div className="map">
        <APIProvider apiKey={apiKey}>
          <Search
            set_building_marker={set_building_marker}
            set_current_lot_marker={set_current_lot_marker}
            setMapCenter={setMapCenter}
            set_lot_markers={set_lot_markers}
            set_highlighted_lot={set_highlighted_lot}
            user={user}
          ></Search>
          <Map
            className="google-map"
            defaultZoom={17}
            center={mapCenter}
            disableDefaultUI={true}
            onCenterChanged={(e) => {
              const newCenter = e.detail.center;
              setMapCenter(newCenter);
            }}
            styles={map_style}
            restriction={{
              latLngBounds: {
                north: 40.930126,
                south: 40.890999,
                west: -73.145663,
                east: -73.100342,
              },
              strictBounds: true,
            }}
          >
            {building_marker && (
              <Marker
                position={building_marker}
                icon={{
                  path: google.maps.SymbolPath.CIRCLE,
                  scale: 7,
                  fillColor: "#FF0000", // Red for building
                  fillOpacity: 1,
                  strokeWeight: 1,
                  strokeColor: "#FFFFFF",
                }}
              />
            )}

            {lotMarkers &&
              lotMarkers.map((lot) => (
                <Marker
                  icon={{
                    path: google.maps.SymbolPath.BACKWARD_CLOSED_ARROW,
                    // Only define scale once
                    scale: highlighted_lot === lot.ID ? 7 : 5,
                    fillColor:
                      highlighted_lot === lot.ID ? "#FF0000" : "#00BFFF",
                    fillOpacity: 1,
                    strokeWeight: 0,
                  }}
                  key={lot.ID}
                  position={{ lat: lot.Location.lat, lng: lot.Location.lng }}
                  onClick={() => {
                    // When clicking on a marker, highlight it and update the current lot marker
                    set_highlighted_lot(lot.ID);
                    set_current_lot_marker({
                      lat: lot.Location.lat,
                      lng: lot.Location.lng,
                    });
                  }}
                />
              ))}
          </Map>
        </APIProvider>
      </div>
    </div>
  );
};
export default Sbu_Map;
