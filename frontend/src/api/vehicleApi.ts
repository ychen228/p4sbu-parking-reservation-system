import axios from 'axios';
import { Vehicle } from '../types/models';
import { getAuthCredentials, isAuthenticated } from './userApi';

// API base URL
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Add request interceptor for authentication
api.interceptors.request.use(config => {
  // Get credentials from localStorage
  if (isAuthenticated()) {
    const { username, password } = getAuthCredentials();
    config.auth = {
      username,
      password
    };
  }
  
  return config;
}, error => {
  return Promise.reject(error);
});

// Add response interceptor to handle auth errors
api.interceptors.response.use(
  response => response,
  error => {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      console.error('Authentication failed. Please log in again.');
      // Could redirect: window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Get a single vehicle by ID
export const getVehicle = async (id: string): Promise<Vehicle> => {
  try {
    const response = await api.get(`/vehicles/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch vehicle with ID ${id}:`, error);
    throw error;
  }
};

// Get all vehicles for a specific user
export const getUserVehicles = async (userId: string): Promise<Vehicle[]> => {
  try {
    const response = await api.get(`/users/${userId}/vehicles`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch vehicles for user ID ${userId}:`, error);
    return [];
  }
};

// Create a new vehicle
export const createVehicle = async (vehicle: Omit<Vehicle, 'ID'>): Promise<{ id: string }> => {
  try {
    console.log("here")
    const response = await api.post('/vehicles', vehicle);
    return response.data;
  } catch (error) {
    console.error('Failed to create vehicle:', error);
    throw error;
  }
};

// Update an existing vehicle
export const updateVehicle = async (id: string, vehicleData: Partial<Vehicle>): Promise<{ message: string }> => {
  try {
    const response = await api.put(`/vehicles/${id}`, vehicleData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update vehicle with ID ${id}:`, error);
    throw error;
  }
};

// Delete a vehicle
export const deleteVehicle = async (id: string): Promise<{ message: string }> => {
  try {
    const response = await api.delete(`/vehicles/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to delete vehicle with ID ${id}:`, error);
    throw error;
  }
};

const vehicleApi = {
  getVehicle,
  getUserVehicles,
  createVehicle,
  updateVehicle,
  deleteVehicle,
};

export default vehicleApi;