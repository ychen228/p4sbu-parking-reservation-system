import axios from 'axios';
import { ViolationRebute } from '../types/models';
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

// Get all rebutes (admin only)
export const getAllRebutes = async (): Promise<ViolationRebute[]> => {
  try {
    const response = await api.get('/admin/violation-rebutes');
    return response.data;
  } catch (error) {
    console.error('Failed to fetch rebutes:', error);
    return [];
  }
};

// Get a single rebute by ID
export const getRebute = async (id: string): Promise<ViolationRebute> => {
  try {
    const response = await api.get(`/violation-rebutes/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch rebute with ID ${id}:`, error);
    throw error;
  }
};

// Get all rebutes for a specific user
export const getUserRebutes = async (userId: string): Promise<ViolationRebute[]> => {
  try {
    const response = await api.get(`/users/${userId}/violation-rebutes`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch rebutes for user ID ${userId}:`, error);
    return [];
  }
};

// Get all rebutes for a specific violation
export const getViolationRebutes = async (violationId: string): Promise<ViolationRebute[]> => {
  try {
    const response = await api.get(`/violations/${violationId}/rebutes`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch rebutes for violation ID ${violationId}:`, error);
    return [];
  }
};

// Create a new rebute
export const createRebute = async (rebute: Omit<ViolationRebute, 'ID'>): Promise<{ id: string }> => {
  try {
    const response = await api.post('/violation-rebutes', rebute);
    return response.data;
  } catch (error) {
    console.error('Failed to create rebute:', error);
    throw error;
  }
};

// Update an existing rebute (admin only)
export const updateRebute = async (id: string, rebuteData: Partial<ViolationRebute>): Promise<{ message: string }> => {
  try {
    const response = await api.put(`/admin/violation-rebutes/${id}`, rebuteData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update rebute with ID ${id}:`, error);
    throw error;
  }
};

// Delete a rebute (admin only)
export const deleteRebute = async (id: string): Promise<{ message: string }> => {
  try {
    const response = await api.delete(`/admin/violation-rebutes/${id}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to delete rebute with ID ${id}:`, error);
    throw error;
  }
};

const violationRebuteApi = {
  getAllRebutes,
  getRebute,
  getUserRebutes,
  getViolationRebutes,
  createRebute,
  updateRebute,
  deleteRebute,
};

export default violationRebuteApi;