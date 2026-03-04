import React, { useState, Dispatch, SetStateAction } from "react";
import Datetime from "react-datetime";
import { ParkingLot } from "../types/models";
import { User } from "../api/userApi";
import reservationApi from "../api/reservationApi";
import "./css_files/parkingLotInfo.css";
import "react-datetime/css/react-datetime.css";

interface ParkingInfoProps {
  parkingLot: ParkingLot | null;
  set_selected_parking: Dispatch<SetStateAction<ParkingLot | null>>;
  user: User;
}

const ParkingLotInfo: React.FC<ParkingInfoProps> = ({
  parkingLot,
  set_selected_parking,
  user,
}) => {
  const [start_time, setStartTime] = useState<Date | undefined>();
  const [end_time, setEndTime] = useState<Date | undefined>();
  const [reservation_error, setReservationError] = useState<string | null>(
    null,
  );

  // Calculate parking lot availability status
  const getAvailabilityStatus = () => {
    if (!parkingLot) return null;

    const totalSpaces = parkingLot.Spaces;
    const occupancy = parkingLot.Occupancy || 0;
    const availableSpaces = totalSpaces - occupancy;

    const availabilityPercent = (availableSpaces / totalSpaces) * 100;

    if (availabilityPercent > 30) return "available";
    if (availabilityPercent > 10) return "limited";
    return "full";
  };

  const availabilityStatus = getAvailabilityStatus();

  const validDate = (current: Date) => {
    const date = new Date(current);
    const today = new Date();
    today.setHours(0, 0, 0, 0); // Normalize today's time to midnight
    const day = date.getDay(); // 0 = Sunday, 6 = Saturday

    return day !== 0 && day !== 6 && date >= today;
  };

  const make_reservation = async (e: React.FormEvent) => {
    e.preventDefault();
    setReservationError(null);

    if (!start_time || !end_time) {
      setReservationError("Please select both start and end times.");
      return;
    }

    if (start_time >= end_time) {
      setReservationError("End time must be after start time.");
      return;
    }

    const currentTime = new Date();
    if (start_time < currentTime) {
      setReservationError("Start time cannot be in the past.");
      return;
    }

    // Here you would make an API call to create the reservation
    // For now, just close the info panel
    try {
      await reservationApi.createReservation({
        ReservedBy: user.ID,
        ParkingLot: parkingLot!.ID,
        StartTime: start_time.toISOString(),
        EndTime: end_time.toISOString(),
      });

      alert("Reservation successfully created!");
      set_selected_parking(null);
      try {
        const data = await reservationApi.getUserReservations(user.ID);
        console.log("My reservations:", data);
      } catch (err) {
        console.error("Failed to fetch reservations:", err);
      }
    } catch (error) {
      console.error("Reservation failed:", error);
      setReservationError("Failed to make reservation. Please try again.");
    }
  };

  const find_cost = () => {
    if (!start_time || !end_time || !parkingLot) return 0;

    // Calculate time difference in hours
    const reserved_time =
      (end_time.getTime() - start_time.getTime()) / (1000 * 60 * 60);

    // Simple calculation for now - will be updated when schema changes
    const cost = reserved_time * parkingLot.Fee;

    return cost.toFixed(2);
  };

  const formatDateTime = (date?: Date) => {
    if (!date) return "Not selected";
    return date.toLocaleString("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatHours = (startDate?: Date, endDate?: Date) => {
    if (!startDate || !endDate) return "0";
    const hours = (endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60);
    return hours.toFixed(1);
  };

  if (!parkingLot) return null;

  return (
    <div className="parking-reserve">
      <div className="parking-header">
        <h2 className="parking-title">{parkingLot.Name}</h2>
        <button
          className="close-button"
          onClick={() => set_selected_parking(null)}
          aria-label="Close"
        >
          ×
        </button>
      </div>

      <div className="parking-details">
        {/* Parking lot basic info */}
        <div className="detail-item">
          <span className="detail-label">Status</span>
          <span className="detail-value">
            <span
              className={`availability-indicator ${availabilityStatus}`}
            ></span>
            {availabilityStatus === "available"
              ? "Available"
              : availabilityStatus === "limited"
                ? "Limited Spaces"
                : "Full"}
          </span>
        </div>

        <div className="detail-item">
          <span className="detail-label">Total Spaces</span>
          <span className="detail-value">{parkingLot.Spaces}</span>
        </div>

        <div className="detail-item">
          <span className="detail-label">Hourly Rate</span>
          <span className="detail-value">${parkingLot.Fee.toFixed(2)}</span>
        </div>

        {/* Special features section */}
        <div className="amenities">
          {parkingLot.Faculty > 0 && (
            <span className="amenity-tag">Faculty: {parkingLot.Faculty}</span>
          )}
          {parkingLot.Premium > 0 && (
            <span className="amenity-tag special">
              Premium: {parkingLot.Premium}
              <span className="badge premium">P</span>
            </span>
          )}
          {parkingLot.Metered > 0 && (
            <span className="amenity-tag">Metered: {parkingLot.Metered}</span>
          )}
          {parkingLot.Resident > 0 && (
            <span className="amenity-tag">Resident: {parkingLot.Resident}</span>
          )}
          {parkingLot.Ada > 0 && (
            <span className="amenity-tag special">
              ADA: {parkingLot.Ada}
              <span className="badge ada">♿</span>
            </span>
          )}
          {parkingLot.Ev && (
            <span className="amenity-tag special">
              EV Charging
              <span className="badge ev">⚡</span>
            </span>
          )}
        </div>
      </div>

      <h3 className="section-title">Make a Reservation</h3>

      <form className="reservation-form" onSubmit={make_reservation}>
        <div className="time-range">
          <div className="form-group">
            <label htmlFor="start-time">Start Time</label>
            {/*
            <input
              id="start-time"
              type="datetime-local"
              onChange={(e) => setStartTime(new Date(e.target.value))}
              required
            />
            */}
            <Datetime
              inputProps={{
                id: "start-time",
                className: "custom-datetime-input",
                required: true,
              }}
              onChange={(date) => setStartTime(new Date(date.toString()))}
              isValidDate={validDate}
            />
          </div>

          <span className="time-separator">→</span>

          <div className="form-group">
            <label htmlFor="end-time">End Time</label>
            {/*
            <input
              id="end-time"
              type="datetime-local"
              onChange={(e) => setEndTime(new Date(e.target.value))}
              required
            />
            */}
            <Datetime
              inputProps={{
                id: "start-time",
                className: "custom-datetime-input",
                required: true,
              }}
              onChange={(date) => setEndTime(new Date(date.toString()))}
              isValidDate={validDate}
            />
          </div>
        </div>

        {reservation_error && (
          <div className="error-message">{reservation_error}</div>
        )}

        {start_time && end_time && (
          <div className="cost-summary">
            <div className="cost-item">
              <span>Duration</span>
              <span>{formatHours(start_time, end_time)} hours</span>
            </div>
            <div className="cost-item">
              <span>From</span>
              <span>{formatDateTime(start_time)}</span>
            </div>
            <div className="cost-item">
              <span>To</span>
              <span>{formatDateTime(end_time)}</span>
            </div>
            <div className="cost-item">
              <span>Rate</span>
              <span>${parkingLot.Fee.toFixed(2)}/hr</span>
            </div>
            <div className="cost-item total-cost">
              <span>Total Cost</span>
              <span>${find_cost()}</span>
            </div>
          </div>
        )}

        <div className="action-buttons">
          <button type="submit" className="reserve-button">
            Reserve Parking
          </button>
          <button
            type="button"
            className="cancel-button"
            onClick={() => set_selected_parking(null)}
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
};

export default ParkingLotInfo;
