package mongodb

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/YimingChen/P4SBU/Backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository provides data access methods for MongoDB
type Repository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

// NewRepository creates a new MongoDB repository
func NewRepository(client *mongo.Client, database string) *Repository {
	return &Repository{
		client:   client,
		database: database,
		timeout:  30 * time.Second,
	}
}

// Generic functions for CRUD operations
func (r *Repository) getCollection(collection string) *mongo.Collection {
	return r.client.Database(r.database).Collection(collection)
}

func (r *Repository) findOne(collection string, filter interface{}) (bson.Raw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result := r.getCollection(collection).FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var doc bson.Raw
	err := result.Decode(&doc)
	return doc, err
}

func (r *Repository) find(collection string, filter interface{}, opts ...*options.FindOptions) ([]bson.Raw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cursor, err := r.getCollection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.Raw
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) insertOne(collection string, document interface{}) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.getCollection(collection).InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("failed to get inserted ID")
	}

	return id, nil
}

func (r *Repository) updateOne(collection string, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.getCollection(collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *Repository) deleteOne(collection string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.getCollection(collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// ParkingLot operations
func (r *Repository) GetParkingLot(id primitive.ObjectID) (*models.ParkingLot, error) {
	doc, err := r.findOne("parking_lots", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var parkingLot models.ParkingLot
	if err := bson.Unmarshal(doc, &parkingLot); err != nil {
		return nil, err
	}

	return &parkingLot, nil
}

func (r *Repository) GetAllParkingLots() ([]models.ParkingLot, error) {

	log.Printf("GetAllParkingLots: Querying collection 'parking_lots' in database '%s'", r.database)

	docs, err := r.find("parking_lots", bson.M{})
	if err != nil {
		return nil, err
	}

	log.Printf("GetAllParkingLots: Query returned %d documents, error: %v", len(docs), err)

	parkingLots := make([]models.ParkingLot, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &parkingLots[i]); err != nil {
			return nil, err
		}
	}

	log.Printf("GetAllParkingLots: Successfully unmarshaled %d parking lots", len(parkingLots))

	return parkingLots, nil
}

func (r *Repository) CreateParkingLot(parkingLot models.ParkingLot) (primitive.ObjectID, error) {
	return r.insertOne("parking_lots", parkingLot)
}

func (r *Repository) UpdateParkingLot(parkingLot models.ParkingLot) error {
	return r.updateOne("parking_lots", bson.M{"_id": parkingLot.ID}, bson.M{"$set": parkingLot})
}

func (r *Repository) DeleteParkingLot(id primitive.ObjectID) error {
	return r.deleteOne("parking_lots", bson.M{"_id": id})
}

// User operations
func (r *Repository) GetUser(id primitive.ObjectID) (*models.User, error) {
	doc, err := r.findOne("users", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := bson.Unmarshal(doc, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	doc, err := r.findOne("users", bson.M{"username": username})
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := bson.Unmarshal(doc, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *Repository) GetAllUsers() ([]models.User, error) {
	docs, err := r.find("users", bson.M{})
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &users[i]); err != nil {
			return nil, err
		}
	}

	return users, nil
}
func (r *Repository) CreateUser(user models.User) (primitive.ObjectID, error) {
	return r.insertOne("users", user)
}

func (r *Repository) UpdateUser(user models.User) error {
	return r.updateOne("users", bson.M{"_id": user.ID}, bson.M{"$set": user})
}

func (r *Repository) DeleteUser(id primitive.ObjectID) error {
	return r.deleteOne("users", bson.M{"_id": id})
}

// Vehicle operations
func (r *Repository) GetVehicle(id primitive.ObjectID) (*models.Vehicle, error) {
	doc, err := r.findOne("vehicles", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var vehicle models.Vehicle
	if err := bson.Unmarshal(doc, &vehicle); err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (r *Repository) CreateVehicle(vehicle models.Vehicle) (primitive.ObjectID, error) {
	return r.insertOne("vehicles", vehicle)
}

func (r *Repository) UpdateVehicle(vehicle models.Vehicle) error {
	return r.updateOne("vehicles", bson.M{"_id": vehicle.ID}, bson.M{"$set": vehicle})
}

func (r *Repository) DeleteVehicle(id primitive.ObjectID) error {
	return r.deleteOne("vehicles", bson.M{"_id": id})
}

// Building operations
func (r *Repository) GetBuilding(id primitive.ObjectID) (*models.Building, error) {
	doc, err := r.findOne("buildings", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var building models.Building
	if err := bson.Unmarshal(doc, &building); err != nil {
		return nil, err
	}

	return &building, nil
}

func (r *Repository) GetAllBuildings() ([]models.Building, error) {
	docs, err := r.find("buildings", bson.M{})
	if err != nil {
		return nil, err
	}

	buildings := make([]models.Building, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &buildings[i]); err != nil {
			return nil, err
		}
	}

	return buildings, nil
}

func (r *Repository) CreateBuilding(building models.Building) (primitive.ObjectID, error) {
	return r.insertOne("buildings", building)
}

func (r *Repository) UpdateBuilding(building models.Building) error {
	return r.updateOne("buildings", bson.M{"_id": building.ID}, bson.M{"$set": building})
}

func (r *Repository) DeleteBuilding(id primitive.ObjectID) error {
	return r.deleteOne("buildings", bson.M{"_id": id})
}

// Reservation operations
func (r *Repository) GetReservation(id primitive.ObjectID) (*models.Reservation, error) {
	doc, err := r.findOne("reservations", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var reservation models.Reservation
	if err := bson.Unmarshal(doc, &reservation); err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *Repository) GetReservationsByUser(userID primitive.ObjectID) ([]models.Reservation, error) {
	docs, err := r.find("reservations", bson.M{"reserved_by": userID})
	if err != nil {
		return nil, err
	}

	reservations := make([]models.Reservation, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &reservations[i]); err != nil {
			return nil, err
		}
	}

	return reservations, nil
}

func (r *Repository) GetReservationsByLot(lotID primitive.ObjectID) ([]models.Reservation, error) {
	docs, err := r.find("reservations", bson.M{"parking_lot": lotID})
	if err != nil {
		return nil, err
	}

	reservations := make([]models.Reservation, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &reservations[i]); err != nil {
			return nil, err
		}
	}

	return reservations, nil
}
func (r *Repository) GetAllReservations() ([]models.Reservation, error) {
	docs, err := r.find("reservations", bson.M{})
	if err != nil {
		return nil, err
	}

	reservations := make([]models.Reservation, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &reservations[i]); err != nil {
			return nil, err
		}
	}

	return reservations, nil
}

func (r *Repository) CreateReservation(reservation models.Reservation) (primitive.ObjectID, error) {
	return r.insertOne("reservations", reservation)
}

func (r *Repository) UpdateReservation(reservation models.Reservation) error {
	return r.updateOne("reservations", bson.M{"_id": reservation.ID}, bson.M{"$set": reservation})
}

func (r *Repository) DeleteReservation(id primitive.ObjectID) error {
	return r.deleteOne("reservations", bson.M{"_id": id})
}

// Violation operations
func (r *Repository) GetViolation(id primitive.ObjectID) (*models.Violation, error) {
	doc, err := r.findOne("violations", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var violation models.Violation
	if err := bson.Unmarshal(doc, &violation); err != nil {
		return nil, err
	}

	return &violation, nil
}

func (r *Repository) GetViolationsByUser(userID primitive.ObjectID) ([]models.Violation, error) {
	docs, err := r.find("violations", bson.M{"user": userID})
	if err != nil {
		return nil, err
	}

	violations := make([]models.Violation, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &violations[i]); err != nil {
			return nil, err
		}
	}

	return violations, nil
}

func (r *Repository) CreateViolation(violation models.Violation) (primitive.ObjectID, error) {
	return r.insertOne("violations", violation)
}

func (r *Repository) UpdateViolation(violation models.Violation) error {
	return r.updateOne("violations", bson.M{"_id": violation.ID}, bson.M{"$set": violation})
}

func (r *Repository) DeleteViolation(id primitive.ObjectID) error {
	return r.deleteOne("violations", bson.M{"_id": id})
}

// ViolationRebute operations
func (r *Repository) GetViolationRebute(id primitive.ObjectID) (*models.ViolationRebute, error) {
	doc, err := r.findOne("violation_rebutes", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var rebute models.ViolationRebute
	if err := bson.Unmarshal(doc, &rebute); err != nil {
		return nil, err
	}

	return &rebute, nil
}

func (r *Repository) GetViolationRebutesByUser(userID primitive.ObjectID) ([]models.ViolationRebute, error) {
	docs, err := r.find("violation_rebutes", bson.M{"user": userID})
	if err != nil {
		return nil, err
	}

	rebutes := make([]models.ViolationRebute, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &rebutes[i]); err != nil {
			return nil, err
		}
	}

	return rebutes, nil
}

func (r *Repository) CreateViolationRebute(rebute models.ViolationRebute) (primitive.ObjectID, error) {
	return r.insertOne("violation_rebutes", rebute)
}

func (r *Repository) UpdateViolationRebute(rebute models.ViolationRebute) error {
	return r.updateOne("violation_rebutes", bson.M{"_id": rebute.ID}, bson.M{"$set": rebute})
}

func (r *Repository) DeleteViolationRebute(id primitive.ObjectID) error {
	return r.deleteOne("violation_rebutes", bson.M{"_id": id})
}

// Node operations
func (r *Repository) GetNode(id primitive.ObjectID) (*models.Node, error) {
	doc, err := r.findOne("nodes", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var node models.Node
	if err := bson.Unmarshal(doc, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (r *Repository) GetAllNodes() ([]models.Node, error) {
	docs, err := r.find("nodes", bson.M{})
	if err != nil {
		return nil, err
	}

	nodes := make([]models.Node, len(docs))
	for i, doc := range docs {
		if err := bson.Unmarshal(doc, &nodes[i]); err != nil {
			return nil, err
		}
	}

	return nodes, nil
}

func (r *Repository) CreateNode(node models.Node) (primitive.ObjectID, error) {
	return r.insertOne("nodes", node)
}

func (r *Repository) UpdateNode(node models.Node) error {
	return r.updateOne("nodes", bson.M{"_id": node.ID}, bson.M{"$set": node})
}

func (r *Repository) DeleteNode(id primitive.ObjectID) error {
	return r.deleteOne("nodes", bson.M{"_id": id})
}

func (r *Repository) GetMessage(id primitive.ObjectID) (*models.Message, error) {
	doc, err := r.findOne("messages", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var message models.Message
	if err := bson.Unmarshal(doc, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *Repository) GetMessageByUser(userID primitive.ObjectID) ([]models.Message, error) {
	doc, err := r.find("messages", bson.M{"user": userID})
	if err != nil {
		return nil, err
	}

	messages := make([]models.Message, len(doc))
	for i, doc := range doc {
		if err := bson.Unmarshal(doc, &messages[i]); err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func (r *Repository) CreateMessage(message models.Message) (primitive.ObjectID, error) {
	return r.insertOne("messages", message)
}

func (r *Repository) UpdateMessage(message models.Message) error {
	return r.updateOne("messages", bson.M{"_id": message.ID}, bson.M{"$set": message})
}

func (r *Repository) DeleteMessage(id primitive.ObjectID) error {
	return r.deleteOne("messages", bson.M{"_id": id})
}
