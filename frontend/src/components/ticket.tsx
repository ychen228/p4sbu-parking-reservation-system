import React, { useState, useEffect } from "react";
import { User, isAuthenticated } from "../api/userApi";
import violationApi from "../api/violationApi";
import { Violation } from "../types/models";
import NavBar from "./top_nav";
import "./css_files/profile.css";
import adminApi from "../api/adminApi";
import "./css_files/tickets.css";

interface ProfileProps {
  user: User;
}

const Tickets: React.FC<ProfileProps> = ({ user }) => {
  const [violations, setViolations] = useState<Violation[]>([]);
  const [all_violations, set_all_Violations] = useState<Violation[]>([]);
  const [users, set_users] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    // Check authentication on component mount
    if (!isAuthenticated()) {
      const temp = [
        {
          ID: "u001",
          Name: { First: "Alice", Last: "Johnson" },
          Role: "admin",
          SbuID: "SBU123456",
          Vehicle: "Toyota Camry",
          Address: {
            Street: "123 Main St",
            City: "Stony Brook",
            State: "NY",
            ZipCode: "11790",
          },
          DriverLicense: {
            Number: "D1234567",
            State: "NY",
            ExpirationDate: "2027-05-01",
          },
          Username: "alicej",
          PasswordHash: "hashedpassword1",
          name: { first: "alice", last: "johnson" },
        },
        {
          ID: "u002",
          Name: { First: "Bob", Last: "Smith" },
          Role: "user",
          SbuID: "SBU654321",
          Vehicle: "Honda Accord",
          Username: "bobsmith",
          PasswordHash: "hashedpassword2",
          name: { first: "bob", last: "smith" },
        },
        {
          ID: "u003",
          Name: { First: "Carol", Last: "Lee" },
          Role: "user",
          SbuID: "SBU112233",
          Vehicle: "Tesla Model 3",
          Address: {
            Street: "456 Elm St",
            City: "Setauket",
            State: "NY",
            ZipCode: "11733",
          },
          DriverLicense: {
            Number: "D7654321",
            State: "NY",
            ExpirationDate: "2026-11-15",
          },
          Username: "carollee",
          PasswordHash: "hashedpassword3",
          name: { first: "carol", last: "lee" },
        },
      ];

      set_users(temp);
      console.log(users);
      return;
    }
    // fetch_info();
    const temp = [
      {
        ID: "u001",
        Name: { First: "Alice", Last: "Johnson" },
        Role: "admin",
        SbuID: "SBU123456",
        Vehicle: "Toyota Camry",
        Address: {
          Street: "123 Main St",
          City: "Stony Brook",
          State: "NY",
          ZipCode: "11790",
        },
        DriverLicense: {
          Number: "D1234567",
          State: "NY",
          ExpirationDate: "2027-05-01",
        },
        Username: "alicej",
        PasswordHash: "hashedpassword1",
        name: { first: "alice", last: "johnson" },
      },
      {
        ID: "u002",
        Name: { First: "Bob", Last: "Smith" },
        Role: "user",
        SbuID: "SBU654321",
        Vehicle: "Honda Accord",
        Username: "bobsmith",
        PasswordHash: "hashedpassword2",
        name: { first: "bob", last: "smith" },
      },
      {
        ID: "u003",
        Name: { First: "Carol", Last: "Lee" },
        Role: "user",
        SbuID: "SBU112233",
        Vehicle: "Tesla Model 3",
        Address: {
          Street: "456 Elm St",
          City: "Setauket",
          State: "NY",
          ZipCode: "11733",
        },
        DriverLicense: {
          Number: "D7654321",
          State: "NY",
          ExpirationDate: "2026-11-15",
        },
        Username: "carollee",
        PasswordHash: "hashedpassword3",
        name: { first: "carol", last: "lee" },
      },
    ];

    set_users(temp);
  }, []);

  useEffect(() => {
    const fetchViolations = async () => {
      try {
        const fetchedViolations = await violationApi.getUserViolations(user.ID);
        setViolations(fetchedViolations);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch violations.");
        setLoading(false);
      }
    };

    fetchViolations();
  }, [user.ID]);
  const fetch_info = async () => {
    try {
      const every_user = await adminApi.getAllUsers();
      set_users(every_user);
      const every_violation = await violationApi.getAllViolations();
      set_all_Violations(every_violation);
      setLoading(false);
    } catch (err) {
      setError("Failed to get information needed.");
      setLoading(false);
      const temp = [
        {
          ID: "u001",
          Name: { First: "Alice", Last: "Johnson" },
          Role: "admin",
          SbuID: "SBU123456",
          Vehicle: "Toyota Camry",
          Address: {
            Street: "123 Main St",
            City: "Stony Brook",
            State: "NY",
            ZipCode: "11790",
          },
          DriverLicense: {
            Number: "D1234567",
            State: "NY",
            ExpirationDate: "2027-05-01",
          },
          Username: "alicej",
          PasswordHash: "hashedpassword1",
          name: { first: "alice", last: "johnson" },
        },
        {
          ID: "u002",
          Name: { First: "Bob", Last: "Smith" },
          Role: "user",
          SbuID: "SBU654321",
          Vehicle: "Honda Accord",
          Username: "bobsmith",
          PasswordHash: "hashedpassword2",
          name: { first: "bob", last: "smith" },
        },
        {
          ID: "u003",
          Name: { First: "Carol", Last: "Lee" },
          Role: "user",
          SbuID: "SBU112233",
          Vehicle: "Tesla Model 3",
          Address: {
            Street: "456 Elm St",
            City: "Setauket",
            State: "NY",
            ZipCode: "11733",
          },
          DriverLicense: {
            Number: "D7654321",
            State: "NY",
            ExpirationDate: "2026-11-15",
          },
          Username: "carollee",
          PasswordHash: "hashedpassword3",
          name: { first: "carol", last: "lee" },
        },
      ];

      set_users(temp);
    }
  };
  return (
    <div className="page-container">
      <NavBar user={user} />
      <div className="profile-container">
        <h2>Your Tickets</h2>
        {loading ? (
          <p>Loading tickets...</p>
        ) : error ? (
          <p>{error}</p>
        ) : (
          <table className="ticket-table">
            <thead>
              <tr>
                <th>Parking Lot</th>
                <th>Fine</th>
                <th>Pay By</th>
              </tr>
            </thead>
            <tbody>
              {violations.length > 0 ? (
                violations.map((violation) => (
                  <tr key={violation.ID}>
                    <td>{violation.ParkingLot}</td>
                    <td>{violation.Fine}</td>
                    <td>{violation.PayBy}</td>
                    <td>
                      <button
                        className="btn btn-primary"
                        // onClick={() => view_violation(pl)}
                      >
                        View
                      </button>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={4}>No violations found.</td>
                </tr>
              )}
            </tbody>
          </table>
        )}
        {user.Role === "admin" && (
          <div className="admin_tickets">
            <div className="section-body">
              <h2>Users</h2>
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
                          {user.Name && typeof user.Name === "object"
                            ? `${user.Name.First || ""} ${user.Name.Last || ""}`
                            : user.name && typeof user.name === "object"
                            ? `${user.name.first || ""} ${user.name.last || ""}`
                            : user.Name || "Unknown"}
                        </td>
                        <td>{user.Username || "Unknown"}</td>
                        <td>{user.Role || "Unknown"}</td>
                        <td>{user.SbuID || "Unknown"}</td>
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
            <h2>All Tickets</h2>
            <table className="ticket-table">
              <thead>
                <tr>
                  <th>Parking Lot</th>
                  <th>Fine</th>
                  <th>Pay By</th>
                </tr>
              </thead>
              <tbody>
                {all_violations.length > 0 ? (
                  all_violations.map((violation) => (
                    <tr key={violation.ID}>
                      <td>{violation.ParkingLot}</td>
                      <td>{violation.Fine}</td>
                      <td>{violation.PayBy}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={4}>No violations found.</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default Tickets;
