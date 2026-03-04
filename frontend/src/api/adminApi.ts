import axios from 'axios';
import { ParkingLot, Building, Node, Reservation, Violation, Vehicle } from '../types/models';
import { User, getAuthCredentials, isAuthenticated } from './userApi';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance without hardcoded credentials
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
      // Handle authentication error - could redirect to login page
      console.error('Authentication failed. Please log in again.');
      // If you have a router, you could redirect: window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const getAllUsers = async (): Promise<User[]> => {
  try {
    const response = await api.get('/admin/users');
    console.log(response.data)
    return response.data;
  } catch (error) {
    console.error('Failed to fetch users:', error);
    return [];
  }
};

// Rest of your API functions remain the same
export const deleteUser = async (id: string): Promise<void> => {
  try {
    await api.delete(`/admin/users/${id}`);
  } catch (error) {
    console.error(`Failed to delete user with ID ${id}:`, error);
  }
};

export const createParkingLot = async (parkingLot: ParkingLot): Promise<void> => {
  try {
    await api.post('/admin/parking-lots', parkingLot);
  } catch (error) {
    console.error('Failed to create parking lot:', error);
  }
};

export const updateParkingLot = async (id: string, parkingLot: ParkingLot): Promise<void> => {
  try {
    await api.put(`/admin/parking-lots/${id}`, parkingLot);
  } catch (error) {
    console.error(`Failed to update parking lot with ID ${id}:`, error);
  }
};

export const deleteParkingLot = async (id: string): Promise<void> => {
  try {
    await api.delete(`/admin/parking-lots/${id}`);
  } catch (error) {
    console.error(`Failed to delete parking lot with ID ${id}:`, error);
  }
};
export const createBuilding = async (building: Building): Promise<void> => {
  try {
    await api.post('/admin/buildings', building);
  } catch (error) {
    console.error('Failed to create building:', error);
  }
};

export const updateBuilding = async (id: string, building: Building): Promise<void> => {
  try {
    await api.put(`/admin/buildings/${id}`, building);
  } catch (error) {
    console.error(`Failed to update building with ID ${id}:`, error);
  }
};

export const deleteBuilding = async (id: string): Promise<void> => {
  try {
    await api.delete(`/admin/buildings/${id}`);
  } catch (error) {
    console.error(`Failed to delete building with ID ${id}:`, error);
  }
};

export const createNode = async (node: Node): Promise<void> => {
  try {
    await api.post('/admin/nodes', node);
  } catch (error) {
    console.error('Failed to create node:', error);
  }
};

export const updateNode = async (id: string, node: Node): Promise<void> => {
  try {
    await api.put(`/admin/nodes/${id}`, node);
  } catch (error) {
    console.error(`Failed to update node with ID ${id}:`, error);
  }
};

export const deleteNode = async (id: string): Promise<void> => {
  try {
    await api.delete(`/admin/nodes/${id}`);
  } catch (error) {
    console.error(`Failed to delete node with ID ${id}:`, error);
  }
};
export const getAllReservations = async (): Promise<Reservation[]> => {
  try {
    const response = await api.get('/admin/reservations');
    return response.data;
  } catch (error) {
    console.error('Failed to fetch reservations:', error);
    return [];
  }
};
export const getReservationsByParkingLot = async (parkingLotId: string): Promise<Reservation[]> => {
  try {
    const response = await api.get(`/admin/parking-lots/${parkingLotId}/reservations`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch reservations for parking lot ID ${parkingLotId}:`, error);
    return [];
  }
};
export default {
  getAllUsers,
  deleteUser,
  createParkingLot,
  updateParkingLot,
  deleteParkingLot,
  createBuilding,
  updateBuilding,
  deleteBuilding,
  createNode,
  updateNode,
  deleteNode,
  getAllReservations,
  getReservationsByParkingLot,
};
