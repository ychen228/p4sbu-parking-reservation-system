import axios from 'axios';
import { Reservation } from '../types/models';

// Base URL for API
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true  // Required for session-based auth
});

// Fetch all reservations (admin-only)
export const getAllReservations = async (): Promise<Reservation[]> => {
  const response = await api.get('/reservations');
  return response.data;
};

// Fetch a single reservation by ID
export const getReservation = async (id: string): Promise<Reservation> => {
  const response = await api.get(`/reservations/${id}`);
  return response.data;
};

// Create a new reservation
export const createReservation = async (
  reservation: Omit<Reservation, 'ID' | 'Status'>
): Promise<Reservation> => {
  const response = await api.post('/reservations', reservation);
  return response.data;
};

// Cancel (delete) a reservation
export const cancelReservation = async (id: string): Promise<{ message: string }> => {
  const response = await api.delete(`/reservations/${id}`);
  return response.data;
};

// Get all reservations for a specific user
export const getUserReservations = async (userId: string): Promise<Reservation[]> => {
  const response = await api.get(`/users/${userId}/reservations`);
  return response.data;
};

// Update an existing reservation
export const updateReservation = async (
  id: string,
  updatedData: Partial<Omit<Reservation, 'ID' | 'ReservedBy'>>
): Promise<{ message: string }> => {
  const response = await api.put(`/reservations/${id}`, updatedData);
  return response.data;
};

const reservationApi = {
  getAllReservations,
  getReservation,
  createReservation,
  cancelReservation,
  getUserReservations,
  updateReservation,
};

export default reservationApi;

