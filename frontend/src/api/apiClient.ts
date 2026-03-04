import axios from 'axios';
import { ParkingLot, Building, Node, User } from '../types/models';

// Base URL for API
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// API functions for parking lots
export const getParkingLots = async (): Promise<ParkingLot[]> => {
  try {
    const response = await api.get<ParkingLot[]>('/parking-lots');
    return response.data;
  } catch (error) {
    console.error('Error fetching parking lots:', error);
    return [];
  }
};

export const getParkingLotByBuilding = async (building: Building): Promise<ParkingLot[]> => {
  try {
    const response = await api.get<ParkingLot[]>(`/parking-lots/${building}`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching parking lot near ${building.Name}:`, error);
    return [];
  }
};

// Function to get nearest parking lots to a building
export interface ParkingLotWithDistance {
  parking_lot: ParkingLot;
  distance_km: number;
}

export const getNearestParkingLots = async (buildingId: string, limit?: number): Promise<ParkingLotWithDistance[]> => {
  try {
    // Add limit parameter only if provided
    const url = limit 
      ? `/wayfind/building/${buildingId}/nearest-parking?limit=${limit}`
      : `/wayfind/building/${buildingId}/nearest-parking`;
      
    const response = await api.get<ParkingLotWithDistance[]>(url);
    return response.data;
  } catch (error) {
    console.error(`Error fetching nearest parking lots for building ${buildingId}:`, error);
    return [];
  }
};

// Alternative method using coordinates
export const getNearestParkingLotsByCoordinates = async (
  lat: number, 
  lng: number, 
  limit?: number
): Promise<ParkingLotWithDistance[]> => {
  try {
    // Build URL based on whether limit is provided
    const url = limit
      ? `/wayfind/nearest-parking?lat=${lat}&lng=${lng}&limit=${limit}`
      : `/wayfind/nearest-parking?lat=${lat}&lng=${lng}`;
      
    const response = await api.get<ParkingLotWithDistance[]>(url);
    return response.data;
  } catch (error) {
    console.error(`Error fetching nearest parking lots for coordinates (${lat}, ${lng}):`, error);
    return [];
  }
};

export const getParkingLotById = async (id: string): Promise<ParkingLot | null> => {
  try {
    const response = await api.get<ParkingLot>(`/parking-lots/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching parking lot with ID ${id}:`, error);
    return null;
  }
};

// API functions for buildings
export const getBuildings = async (): Promise<Building[]> => {
  try {
    const response = await api.get<Building[]>('/buildings');
    return response.data;
  } catch (error) {
    console.error('Error fetching buildings:', error);
    return [];
  }
};

// API functions for nodes
export const getNodes = async (): Promise<Node[]> => {
  try {
    const response = await api.get<Node[]>('/nodes');
    return response.data;
  } catch (error) {
    console.error('Error fetching nodes:', error);
    return [];
  }
};

export default {
  getParkingLots,
  getParkingLotById,
  getBuildings,
  getNodes,
  getNearestParkingLots,
  getNearestParkingLotsByCoordinates,
};