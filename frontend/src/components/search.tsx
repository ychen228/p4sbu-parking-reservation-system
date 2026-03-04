import React, { useState, useEffect, Dispatch, SetStateAction } from "react";
import { 
  getBuildings, 
  getParkingLots, 
  getNearestParkingLots,
  ParkingLotWithDistance
} from "../api/apiClient";
import { Building, ParkingLot } from "../types/models";
import { User } from "../api/userApi";
import "./css_files/search.css";
import ParkingLotInfo from "./Parking_Lot_Info";

interface MapProps {
  set_building_marker: Dispatch<SetStateAction<{lat: number, lng: number} | null>>;
  setMapCenter: Dispatch<SetStateAction<{lat: number, lng: number}>>;
  set_lot_markers: Dispatch<SetStateAction<ParkingLot[] | null>>;
  set_current_lot_marker: Dispatch<SetStateAction<{lat: number, lng: number} | null>>;
  set_highlighted_lot: Dispatch<SetStateAction<string | null>>;
  user: User;
}

const Search: React.FC<MapProps> = ({
  set_building_marker,
  set_current_lot_marker,
  setMapCenter,
  set_lot_markers,
  set_highlighted_lot,
  user
}) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [buildings, setBuildings] = useState<Building[]>([]);
  const [building_selected, set_selected_building] = useState<Building | null>(null);
  const [filtered_list, set_filtered] = useState<Building[]>([]);
  
  // Store all nearest parking lots in sorted order
  const [all_nearest_lots, set_all_nearest_lots] = useState<ParkingLotWithDistance[]>([]);
  
  // Get all available parking lots 
  const [all_parking_lots, set_all_parking_lots] = useState<ParkingLot[]>([]);

  // Track the currently selected parking lot
  const [selected_parking, set_selected_parking] = useState<ParkingLot | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetch_Buildings_ParkingLots = async () => {
      try {
        console.log("fetching buildings and parking lots");
        setLoading(true);
        const building_data = await getBuildings();
        setBuildings(building_data);
        
        // Fetch all parking lots so we have the complete dataset
        const parkingLot_data = await getParkingLots();
        set_all_parking_lots(parkingLot_data);
        
        setError(null);
      } catch (err) {
        setError("Failed to fetch data");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetch_Buildings_ParkingLots();
  }, []);

  useEffect(() => {
    if (searchQuery.trim().length > 0) {
      const filtered = buildings.filter((building) =>
        building.Name.toLowerCase().includes(searchQuery.trim().toLowerCase())
      );
      if (filtered.length == 1 && searchQuery.trim().toLowerCase() === filtered[0].Name.toLowerCase()){
        set_filtered([]);
      } else {
        set_selected_building(null);
        set_lot_markers(null);
        set_highlighted_lot(null);
        if (all_nearest_lots.length > 0) {
          set_all_nearest_lots([]);
        }
        set_filtered(filtered);
      }
    } else {
      set_filtered([]);
    }
  }, [searchQuery]);

  useEffect(() => {
    if (selected_parking) {
      // Update the highlighted marker
      set_highlighted_lot(selected_parking.ID);
      
      // Center map on selected parking lot
      set_current_lot_marker({ 
        lat: selected_parking.Location.lat, 
        lng: selected_parking.Location.lng 
      });
      
      setMapCenter({ 
        lat: selected_parking.Location.lat, 
        lng: selected_parking.Location.lng 
      });
    } else {
      // Clear the highlighted marker when no parking lot is selected
      set_highlighted_lot(null);
    }
  }, [selected_parking]);

  useEffect(() => {
    const calculateNearestLots = async () => {
      try {
        setLoading(true);
        if (building_selected && all_parking_lots.length > 0) {
          // Calculate distances client-side to sort all parking lots
          const withDistances = all_parking_lots.map(lot => {
            // Calculate Euclidean distance
            const dx = lot.Location.lng - building_selected.Location.lng;
            const dy = lot.Location.lat - building_selected.Location.lat;
            const distance = Math.sqrt(dx*dx + dy*dy);
            
            return {
              parking_lot: lot,
              distance_km: distance
            };
          });
          
          // Sort by distance
          withDistances.sort((a, b) => a.distance_km - b.distance_km);
          
          // Set all lots with distances
          set_all_nearest_lots(withDistances);
          
          // Only show top 5 parking lots on the map initially
          const topFiveLots = withDistances.slice(0, 5).map(item => item.parking_lot);
          
          // Set the building marker and center the map
          set_building_marker({ 
            lat: building_selected.Location.lat, 
            lng: building_selected.Location.lng 
          });
          
          setMapCenter({ 
            lat: building_selected.Location.lat, 
            lng: building_selected.Location.lng 
          });
          
          // Set the parking lot markers on the map
          set_lot_markers(topFiveLots);
          
          // Clear any previously selected lot
          set_selected_parking(null);
          set_highlighted_lot(null);
        }
        setError(null);
      } catch (err) {
        setError("Failed to calculate nearest parking lots");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    
    console.log("building is selected", building_selected);
    calculateNearestLots();
  }, [building_selected, all_parking_lots]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (filtered_list.length > 0) {
      getInfo(filtered_list[0]);
    }
  };

  const getInfo = (building: Building) => {
    setSearchQuery(building.Name);
    set_selected_building(building);
  };

  const handleParkingLotSelect = (parkingLotItem: ParkingLotWithDistance) => {
    const selectedLot = parkingLotItem.parking_lot;
    set_selected_parking(selectedLot);
    
    // Check if the selected lot is already in the displayed markers
    const topFiveIds = all_nearest_lots.slice(0, 5).map(item => item.parking_lot.ID);
    
    // If not in the top 5, update the map to include this specific lot
    if (!topFiveIds.includes(selectedLot.ID)) {
      // Get the first 4 lots plus the selected one
      const topFourLots = all_nearest_lots.slice(0, 4).map(item => item.parking_lot);
      set_lot_markers([...topFourLots, selectedLot]);
    }
  };

  return (
    <div className="search">
      <form onSubmit={handleSearch} className="search-form">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Enter a building"
          className="search-input"
        />
        <button type="submit" className="search-button">
          Search
        </button>
      </form>
      
      {/* Display filtered buildings when typing */}
      {filtered_list.length > 0 && (
        <div className="search-results-container">
          <h3>Buildings</h3>
          <ul className="search-results">
            {filtered_list.map((building) => (
              <li 
                className="search-info" 
                key={building.ID} 
                onClick={() => getInfo(building)}
              >
                {building.Name}
              </li>
            ))}
          </ul>
        </div>
      )}
      
      {/* Display all nearest parking lots when a building is selected */}
      {all_nearest_lots.length > 0 && (
        <div className="search-results-container">
          <h3>Nearest Parking Lots ({all_nearest_lots.length})</h3>
          <ul className="search-results parking-lots-list">
            {all_nearest_lots.map((item, index) => (
              <li 
                className={`search-info ${selected_parking && selected_parking.ID === item.parking_lot.ID ? 'selected' : ''}`}
                key={item.parking_lot.ID} 
                onClick={() => handleParkingLotSelect(item)}
              >
                <div className="lot-info">
                  <span className="lot-rank">{index + 1}</span>
                  <span className="lot-name">{item.parking_lot.Name}</span>
                </div>
                <span className="lot-distance">{item.distance_km.toFixed(4)} km</span>
              </li>
            ))}
          </ul>
        </div>
      )}
      
      {/* Display loading state */}
      {loading && <div className="loading">Loading...</div>}
      
      {/* Display error if any */}
      {error && <div className="error">{error}</div>}
      
      {/* Parking lot details */}
      <ParkingLotInfo 
        parkingLot={selected_parking} 
        set_selected_parking={set_selected_parking} 
        user={user}
      />
    </div>
  );
};

export default Search;