import React, { useState, useEffect } from "react";
import { User } from "../api/userApi";
import userApi from "../api/userApi";
import { useNavigate } from "react-router-dom";
import reservationApi from "../api/reservationApi";
import { getParkingLots } from "../api/apiClient";
import vehicleApi from "../api/vehicleApi";
import Datetime from "react-datetime";
import "react-datetime/css/react-datetime.css";
import { Reservation, ParkingLot, Vehicle} from "../types/models";
import NavBar from "./top_nav";
import "./css_files/profile.css";

interface ProfileProps {
  user: User;
  set_user: (user: User | null) => void;
}

const Profile: React.FC<ProfileProps> = ({ user, set_user }) => {
  const navigate = useNavigate();
  const [editingSection, setEditingSection] = useState<string | null>(null);
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [confirmCancelId, setConfirmCancelId] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editedStartTime, setEditedStartTime] = useState<Date | null>(null);
  const [editedEndTime, setEditedEndTime] = useState<Date | null>(null);
  const [parkingLots, setParkingLots] = useState<ParkingLot[]>([]);
  const [start_time, setStartTime] = useState<Date | undefined>();
  const [end_time, setEndTime] = useState<Date | undefined>();
  const [editedAddress, setEditedAddress] = useState(user.Address);
  const [editedVehicle, setEditedVehicle] = useState<Vehicle | null>(null);
  const [editedLicense, setEditedLicense] = useState(user.DriverLicense);

  useEffect(() => {
    const fetchReservations = async () => {
      try {
        const res = await reservationApi.getUserReservations(user.ID);
        setReservations(res);
        console.log(res[0]);
      } catch (err) {
        console.error("Failed to fetch reservations:", err);
      }
    };

    fetchReservations();
  }, [user.ID]);

  useEffect(() => {
    const fetchLots = async () => {
      try {
        const data = await getParkingLots();
        setParkingLots(data);
      } catch (err) {
        console.error("Failed to fetch parking lots:", err);
      }
    };
    fetchLots();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("user");
    navigate("/login");
  };
  const handleEditClick= async (section: string) =>  {
    if (editingSection === section) {
      setEditingSection(null);
    } else {
      setEditingSection(section);
      if (section === "address") setEditedAddress(user.Address);
      if (section === "vehicle") {
        if (user.Vehicle === '000000000000000000000000') {
          setEditedVehicle({
            ID: '',
            Name: '',
            Model: '',
            Year: 0,
            PlateNumber: ''
          });
        } else {
          try {
            const vehicle = await vehicleApi.getVehicle(user.Vehicle);  
            setEditedVehicle(vehicle);
          } catch (err) {
            console.error("Failed to fetch vehicle:", err);
          }
        }
      }
        if (section === "license") setEditedLicense(user.DriverLicense);
    }
  };
  const handleSaveAddress = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const currentUser = await userApi.getUserProfile(user.ID);
      const updatedUser = {
        ...currentUser,  
        Address: editedAddress,  
      };
      await userApi.updateUserProfile(user.ID, updatedUser);
      set_user(updatedUser);
      setEditingSection(null);  
    } catch (err) {
      console.error("Failed to update address:", err);
    }
  };
    
  const handleSaveVehicle = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editedVehicle?.ID === '000000000000000000000000') {
        const { id } = await vehicleApi.createVehicle({
          Name: editedVehicle.Name,
          Model: editedVehicle.Model,
          Year: editedVehicle.Year,
          PlateNumber: editedVehicle.PlateNumber,
        });
  
        const updatedUser = { 
          ...user,
          Vehicle: id
        };
  
        await userApi.updateUserProfile(user.ID, updatedUser);
        set_user(updatedUser);
      } else if (editedVehicle) {
        await vehicleApi.updateVehicle(editedVehicle.ID, {
          Name: editedVehicle.Name,
          Model: editedVehicle.Model,
          Year: editedVehicle.Year,
          PlateNumber: editedVehicle.PlateNumber,
        });
      }
      
      setEditingSection(null);
    } catch (err) {
      console.error("Failed to update vehicle:", err);
    }
  };
      const handleSaveLicense = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editedLicense?.ExpirationDate) {
      return;
    }
    const expirationDate = new Date(editedLicense?.ExpirationDate);
    try {
      const expirationDateISO = expirationDate.toISOString();
      const currentUser = await userApi.getUserProfile(user.ID);
      const updatedUser = {
        ...currentUser,
        DriverLicense: {
          ...currentUser.DriverLicense,  
          Number: editedLicense?.Number,
          State: editedLicense?.State,
          ExpirationDate: expirationDateISO,
        },
      };
      await userApi.updateUserProfile(user.ID, updatedUser);
      set_user(updatedUser);
      setEditingSection(null);  
    } catch (err) {
      console.error("Failed to update license:", err);
    }
  };
  
  const validDate = (current: any) => {
    const date = new Date(current);
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const day = date.getDay(); // Sunday = 0, Saturday = 6

    return day !== 0 && day !== 6 && date >= today;
  };

  const calculateCost = (
    start: Date | null,
    end: Date | null,
    fee: number
  ): string => {
    if (!start || !end) return "0.00";
    const durationInHours =
      (end.getTime() - start.getTime()) / (1000 * 60 * 60);
    return `$${(durationInHours * fee).toFixed(2)}`;
  };
  const changeAddress = (addr?: {
    Street: string;
    City: string;
    State: string;
    ZipCode: string;
  }) => ({
    Street: addr?.Street ?? "",
    City: addr?.City ?? "",
    State: addr?.State ?? "",
    ZipCode: addr?.ZipCode ?? "",
  });
  const change_liscense = (license?: {
    Number: string;
    State: string;
    ExpirationDate: string;
  }) => ({
    Number: license?.Number ?? "",
    State: license?.State ?? "",
    ExpirationDate: license?.ExpirationDate ?? "",
  });

  return (
    <div className="page-container">
      <NavBar user={user} />
      <div className="profile-container">
        <h1>Profile Information</h1>

        <div className="user-info">
          <h2>{`${user.Name.First} ${user.Name.Last}`}</h2>
          <p>Username: {user.Username}</p>
          <p>Sbu ID: {user.SbuID}</p>
          <p>Vehicle Plate: {user.Vehicle}</p>
          <h3>Address:</h3>
          {user.Address && user.Address?.City !== "" ? (
            <div>
              <p>{user.Address.Street}</p>
              <p>
                {user.Address.City}, {user.Address.State} {user.Address.ZipCode}
              </p>
            </div>
          ) : (
            <p>No address provided</p>
          )}
          <h3>Driver's License:</h3>
          {user.DriverLicense && user.DriverLicense.Number !== "" ? (
            <div>
              <p>{user.DriverLicense.Number}</p>
              <p>{user.DriverLicense.State}</p>
              <p>
                Expires:{" "}
                {new Date(user.DriverLicense.ExpirationDate).toLocaleDateString('en-US', {timeZone: 'UTC'})
                }
              </p>
            </div>
          ) : (
            <p>No driver's license information provided</p>
          )}
        </div>

        <div className="edit-options">
          <button onClick={() => handleEditClick("address")}>
            {editingSection === "address" ? "Cancel" : "Edit Address"}
          </button>
          <button onClick={() => handleEditClick("vehicle")}>
            {editingSection === "vehicle" ? "Cancel" : "Edit Vehicle"}
          </button>
          <button onClick={() => handleEditClick("license")}>
            {editingSection === "license" ? "Cancel" : "Edit License"}
          </button>
        </div>

        {editingSection === "address" && user.Address && (
          <div className="edit-form">
            <h3>Edit Address</h3>
            <form onSubmit={handleSaveAddress}>
              <label>Street:</label>
              <input
                type="text"
                value={editedAddress?.Street || ""}
                onChange={(e) =>
                  setEditedAddress((prev) => ({
                    ...changeAddress(prev),
                    Street: e.target.value,
                  }))
                }
              />
              <label>City:</label>
              <input
                type="text"
                value={editedAddress?.City || ""}
                onChange={(e) =>
                  setEditedAddress((prev) => ({
                    ...changeAddress(prev),
                    City: e.target.value,
                  }))
                }
              />
              <label>State:</label>
              <input
                type="text"
                value={editedAddress?.State || ""}
                onChange={(e) =>
                  setEditedAddress((prev) => ({
                    ...changeAddress(prev),
                    State: e.target.value,
                  }))
                }
              />
              <label>Zip Code:</label>
              <input
                type="text"
                value={editedAddress?.ZipCode || ""}
                onChange={(e) =>
                  setEditedAddress((prev) => ({
                    ...changeAddress(prev),
                    ZipCode: e.target.value,
                  }))
                }
              />
              <button type="submit">Save Address</button>
            </form>
          </div>
        )}

{editingSection === "vehicle" && editedVehicle && (
  <div className="edit-form">
    <h3>Edit Vehicle</h3>
    <form onSubmit={handleSaveVehicle}>
      <label>Vehicle Name:</label>
      <input
        type="text"
        value={editedVehicle.Name}
        onChange={(e) => setEditedVehicle((prev) => prev === null ? {
          ID: '',
          Name: '',
          Model: '',
          Year: 0,
          PlateNumber: ''
        } : {
          ...prev,
          Name: e.target.value,
        })
      }
    
      />
      <label>Model:</label>
      <input
        type="text"
        value={editedVehicle.Model}
        onChange={(e) => setEditedVehicle((prev) => prev === null ? {
          ID: '',
          Name: '',
          Model: '',
          Year: 0,
          PlateNumber: ''
        } : {
          ...prev,
          Model: e.target.value,
        })
      }
    
      />
      <label>Year:</label>
      <input
        type="number"
        value={editedVehicle.Year}
        onChange={(e) => setEditedVehicle((prev) => prev === null ? {
          ID: '',
          Name: '',
          Model: '',
          Year: 0,
          PlateNumber: ''
        } : {
          ...prev,
          Year: Number(e.target.value),
        })
      }
    
      />
      <label>Plate Number:</label>
      <input
        type="text"
        value={editedVehicle.PlateNumber}
        onChange={(e) => setEditedVehicle((prev) => prev === null ? {
          ID: '',
          Name: '',
          Model: '',
          Year: 0,
          PlateNumber: ''
        } : {
          ...prev,
          PlateNumber: e.target.value,
        })
      }
    
      />
      <button type="submit">Save Vehicle</button>
    </form>
  </div>
)}

        {editingSection === "license" && user.DriverLicense && (
          <div className="edit-form">
            <h3>Edit Driver's License</h3>
            <form onSubmit={handleSaveLicense}>
              <label>License Number:</label>
              <input
                type="text"
                value={editedLicense?.Number || ""}
                onChange={(e) =>
                  setEditedLicense((prev) => ({
                    ...change_liscense(prev),
                    Number: e.target.value,
                  }))
                }
              />
              <label>State:</label>
              <input
                type="text"
                value={editedLicense?.State || ""}
                onChange={(e) =>
                  setEditedLicense((prev) => ({
                    ...change_liscense(prev),
                    State: e.target.value,
                  }))
                }
              />
              <label>Expiration Date:</label>
              <input
                type="date"
                value={
                  editedLicense?.ExpirationDate
                    ? editedLicense?.ExpirationDate.split("T")[0]
                    : ""
                }
                onChange={(e) =>
                  setEditedLicense((prev) => ({
                    ...change_liscense(prev),
                    ExpirationDate: e.target.value,
                  }))
                }
              />
              <button type="submit">Save License</button>
            </form>
          </div>
        )}
        <div className="reservations-table">
          <h3>My Reservations</h3>
          {reservations.length === 0 ? (
            <p>No reservations found.</p>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Start Time</th>
                  <th>End Time</th>
                  <th>Status</th>
                  <th>Parking Lot</th>
                  <th>Cost</th>
                </tr>
              </thead>
              <tbody>
                {reservations.map((res) => (
                  <tr key={res.ID}>
                    {editingId === res.ID ? (
                      <>
                        <td>
                          <Datetime
                            value={editedStartTime ?? undefined}
                            onChange={(date) => {
                              setEditedStartTime(new Date(date.toString()));
                            }}
                            isValidDate={validDate}
                          />
                        </td>

                        <td>
                          <Datetime
                            value={editedEndTime ?? undefined}
                            onChange={(date) => {
                              setEditedEndTime(new Date(date.toString()));
                            }}
                            isValidDate={(current) => {
                              const date = new Date(current);
                              return (
                                validDate(current) &&
                                (!editedStartTime || date > editedStartTime)
                              );
                            }}
                          />
                        </td>

                        <td>{res.Status}</td>

                        <td>
                          {parkingLots.find((lot) => lot.ID === res.ParkingLot)
                            ?.Name || res.ParkingLot}
                        </td>
                        <td>
                          {(() => {
                            const lot = parkingLots.find(
                              (lot) => lot.ID === res.ParkingLot
                            );
                            if (!lot || !editedStartTime || !editedEndTime)
                              return "0.00";
                            return calculateCost(
                              editedStartTime,
                              editedEndTime,
                              lot.Fee
                            );
                          })()}
                        </td>

                        <td>
                          <button
                            className="confirm-btn"
                            onClick={async () => {
                              if (!editedStartTime || !editedEndTime) return;

                              try {
                                await reservationApi.updateReservation(res.ID, {
                                  StartTime: editedStartTime.toISOString(),
                                  EndTime: editedEndTime.toISOString(),
                                });
                                setReservations((prev) =>
                                  prev.map((r) =>
                                    r.ID === res.ID
                                      ? {
                                          ...r,
                                          StartTime:
                                            editedStartTime.toISOString(),
                                          EndTime: editedEndTime.toISOString(),
                                        }
                                      : r
                                  )
                                );
                                setEditingId(null);
                              } catch (err) {
                                console.error("Update failed:", err);
                              }
                            }}
                          >
                            Save
                          </button>
                        </td>
                        <td>
                          <button
                            className="cancel-confirm-btn"
                            onClick={() => setEditingId(null)}
                          >
                            Cancel
                          </button>
                        </td>
                      </>
                    ) : (
                      <>
                        <td>{new Date(res.StartTime).toLocaleString()}</td>
                        <td>{new Date(res.EndTime).toLocaleString()}</td>
                        <td>{res.Status}</td>
                        <td>
                          {
                            parkingLots.find((lot) => lot.ID === res.ParkingLot)
                              ?.Name || res.ParkingLot // fallback to ID if not found
                          }
                        </td>
                        <td>
                          {(() => {
                            const lot = parkingLots.find(
                              (lot) => lot.ID === res.ParkingLot
                            );
                            if (!lot) return "0.00";
                            return calculateCost(
                              new Date(res.StartTime),
                              new Date(res.EndTime),
                              lot.Fee
                            );
                          })()}
                        </td>

                        <td>
                          {confirmCancelId === res.ID ? (
                            <div
                              style={{
                                display: "flex",
                                flexDirection: "column",
                                gap: "5px",
                              }}
                            >
                              <span>Cancel this reservation?</span>
                              <div>
                                <button
                                  className="confirm-btn"
                                  onClick={async () => {
                                    try {
                                      await reservationApi.cancelReservation(
                                        res.ID
                                      );
                                      setReservations((prev) =>
                                        prev.filter((r) => r.ID !== res.ID)
                                      );
                                    } catch (err) {
                                      console.error("Cancel failed:", err);
                                    } finally {
                                      setConfirmCancelId(null);
                                    }
                                  }}
                                >
                                  Yes, Cancel
                                </button>
                                <button
                                  className="cancel-confirm-btn"
                                  onClick={() => setConfirmCancelId(null)}
                                >
                                  No, Go Back
                                </button>
                              </div>
                            </div>
                          ) : (
                            <button
                              className="cancel-btn"
                              onClick={() => setConfirmCancelId(res.ID)}
                            >
                              Cancel
                            </button>
                          )}
                        </td>
                        <td>
                          <button
                            className="edit-btn"
                            onClick={() => {
                              setEditingId(res.ID);
                              setEditedStartTime(new Date(res.StartTime));
                              setEditedEndTime(new Date(res.EndTime));
                              setConfirmCancelId(null);
                              console.log(reservations);
                            }}
                          >
                            Edit
                          </button>
                        </td>
                      </>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        <button onClick={handleLogout} className="logout-button">
          Logout
        </button>
      </div>
    </div>
  );
};

export default Profile;
