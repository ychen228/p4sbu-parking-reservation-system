import React, { useEffect, useState } from "react";
import {
  getAllUsers,
  deleteUser,
  getAllReservations,
  updateParkingLot,
  deleteParkingLot,
  updateBuilding,
  deleteBuilding,
  createParkingLot,
  createBuilding,
} from "../api/adminApi";
import { getBuildings, getParkingLots } from "../api/apiClient";
import { User, isAuthenticated } from "../api/userApi";
import { ParkingLot, Building, Reservation } from "../types/models";
import NavBar from "./top_nav";
import "./css_files/adminEdit.css";

interface AdminProps {
  user: User | null;
}

// Modal component for creating/editing items
interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
}

const Modal: React.FC<ModalProps> = ({ isOpen, onClose, title, children }) => {
  if (!isOpen) return null;

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>{title}</h3>
          <button className="btn btn-secondary" onClick={onClose}>
            ✕
          </button>
        </div>
        {children}
      </div>
    </div>
  );
};

const AdminEdit: React.FC<AdminProps> = ({ user }) => {
  const [users, setUsers] = useState<User[]>([]);
  const [parkingLots, setParkingLots] = useState<ParkingLot[]>([]);
  const [buildings, setBuildings] = useState<Building[]>([]);
  const [reservations, setReservations] = useState<Reservation[]>([]);
  
  // State for modals
  const [isParkingLotModalOpen, setIsParkingLotModalOpen] = useState(false);
  const [isBuildingModalOpen, setIsBuildingModalOpen] = useState(false);
  
  // State for edit/create objects
  const [editingParkingLot, setEditingParkingLot] = useState<ParkingLot | null>(null);
  const [editingBuilding, setEditingBuilding] = useState<Building | null>(null);
  
  // Status messages
  const [statusMessage, setStatusMessage] = useState<{ text: string; type: 'success' | 'error' } | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // Check authentication on component mount
    if (!isAuthenticated()) {
      setStatusMessage({ 
        text: "You need to be logged in as an admin to access this page", 
        type: "error" 
      });
      return;
    }
    
    fetchData();
  }, []);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      await Promise.all([
        fetchUsers(),
        fetchReservations(),
        fetchBuilding(),
        fetchParkingLot(),
      ]);
    } catch (error) {
      console.error("Error fetching data:", error);
      setStatusMessage({ 
        text: "Failed to load data. Please try again.", 
        type: "error" 
      });
    } finally {
      setIsLoading(false);
    }
  };

  const fetchUsers = async () => {
    try {
      const data = await getAllUsers();
      console.log("Fetched users:", data);
      // Make sure we're properly handling the user data structure
      console.log(data)
      if (Array.isArray(data)) {
        setUsers(data as User[]);
      } else if (data && typeof data === 'object') {
        // If it's not an array but an object, check if it has expected user properties
        const typedData = data as Record<string, any>;
        if (typedData.ID || typedData.id) {
          setUsers([typedData as User]); // Single user as array
        } else {
          // It might be an object containing users array
          const usersArray = Object.values(typedData).filter(item => 
            item && typeof item === 'object' && 
            ((item as any).ID != null || (item as any).Name != null || (item as any).Username != null)
          );
          if (usersArray.length > 0) {
            setUsers(usersArray as User[]);
          } else {
            console.error("Unexpected user data structure:", data);
            setUsers([]);
          }
        }
      } else {
        setUsers([]);
      }
    } catch (error) {
      console.error("Error fetching users:", error);
      setUsers([]); // Set empty array instead of throwing
    }
  };

  const fetchReservations = async () => {
    try {
      const data = await getAllReservations();
      console.log("Fetched reservations:", data);
      // Handle different data structures
      if (Array.isArray(data)) {
        setReservations(data as Reservation[]);
      } else if (data && typeof data === 'object') {
        // If it's not an array but an object with reservation properties
        const typedData = data as Record<string, any>;
        if (typedData.ID || typedData.StartTime) {
          setReservations([typedData as Reservation]); // Single reservation as array
        } else {
          // It might be an object containing reservations array
          const reservationsArray = Object.values(typedData).filter(item => 
            item && typeof item === 'object' && 
            ((item as any).ID != null || (item as any).StartTime != null || (item as any).ReservedBy != null)
          );
          if (reservationsArray.length > 0) {
            setReservations(reservationsArray as Reservation[]);
          } else {
            console.error("Unexpected reservation data structure:", data);
            setReservations([]);
          }
        }
      } else {
        setReservations([]);
      }
    } catch (error) {
      console.error("Error fetching reservations:", error);
      setReservations([]); // Set empty array instead of throwing
    }
  };

  const fetchBuilding = async () => {
    try {
      const data = await getBuildings();
      setBuildings(data);
    } catch (error) {
      console.error("Error fetching buildings:", error);
      throw error;
    }
  };

  const fetchParkingLot = async () => {
    try {
      const data = await getParkingLots();
      setParkingLots(data);
    } catch (error) {
      console.error("Error fetching parking lots:", error);
      throw error;
    }
  };

  const handleDeleteUser = async (id: string) => {
    try {
      await deleteUser(id);
      setUsers(users.filter((user) => user.ID !== id));
      setStatusMessage({ text: "User deleted successfully", type: "success" });
    } catch (error) {
      console.error(`Error deleting user with ID ${id}:`, error);
      setStatusMessage({ text: "Failed to delete user", type: "error" });
    }
  };

  // Parking Lot Functions
  const handleCreateParkingLot = async () => {
    if (!editingParkingLot) return;
    
    try {
      await createParkingLot(editingParkingLot);
      await fetchParkingLot(); // Refetch to get the ID from server
      setIsParkingLotModalOpen(false);
      setEditingParkingLot(null);
      setStatusMessage({ text: "Parking lot created successfully", type: "success" });
    } catch (error) {
      console.error("Error creating parking lot:", error);
      setStatusMessage({ text: "Failed to create parking lot", type: "error" });
    }
  };

  const handleUpdateParkingLot = async () => {
    if (!editingParkingLot || !editingParkingLot.ID) return;
    
    try {
      await updateParkingLot(editingParkingLot.ID, editingParkingLot);
      setParkingLots((prev) =>
        prev.map((pl) =>
          pl.ID === editingParkingLot.ID ? editingParkingLot : pl
        )
      );
      setIsParkingLotModalOpen(false);
      setEditingParkingLot(null);
      setStatusMessage({ text: "Parking lot updated successfully", type: "success" });
    } catch (error) {
      console.error("Error updating parking lot:", error);
      setStatusMessage({ text: "Failed to update parking lot", type: "error" });
    }
  };

  const handleDeleteParkingLot = async (id: string) => {
    try {
      await deleteParkingLot(id);
      setParkingLots((prev) => prev.filter((pl) => pl.ID !== id));
      setStatusMessage({ text: "Parking lot deleted successfully", type: "success" });
    } catch (error) {
      console.error(`Error deleting parking lot ${id}:`, error);
      setStatusMessage({ text: "Failed to delete parking lot", type: "error" });
    }
  };

  // Building Functions
  const handleCreateBuilding = async () => {
    if (!editingBuilding) return;
    
    try {
      await createBuilding(editingBuilding);
      await fetchBuilding(); // Refetch to get the ID from server
      setIsBuildingModalOpen(false);
      setEditingBuilding(null);
      setStatusMessage({ text: "Building created successfully", type: "success" });
    } catch (error) {
      console.error("Error creating building:", error);
      setStatusMessage({ text: "Failed to create building", type: "error" });
    }
  };

  const handleUpdateBuilding = async () => {
    if (!editingBuilding || !editingBuilding.ID) return;
    
    try {
      await updateBuilding(editingBuilding.ID, editingBuilding);
      setBuildings((prev) =>
        prev.map((b) => (b.ID === editingBuilding.ID ? editingBuilding : b))
      );
      setIsBuildingModalOpen(false);
      setEditingBuilding(null);
      setStatusMessage({ text: "Building updated successfully", type: "success" });
    } catch (error) {
      console.error("Error updating building:", error);
      setStatusMessage({ text: "Failed to update building", type: "error" });
    }
  };

  const handleDeleteBuilding = async (id: string) => {
    try {
      await deleteBuilding(id);
      setBuildings((prev) => prev.filter((b) => b.ID !== id));
      setStatusMessage({ text: "Building deleted successfully", type: "success" });
    } catch (error) {
      console.error(`Error deleting building ${id}:`, error);
      setStatusMessage({ text: "Failed to delete building", type: "error" });
    }
  };

  // Helper function to open parking lot modal
  const openParkingLotModal = (parkingLot?: ParkingLot) => {
    if (parkingLot) {
      setEditingParkingLot({ ...parkingLot });
    } else {
      // Create a new parking lot with default values
      setEditingParkingLot({
        ID: "",
        Name: "",
        Spaces: 0,
        Faculty: 0,
        Premium: 0,
        Metered: 0,
        Resident: 0,
        Ada: 0,
        Ev: false,
        Active: true,
        Location: { lat: 0, lng: 0 },
        Node: "",
        Fee: 0,
        Occupancy: 0,
      });
    }
    setIsParkingLotModalOpen(true);
  };

  // Helper function to open building modal
  const openBuildingModal = (building?: Building) => {
    if (building) {
      setEditingBuilding({ ...building });
    } else {
      // Create a new building with default values
      setEditingBuilding({
        ID: "",
        Name: "",
        Location: { lat: 0, lng: 0 },
        Node: "",
      });
    }
    setIsBuildingModalOpen(true);
  };

  // Render loading state
  if (isLoading) {
    return (
      <>
        <NavBar user={user} />
        <div className="admin-panel">
          <h1>Admin Panel</h1>
          <div className="loading-indicator">Loading data...</div>
        </div>
      </>
    );
  }

  // Debug information for troubleshooting
  console.log("Current state:", {
    users: users,
    parkingLots: parkingLots,
    buildings: buildings,
    reservations: reservations
  });

  return (
    <div className="page-container">
      <NavBar user={user} />
      <div className="admin-panel">
        <h1>Admin Panel</h1>

        {statusMessage && (
          <div className={`status-message status-${statusMessage.type}`}>
            {statusMessage.text}
          </div>
        )}
        
        {/* Users Section */}
        <section className="admin-section">
          <div className="section-header">
            <h2>Users</h2>
          </div>
          <div className="section-body">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Username</th>
                  <th>Role</th>
                  <th>SBU ID</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {users && users.length > 0 ? (
                  users.map((user) => (
                    <tr key={user.ID || `user-${Math.random()}`}>
                      <td>
                        {user.Name && typeof user.Name === 'object' 
                          ? `${user.Name.First || ''} ${user.Name.Last || ''}` 
                          : (user.name && typeof user.name === 'object'
                              ? `${user.name.first || ''} ${user.name.last || ''}` 
                              : (user.Name || 'Unknown'))}
                      </td>
                      <td>{user.Username || 'Unknown'}</td>
                      <td>{user.Role || 'Unknown'}</td>
                      <td>{user.SbuID || 'Unknown'}</td>
                      <td>
                        <button 
                          className="btn btn-danger" 
                          onClick={() => handleDeleteUser(user.ID || '')}
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={5} style={{ textAlign: "center" }}>
                      No users found or unable to display user data
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </section>

        {/* Parking Lots Section */}
        <section className="admin-section">
          <div className="section-header">
            <h2>Parking Lots</h2>
            <button 
              className="btn btn-primary" 
              onClick={() => openParkingLotModal()}
            >
              Add New Parking Lot
            </button>
          </div>
          <div className="section-body">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Total Spaces</th>
                  <th>Occupancy</th>
                  <th>Fee</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {parkingLots.map((pl) => (
                  <tr key={pl.ID}>
                    <td>{pl.Name}</td>
                    <td>{pl.Spaces}</td>
                    <td>{pl.Occupancy}</td>
                    <td>${pl.Fee.toFixed(2)}</td>
                    <td>{pl.Active ? "Active" : "Inactive"}</td>
                    <td>
                      <div className="btn-group">
                        <button 
                          className="btn btn-primary" 
                          onClick={() => openParkingLotModal(pl)}
                        >
                          Edit
                        </button>
                        <button 
                          className="btn btn-danger" 
                          onClick={() => handleDeleteParkingLot(pl.ID)}
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Buildings Section */}
        <section className="admin-section">
          <div className="section-header">
            <h2>Buildings</h2>
            <button 
              className="btn btn-primary" 
              onClick={() => openBuildingModal()}
            >
              Add New Building
            </button>
          </div>
          <div className="section-body">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Location</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {buildings.map((b) => (
                  <tr key={b.ID}>
                    <td>{b.Name}</td>
                    <td>
                      {b.Location?.lat?.toFixed(6) || 0}, {b.Location?.lng?.toFixed(6) || 0}
                    </td>
                    <td>
                      <div className="btn-group">
                        <button 
                          className="btn btn-primary" 
                          onClick={() => openBuildingModal(b)}
                        >
                          Edit
                        </button>
                        <button 
                          className="btn btn-danger" 
                          onClick={() => handleDeleteBuilding(b.ID)}
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Reservations Section */}
        <section className="admin-section">
          <div className="section-header">
            <h2>Reservations</h2>
          </div>
          <div className="section-body">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Reserved By</th>
                  <th>Parking Lot</th>
                  <th>Start Time</th>
                  <th>End Time</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {reservations && reservations.length > 0 ? (
                  reservations.map((res) => (
                    <tr key={res.ID || `res-${Math.random()}`}>
                      <td>
                        {typeof res.ReservedBy === 'object' 
                          ? JSON.stringify(res.ReservedBy).substring(0, 20) 
                          : (res.ReservedBy || 'Unknown')}
                      </td>
                      <td>
                        {typeof res.ParkingLot === 'object' 
                          ? JSON.stringify(res.ParkingLot).substring(0, 20) 
                          : (res.ParkingLot || 'Unknown')}
                      </td>
                      <td>
                        {res.StartTime 
                          ? new Date(res.StartTime).toLocaleString()
                          : 'Unknown'}
                      </td>
                      <td>
                        {res.EndTime 
                          ? new Date(res.EndTime).toLocaleString() 
                          : 'Unknown'}
                      </td>
                      <td>{res.Status || 'Unknown'}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={5} style={{ textAlign: "center" }}>
                      No reservations found or unable to display reservation data
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </section>

        {/* Modal for Parking Lot Edit/Create */}
        <Modal
          isOpen={isParkingLotModalOpen}
          onClose={() => {
            setIsParkingLotModalOpen(false);
            setEditingParkingLot(null);
          }}
          title={editingParkingLot?.ID ? "Edit Parking Lot" : "Add New Parking Lot"}
        >
          <div className="modal-body">
            <div className="admin-form">
              <div className="form-group">
                <label htmlFor="parking-name">Name</label>
                <input
                  id="parking-name"
                  type="text"
                  className="form-control"
                  value={editingParkingLot?.Name || ""}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Name: e.target.value } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-spaces">Total Spaces</label>
                <input
                  id="parking-spaces"
                  type="number"
                  className="form-control"
                  value={editingParkingLot?.Spaces || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Spaces: parseInt(e.target.value) || 0 } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-faculty">Faculty Spaces</label>
                <input
                  id="parking-faculty"
                  type="number"
                  className="form-control"
                  value={editingParkingLot?.Faculty || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Faculty: parseInt(e.target.value) || 0 } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-premium">Premium Spaces</label>
                <input
                  id="parking-premium"
                  type="number"
                  className="form-control"
                  value={editingParkingLot?.Premium || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Premium: parseInt(e.target.value) || 0 } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-metered">Metered Spaces</label>
                <input
                  id="parking-metered"
                  type="number"
                  className="form-control"
                  value={editingParkingLot?.Metered || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Metered: parseInt(e.target.value) || 0 } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-resident">Resident Spaces</label>
                <input
                  id="parking-resident"
                  type="number"
                  className="form-control"
                  value={editingParkingLot?.Resident || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Resident: parseInt(e.target.value) || 0 } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-ada">ADA Spaces</label>
                <input
                  id="parking-ada"
                  type="number"
                  className="form-control"
                  value={editingParkingLot?.Ada || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Ada: parseInt(e.target.value) || 0 } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-fee">Fee ($)</label>
                <input
                  id="parking-fee"
                  type="number"
                  step="0.01"
                  className="form-control"
                  value={editingParkingLot?.Fee || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Fee: parseFloat(e.target.value) || 0 } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-lat">Latitude</label>
                <input
                  id="parking-lat"
                  type="number"
                  step="0.000001"
                  className="form-control"
                  value={editingParkingLot?.Location?.lat || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { 
                        ...prev, 
                        Location: {
                          ...prev.Location,
                          lat: parseFloat(e.target.value) || 0
                        } 
                      } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="parking-lng">Longitude</label>
                <input
                  id="parking-lng"
                  type="number"
                  step="0.000001"
                  className="form-control"
                  value={editingParkingLot?.Location?.lng || 0}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { 
                        ...prev, 
                        Location: {
                          ...prev.Location,
                          lng: parseFloat(e.target.value) || 0
                        } 
                      } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group" style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                <label htmlFor="parking-ev">EV Charging Available</label>
                <input
                  id="parking-ev"
                  type="checkbox"
                  checked={editingParkingLot?.Ev || false}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Ev: e.target.checked } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group" style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                <label htmlFor="parking-active">Active</label>
                <input
                  id="parking-active"
                  type="checkbox"
                  checked={editingParkingLot?.Active || false}
                  onChange={(e) =>
                    setEditingParkingLot(prev => 
                      prev ? { ...prev, Active: e.target.checked } : null
                    )
                  }
                />
              </div>
            </div>
          </div>
          <div className="modal-footer">
            <button
              className="btn btn-secondary"
              onClick={() => {
                setIsParkingLotModalOpen(false);
                setEditingParkingLot(null);
              }}
            >
              Cancel
            </button>
            <button
              className="btn btn-primary"
              onClick={editingParkingLot?.ID ? handleUpdateParkingLot : handleCreateParkingLot}
            >
              {editingParkingLot?.ID ? "Save Changes" : "Create"}
            </button>
          </div>
        </Modal>

        {/* Modal for Building Edit/Create */}
        <Modal
          isOpen={isBuildingModalOpen}
          onClose={() => {
            setIsBuildingModalOpen(false);
            setEditingBuilding(null);
          }}
          title={editingBuilding?.ID ? "Edit Building" : "Add New Building"}
        >
          <div className="modal-body">
            <div className="admin-form">
              <div className="form-group">
                <label htmlFor="building-name">Name</label>
                <input
                  id="building-name"
                  type="text"
                  className="form-control"
                  value={editingBuilding?.Name || ""}
                  onChange={(e) =>
                    setEditingBuilding(prev => 
                      prev ? { ...prev, Name: e.target.value } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="building-lat">Latitude</label>
                <input
                  id="building-lat"
                  type="number"
                  step="0.000001"
                  className="form-control"
                  value={editingBuilding?.Location?.lat || 0}
                  onChange={(e) =>
                    setEditingBuilding(prev => 
                      prev ? { 
                        ...prev, 
                        Location: {
                          ...prev.Location,
                          lat: parseFloat(e.target.value) || 0
                        } 
                      } : null
                    )
                  }
                />
              </div>
              
              <div className="form-group">
                <label htmlFor="building-lng">Longitude</label>
                <input
                  id="building-lng"
                  type="number"
                  step="0.000001"
                  className="form-control"
                  value={editingBuilding?.Location?.lng || 0}
                  onChange={(e) =>
                    setEditingBuilding(prev => 
                      prev ? { 
                        ...prev, 
                        Location: {
                          ...prev.Location,
                          lng: parseFloat(e.target.value) || 0
                        } 
                      } : null
                    )
                  }
                />
              </div>
            </div>
          </div>
          <div className="modal-footer">
            <button
              className="btn btn-secondary"
              onClick={() => {
                setIsBuildingModalOpen(false);
                setEditingBuilding(null);
              }}
            >
              Cancel
            </button>
            <button
              className="btn btn-primary"
              onClick={editingBuilding?.ID ? handleUpdateBuilding : handleCreateBuilding}
            >
              {editingBuilding?.ID ? "Save Changes" : "Create"}
            </button>
          </div>
        </Modal>
      </div>
    </div>
  );
};

export default AdminEdit;