package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
    Lat float64 `bson:"lat" json:"lat"`
    Lng float64 `bson:"lng" json:"lng"`
}

// ParkingLot represents a parking lot document in MongoDB
type ParkingLot struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Spaces   int                `bson:"spaces"`
	Faculty  int                `bson:"faculty"`
	Premium  int                `bson:"premium"`
	Metered  int                `bson:"metered"`
	Resident int                `bson:"resident"`
	Ada      int                `bson:"ada"`
	Ev       bool               `bson:"ev"`
	Active   bool               `bson:"active"`
	Location Location          `bson:"location"`
	Node     primitive.ObjectID `bson:"node"`
	Fee      float64            `bson:"fee`
	Occupancy int               `bson:"occupancy"`
}

// User represents a user document in MongoDB
type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name struct {
		First string `bson:"first"`
		Last  string `bson:"last"`
	} `bson:"name"`
	Role    string             `bson:"role"`
	SbuID   string             `bson:"sbuId"`
	Vehicle primitive.ObjectID `bson:"vehicle"`
	Address struct {
		Street  string `bson:"street"`
		City    string `bson:"city"`
		State   string `bson:"state"`
		ZipCode string `bson:"zipCode"`
	} `bson:"address"`
	DriverLicense struct {
		Number         string    `bson:"number"`
		State          string    `bson:"state"`
		ExpirationDate time.Time `bson:"expirationDate"`
	} `bson:"driverLicense"`
	Username     string `bson:"username,omitempty"`
	PasswordHash string `bson:"passwordHash"`
}

// Vehicle represents a vehicle document in MongoDB
type Vehicle struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Model       string             `bson:"model"`
	Year        int                `bson:"year"`
	PlateNumber string             `bson:"plateNumber"`
}

// Building represents a building document in MongoDB


type Building struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Location Location          `bson:"location"`
	Node     primitive.ObjectID `bson:"node"`
}

// Reservation represents a reservation document in MongoDB
type Reservation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ReservedBy primitive.ObjectID `bson:"reserved_by"`
	ParkingLot primitive.ObjectID `bson:"parking_lot"`
	StartTime  time.Time          `bson:"start_time"`
	EndTime    time.Time          `bson:"end_time"`
	Status     string             `bson:"status"`
}

// Violation represents a violation document in MongoDB
type Violation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	User       primitive.ObjectID `bson:"user"`
	ParkingLot primitive.ObjectID `bson:"parking_lot"`
	Reason     string             `bson:"reason"`
	Fine       float64            `bson:"fine"`
	PayBy      time.Time          `bson:"pay_by"`
}

// ViolationRebute represents a violation rebute document in MongoDB
type ViolationRebute struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	User       primitive.ObjectID `bson:"user"`
	ParkingLot primitive.ObjectID `bson:"parking_lot"`
	Violation  primitive.ObjectID `bson:"violation"`
	Reason     string             `bson:"reason"`
}

// Node represents a node document in MongoDB
type Node struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Type       string             `bson:"type"`
	ParkingLot primitive.ObjectID `bson:"parking_lot,omitempty"`
	Building   primitive.ObjectID `bson:"building,omitempty"`
	Neighbors  struct {
		Members []primitive.ObjectID `bson:"members"`
	} `bson:"neighbors"`
	Location   Location `bson:"location"`
}

// Feedback from users
type Feedbacks struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"` 
	User       primitive.ObjectID `bson:"user"`
	Feedback   string 			  `bson:"feedback"`
}

// Message for users and admin
type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Type       string             `bson:"type"`
	User       primitive.ObjectID `bson:"user"`
	Reserve    primitive.ObjectID `bson:"reserved_by"`
	Violation  primitive.ObjectID `bson:"violation"`
	VioRebut   primitive.ObjectID `bson:"violation_rebute"`
	FeedBack   primitive.ObjectID `bson:"_id,omitempty"` 
}

