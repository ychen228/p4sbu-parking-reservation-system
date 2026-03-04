import axios from 'axios';
import { Violation } from '../types/models';
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

// Get all violations (admin only)
export const getAllViolations = async (): Promise<Violation[]> => {
  try {
    const response = await api.get('/admin/violations');
    return response.data;
  } catch (error) {
    console.error('Failed to fetch violations:', error);
    return [];
  }
};

// Get a single violation by ID
export const getViolation = async (id: string): Promise<Violation> => {
  try {
    const response = await api.get(`/violations/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch violation with ID ${id}:`, error);
    throw error;
  }
};

// Get all violations for a specific user
export const getUserViolations = async (userId: string): Promise<Violation[]> => {
  try {
    const response = await api.get(`/users/${userId}/violations`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch violations for user ID ${userId}:`, error);
    return [];
  }
};

// Create a new violation (admin only)
export const createViolation = async (violation: Omit<Violation, 'ID'>): Promise<{ id: string }> => {
  try {
    const response = await api.post('/admin/violations', violation);
    return response.data;
  } catch (error) {
    console.error('Failed to create violation:', error);
    throw error;
  }
};

// Update an existing violation (admin only)
export const updateViolation = async (id: string, violationData: Partial<Violation>): Promise<{ message: string }> => {
  try {
    const response = await api.put(`/admin/violations/${id}`, violationData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update violation with ID ${id}:`, error);
    throw error;
  }
};

// Delete a violation (admin only)
export const deleteViolation = async (id: string): Promise<{ message: string }> => {
  try {
    const response = await api.delete(`/admin/violations/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to delete violation with ID ${id}:`, error);
    throw error;
  }
};

const violationApi = {
  getAllViolations,
  getViolation,
  getUserViolations,
  createViolation,
  updateViolation,
  deleteViolation,
};

export default violationApi;