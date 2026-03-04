package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/YimingChen/P4SBU/Backend/config"
	"github.com/YimingChen/P4SBU/Backend/database"
	"github.com/YimingChen/P4SBU/Backend/models"
)

func RunDataAnalysis() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg := config.NewConfig()
	client, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	db := client.Database("campus_parking")
	lotCol := db.Collection("parkinglots")
	violCol := db.Collection("violations")
	resCol := db.Collection("reservations")

	analyzeLotUtilization(ctx, lotCol)
	analyzeViolationsByReason(ctx, violCol)
	analyzeReservationsByUserRole(ctx, resCol, db)
}

func analyzeLotUtilization(ctx context.Context, col *mongo.Collection) {
	cursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	fmt.Println("\n--- Parking Lot Utilization ---")
	for cursor.Next(ctx) {
		var lot models.ParkingLot
		if err := cursor.Decode(&lot); err != nil {
			log.Fatal(err)
		}
		if lot.Spaces == 0 {
			continue
		}
		percent := (float64(lot.Occupancy) / float64(lot.Spaces)) * 100
		fmt.Printf("%-35s | %3d/%3d occupied (%.1f%%)\n", lot.Name, lot.Occupancy, lot.Spaces, percent)
	}
}

func analyzeViolationsByReason(ctx context.Context, col *mongo.Collection) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$reason"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	fmt.Println("\n--- Violations by Reason ---")
	for cursor.Next(ctx) {
		var result struct {
			Reason string `bson:"_id"`
			Count  int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-40s | %d\n", result.Reason, result.Count)
	}
}

func analyzeReservationsByUserRole(ctx context.Context, col *mongo.Collection, db *mongo.Database) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "reservedBy"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: "$user"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user.role"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	fmt.Println("\n--- Reservations by User Role ---")
	for cursor.Next(ctx) {
		var result struct {
			Role  string `bson:"_id"`
			Total int    `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Role: %-10s | Total Reservations: %d\n", result.Role, result.Total)
	}
}
