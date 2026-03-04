// userApi.ts - Updated to manage authentication
import axios from 'axios';

// API base URL
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

// User model interfaces - Updated to match Go capitalization
export interface User {
  ID: string;
  Name: {
    First: string;
    Last: string;
  };
  Role: string;
  SbuID: string;
  Vehicle: string;
  Address?: {
    Street: string;
    City: string;
    State: string;
    ZipCode: string;
  };
  DriverLicense?: {
    Number: string;
    State: string;
    ExpirationDate: string;
  };
  Username: string;
  PasswordHash: string;
  // Keep lowercase for backward compatibility with our components
  name?: {
    first: string;
    last: string;
  };
}

// Login request and response interfaces
export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  userId: string;
  role: string;
  user: User;
}

// Registration interfaces
export interface RegisterRequest {
  username: string;
  passwordHash: string;
  name: {
    first: string;
    last: string;
  };
  sbuId: string;
  role?: string;
  address?: {
    street: string;
    city: string;
    state: string;
    zipCode: string;
  };
  driverLicense?: {
    number: string;
    state: string;
    expirationDate: string;
  };
}

export interface RegisterResponse {
  id: string;
}

// Authentication helpers
export const saveAuthCredentials = (username: string, password: string) => {
  localStorage.setItem('auth_username', username);
  localStorage.setItem('auth_password', password);
};

export const getAuthCredentials = () => {
  return {
    username: localStorage.getItem('auth_username') || '',
    password: localStorage.getItem('auth_password') || ''
  };
};

export const clearAuthCredentials = () => {
  localStorage.removeItem('auth_username');
  localStorage.removeItem('auth_password');
};

export const isAuthenticated = () => {
  const { username, password } = getAuthCredentials();
  return !!username && !!password;
};

// API client functions
const userApi = {
  // Login function
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    try {
      const response = await axios.post<LoginResponse>(
        `${API_URL}/login`, 
        credentials,
        {
          auth: {
            username: credentials.username,
            password: credentials.password
          }
        }
      );
      
      // Save credentials on successful login
      saveAuthCredentials(credentials.username, credentials.password);
      
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data || 'Login failed');
      }
      throw new Error('Login failed. Please try again.');
    }
  },

  // Logout function
  logout: () => {
    clearAuthCredentials();
  },

  // Register function (unchanged)
  register: async (userData: RegisterRequest): Promise<RegisterResponse> => {
    try {
      // Transform to match backend expectations - capitalize field names
      const transformedData = {
        Username: userData.username,
        PasswordHash: userData.passwordHash,
        Name: {
          First: userData.name.first,
          Last: userData.name.last
        },
        SbuID: userData.sbuId,
        Role: userData.role || 'user',
        Address: userData.address ? {
          Street: userData.address.street,
          City: userData.address.city,
          State: userData.address.state,
          ZipCode: userData.address.zipCode
        } : undefined,
        DriverLicense: userData.driverLicense ? {
          Number: userData.driverLicense.number,
          State: userData.driverLicense.state,
          ExpirationDate: userData.driverLicense.expirationDate
        } : undefined
      };

      const response = await axios.post<RegisterResponse>(
        `${API_URL}/register`, 
        transformedData
      );
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        if (error.response.status === 409) {
          throw new Error('Username already exists');
        }
        throw new Error(error.response.data || 'Registration failed');
      }
      throw new Error('Registration failed. Please try again.');
    }
  },

  // Get user profile - now uses stored credentials
  getUserProfile: async (userId: string): Promise<User> => {
    try {
      const { username, password } = getAuthCredentials();
      
      if (!username || !password) {
        throw new Error('Not authenticated');
      }

      const response = await axios.get<User>(
        `${API_URL}/users/${userId}`,
        { 
          auth: {
            username,
            password
          }
        }
      );
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data || 'Failed to get user profile');
      }
      throw new Error('Failed to get user profile. Please try again.');
    }
  },

  // Update user profile - now uses stored credentials
  updateUserProfile: async (userId: string, userData: Partial<User>): Promise<void> => {
    try {
      const { username, password } = getAuthCredentials();
      
      if (!username || !password) {
        throw new Error('Not authenticated');
      }

      await axios.put(
        `${API_URL}/users/${userId}`,
        userData,
        { 
          auth: {
            username,
            password
          }
        }
      );
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data || 'Failed to update user profile');
      }
      throw new Error('Failed to update user profile. Please try again.');
    }
  }
};

export default userApi;