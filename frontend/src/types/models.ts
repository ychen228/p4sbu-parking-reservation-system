export interface Location {
  lat: number;
  lng: number;
}

export interface ParkingLot {
  ID: string;
  Name: string;
  Spaces: number;
  Faculty: number;
  Premium: number;
  Metered: number;
  Resident: number;
  Ada: number;
  Ev: boolean;
  Active: boolean;
  Location: Location;
  Node: string;
  Fee: number;
  Occupancy: number;
}

export interface Building {
  ID: string;
  Name: string;
  Location: Location;
  Node: string;
}

export interface Node {
  ID: string;
  Type: string;
  ParkingLot?: string;
  Building?: string;
  Neighbors: {
    Members: string[];
  };
  Location: Location;
}

export interface User {
  ID: string;
  Name: {
    First: string;
    Last: string;
  };
  Role: string;
  SbuID: string;
  Vehicle: string;
  Address: {
    Street: string;
    City: string;
    State: string;
    ZipCode: string;
  };
  DriverLicense: {
    Number: string;
    State: string;
    ExpirationDate: string;
  };
  Username?: string;
}

export interface Vehicle {
  ID: string;
  Name: string;
  Model: string;
  Year: number;
  PlateNumber: string;
}

export interface Reservation {
  ID: string;
  ReservedBy: string;
  ParkingLot: string;
  StartTime: string;
  EndTime: string;
  Status: string;
}

// Adding missing models based on the backend schema
export interface Violation {
  ID: string;
  User: string;
  ParkingLot: string;
  Reason: string;
  Fine: number;
  PayBy: string;
}

export interface ViolationRebute {
  ID: string;
  User: string;
  ParkingLot: string;
  Violation: string;
  Reason: string;
}

// Auth-related request/response interfaces
export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  userId: string;
  role: string;
  user: User;
  expiresIn: number;
}

export interface RegisterRequest {
  Name: {
    First: string;
    Last: string;
  };
  Role?: string;
  SbuID: string;
  Address: {
    Street: string;
    City: string;
    State: string;
    ZipCode: string;
  };
  DriverLicense: {
    Number: string;
    State: string;
    ExpirationDate: Date;
  };
  Username: string;
  PasswordHash: string; // Will be hashed on server
}

// Generic response for created resources
export interface CreateResponse {
  ID: string;
  message?: string;
}
