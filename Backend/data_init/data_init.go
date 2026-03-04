package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/YimingChen/P4SBU/Backend/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using hardcoded connection string")
	}

	// Get MongoDB connection parameters
	var connectionString string

	// If environment variables are available, use them
	if os.Getenv("MONGODB_USERNAME") != "" {
		username := os.Getenv("MONGODB_USERNAME")
		password := os.Getenv("MONGODB_PASSWORD")
		host := os.Getenv("MONGODB_HOST")
		database := os.Getenv("MONGODB_DATABASE")
		if database == "" {
			database = "admin"
		}
		connectionOptions := os.Getenv("MONGODB_OPTIONS")

		// Construct connection string
		connectionString = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?%s",
			username, password, host, database, connectionOptions)
	} else {
		// Provide a fallback for testing
		connectionString = "mongodb://localhost:27017"
		log.Println("Using local MongoDB connection: " + connectionString)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Get database
	database := client.Database("campus_parking")

	// Create node first (since other documents reference it)
	nodeCollection := database.Collection("nodes") //models.Location{Lat: number, Lng: number}

	nodeID := primitive.NewObjectID()
	buildingID := primitive.NewObjectID()
	parkingLotID := primitive.NewObjectID()

	// Create building
	buildingCollection := database.Collection("buildings")
	engineering := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Engineering",
		Location: models.Location{Lat: 40.91291783527936, Lng: -73.12462336456686},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, engineering)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", engineering.ID)

	// Create building
	lightEngineering := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Light Engineering",
		Location: models.Location{Lat: 40.91356941186373, Lng: -73.12544385809436},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, lightEngineering)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", lightEngineering.ID)

	// Create building
	heavyEngineering := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Heavy Engineering",
		Location: models.Location{Lat: 40.91251847871489, Lng: -73.12583788040982},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, heavyEngineering)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", heavyEngineering.ID)

	// Create building
	computingCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Computing Center",
		Location: models.Location{Lat: 40.91309649401406, Lng: -73.12610210713902},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, computingCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", computingCenter.ID)

	// Create building
	newComputerScience := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "New Computer Science",
		Location: models.Location{Lat: 40.91274618231478, Lng: -73.12362671988656},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, newComputerScience)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", newComputerScience.ID)

	// Create building
	computerScience := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Computer Science",
		Location: models.Location{Lat: 40.91241688761761, Lng: -73.12208771505827},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, computerScience)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", computerScience.ID)

	// Create building
	javits := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Javits Lecture Center",
		Location: models.Location{Lat: 40.91292133838911, Lng: -73.12205063060078},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, javits)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", javits.ID)

	// Create building
	ECC := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Educational Communications Center",
		Location: models.Location{Lat: 40.913412474851526, Lng: -73.12274843618637},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, ECC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", ECC.ID)

	// Create building
	IACSLauferCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "IACS Laufer Center",
		Location: models.Location{Lat: 40.91201622136837, Lng: -73.12079533081496},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, IACSLauferCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", IACSLauferCenter.ID)

	// Create building
	bioengineering := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Bioengineering",
		Location: models.Location{Lat: 40.912247418930555, Lng: -73.1200781283574},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, bioengineering)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", bioengineering.ID)

	// Create building
	CMM := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Centers for Molecular Medicine",
		Location: models.Location{Lat: 40.91166232171574, Lng: -73.11936345628176},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, CMM)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", CMM.ID)

	// Create building
	lifeSciences := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "LifeSciences",
		Location: models.Location{Lat: 40.911587308877834, Lng: -73.11997886834689},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, lifeSciences)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", lifeSciences.ID)

	// Create building
	SBS := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Social and Behavioral Sciences",
		Location: models.Location{Lat: 40.91299816157528, Lng: -73.12011045108952},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, SBS)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", SBS.ID)

	// Create building
	hummanities := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Hummanities",
		Location: models.Location{Lat: 40.914391077498216, Lng: -73.12113041482887},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, hummanities)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", hummanities.ID)

	// Create building
	administration := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Administration",
		Location: models.Location{Lat: 40.91475994756114, Lng: -73.1203800129342},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, administration)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", administration.ID)

	// Create building
	hiltonGardenInn := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Hilton Garden Inn",
		Location: models.Location{Lat: 40.91380198227215, Lng: -73.11836194182943},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, hiltonGardenInn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", hiltonGardenInn.ID)

	// Create building
	psychologyA := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "PsychologyA",
		Location: models.Location{Lat: 40.91430298883871, Lng: -73.12251465130853},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, psychologyA)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", psychologyA.ID)

	// Create building
	psychologyB := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Psychology B",
		Location: models.Location{Lat: 40.91389557719384, Lng: -73.12210666581586},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, psychologyB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", psychologyB.ID)

	// Create building
	SAC := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Student Activities Center",
		Location: models.Location{Lat: 40.91442411073418, Lng: -73.12430687333803},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, SAC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", SAC.ID)

	// Create building
	ESS := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Earth and Space Sciences",
		Location: models.Location{Lat: 40.91492511259485, Lng: -73.12571296625246},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, ESS)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", ESS.ID)

	// Create building
	simonsCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Simons Center",
		Location: models.Location{Lat: 40.915883061626126, Lng: -73.12704620462928},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, simonsCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", simonsCenter.ID)

	// Create building
	mathTower := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Math Tower",
		Location: models.Location{Lat: 40.91570688814587, Lng: -73.12636865732819},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, mathTower)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", mathTower.ID)

	// Create building
	physics := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Physics",
		Location: models.Location{Lat: 40.91635102018777, Lng: -73.12608452457435},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, physics)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", physics.ID)

	// Create building
	harrimanHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Harriman Hall",
		Location: models.Location{Lat: 40.915591274038086, Lng: -73.1252466972233},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, harrimanHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", harrimanHall.ID)

	// Create building
	accelerator := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Van de Graaff Accelerator",
		Location: models.Location{Lat: 40.91608125625032, Lng: -73.12476585717835},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, accelerator)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", accelerator.ID)

	// Create building
	freyHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Frey Hall",
		Location: models.Location{Lat: 40.91567936099572, Lng: -73.12389160255117},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, freyHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", freyHall.ID)

	// Create building
	chemistry := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Chemistry",
		Location: models.Location{Lat: 40.91648314906162, Lng: -73.12371675162574},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, chemistry)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", chemistry.ID)

	// Create building
	library := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Frank Melville Jr. Memorial Library",
		Location: models.Location{Lat: 40.91555824139875, Lng: -73.12269678789403},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, library)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", library.ID)

	// Create building
	stallerCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Staller Center for the Arts",
		Location: models.Location{Lat: 40.916070245453966, Lng: -73.12146554596075},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, stallerCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", stallerCenter.ID)

	// Create building
	wangCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Charles B.Wang Center",
		Location: models.Location{Lat: 40.91598215901012, Lng: -73.11978260565398},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, wangCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", wangCenter.ID)

	// Create building
	union := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Stony Brook Union",
		Location: models.Location{Lat: 40.91719009360351, Lng: -73.12253286848828},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, union)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", union.ID)

	// Create building
	recreationCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Walter J. Hawrys Campus Recreation Center",
		Location: models.Location{Lat: 40.91687585166851, Lng: -73.12356400995459},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, recreationCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", recreationCenter.ID)

	// Create building
	sportsComplex := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Sports Complex",
		Location: models.Location{Lat: 40.91711273897799, Lng: -73.12460668598156},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, sportsComplex)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", sportsComplex.ID)

	// Create building
	arena := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Island Federal Arena",
		Location: models.Location{Lat: 40.91749896646663, Lng: -73.1259832909322},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, arena)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", arena.ID)

	// Create building
	eastDinning := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "East Side Dining",
		Location: models.Location{Lat: 40.91702637970153, Lng: -73.1209004258759},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, eastDinning)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", eastDinning.ID)

	// Create building
	westDinning := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Side Dining",
		Location: models.Location{Lat: 40.913124765048934, Lng: -73.13051992321891},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westDinning)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westDinning.ID)

	// Create building
	chavezHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Chávez Hall",
		Location: models.Location{Lat: 40.916874349964836, Lng: -73.11959940958609},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, chavezHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", chavezHall.ID)

	// Create building
	tubmanHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Tubman Hall",
		Location: models.Location{Lat: 40.91721895020736, Lng: -73.11893548890274},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, tubmanHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", tubmanHall.ID)

	// Create building
	grayHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Gray Hall",
		Location: models.Location{Lat: 40.91790814531633, Lng: -73.12124915188524},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, grayHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", grayHall.ID)

	// Create building
	irvingHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Irving Hall",
		Location: models.Location{Lat: 40.91783909074487, Lng: -73.11959086229282},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, irvingHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", irvingHall.ID)

	// Create building
	oNeillHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "O'Neill Hall",
		Location: models.Location{Lat: 40.918557519080125, Lng: -73.11950758510791},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, oNeillHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", oNeillHall.ID)

	// Create building
	benedictHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Benedict Hall",
		Location: models.Location{Lat: 40.91958533001034, Lng: -73.11866093369929},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, benedictHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", benedictHall.ID)

	// Create building
	jamesHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "James Hall",
		Location: models.Location{Lat: 40.9193650861622, Lng: -73.12023626044734},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, jamesHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", jamesHall.ID)

	// Create building
	langmuirHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Langmuir Hall",
		Location: models.Location{Lat: 40.92021459410414, Lng: -73.12008358560833},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, langmuirHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", langmuirHall.ID)

	// Create building
	ammannHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Ammann Hall",
		Location: models.Location{Lat: 40.91875789230143, Lng: -73.1210895704151},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, ammannHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", ammannHall.ID)

	// Create building
	mendelsohnComm := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Mendelsohn Community",
		Location: models.Location{Lat: 40.91825973235171, Lng: -73.12018977783566},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, mendelsohnComm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", mendelsohnComm.ID)

	// Create building
	altheticIndoorPractice := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Athletic Indoor Practice Facility",
		Location: models.Location{Lat: 40.92096094549492, Lng: -73.12322616340042},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, altheticIndoorPractice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", altheticIndoorPractice.ID)

	// Create building
	environmentalConservation := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "NYS Department of Environmental Conservation",
		Location: models.Location{Lat: 40.92209526105391, Lng: -73.11984307923392},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, environmentalConservation)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", environmentalConservation.ID)

	// Create building
	healthCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Student Health Center",
		Location: models.Location{Lat: 40.919339012814895, Lng: -73.12177501053392},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, healthCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", healthCenter.ID)

	// Create building
	centralStores := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Central Stores (Warehouse)",
		Location: models.Location{Lat: 40.91712509487659, Lng: -73.128463209992},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, centralStores)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", centralStores.ID)

	// Create building
	centralServices := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Central Services (Receiving)",
		Location: models.Location{Lat: 40.916761390031866, Lng: -73.12754391242264},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, centralServices)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", centralServices.ID)

	// Create building
	coGenPlant := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "CoGen Plant",
		Location: models.Location{Lat: 40.91691414631039, Lng: -73.12872792918736},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, coGenPlant)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", coGenPlant.ID)

	// Create building
	SBVAC := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "SBVAC",
		Location: models.Location{Lat: 40.91596081410035, Lng: -73.12812099616193},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, SBVAC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", SBVAC.ID)

	// Create building
	serviceGroup := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Service Group",
		Location: models.Location{Lat: 40.91615546415202, Lng: -73.12886380249273},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, serviceGroup)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", serviceGroup.ID)

	// Create building
	toscaniniHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Toscanini Hall",
		Location: models.Location{Lat: 40.9101985977337, Lng: -73.12833517688767},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, toscaniniHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", toscaniniHall.ID)

	// Create building
	chinnHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Chinn Hall",
		Location: models.Location{Lat: 40.909830701670295, Lng: -73.12865401298465},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, chinnHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", chinnHall.ID)

	// Create building
	tablerCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Tabler Center",
		Location: models.Location{Lat: 40.91012401804897, Lng: -73.12703847022057},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, tablerCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", tablerCenter.ID)

	// Create building
	handHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Hand Hall",
		Location: models.Location{Lat: 40.90962311590286, Lng: -73.12617153732793},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, handHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", handHall.ID)

	// Create building
	douglassHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Douglass Hall",
		Location: models.Location{Lat: 40.90886110701905, Lng: -73.12664528346346},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, douglassHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", douglassHall.ID)

	// Create building
	dreiserHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Dreiser Hall",
		Location: models.Location{Lat: 40.90864425024119, Lng: -73.12763877868426},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, dreiserHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", dreiserHall.ID)

	// Create building
	gershwinHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Gershwin Hall",
		Location: models.Location{Lat: 40.91141471488535, Lng: -73.12252476879256},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, gershwinHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", gershwinHall.ID)

	// Create building
	hendrixHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Hendrix Hall",
		Location: models.Location{Lat: 40.911592031208016, Lng: -73.12320968659037},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, hendrixHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", hendrixHall.ID)

	// Create building
	mountHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Mount Hall",
		Location: models.Location{Lat: 40.91151880111218, Lng: -73.12449232944562},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, mountHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", mountHall.ID)

	// Create building
	whitmanHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Whitman Hall",
		Location: models.Location{Lat: 40.910701982331055, Lng: -73.12275027110068},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, whitmanHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", whitmanHall.ID)

	// Create building
	cardozoHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Cardozo Hall",
		Location: models.Location{Lat: 40.91087976140076, Lng: -73.12493102295917},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, cardozoHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", cardozoHall.ID)

	// Create building
	rothCafe := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Roth Café",
		Location: models.Location{Lat: 40.91075964045432, Lng: -73.12372302630207},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, rothCafe)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", rothCafe.ID)

	// Create building
	scanCenter := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "SCAN Center",
		Location: models.Location{Lat: 40.91079807917802, Lng: -73.11997823653383},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, scanCenter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", scanCenter.ID)

	// Create building
	yangHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Yang Hall",
		Location: models.Location{Lat: 40.912470112546266, Lng: -73.12907897000453},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, yangHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", yangHall.ID)

	// Create building
	stimsonHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Stimson Hall",
		Location: models.Location{Lat: 40.91193287878291, Lng: -73.12914168664138},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, stimsonHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", stimsonHall.ID)

	// Create building
	greeleyHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Greeley Hall",
		Location: models.Location{Lat: 40.91203442186987, Lng: -73.13106595138132},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, greeleyHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", greeleyHall.ID)

	// Create building
	kellerHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Keller Hall",
		Location: models.Location{Lat: 40.91148192204807, Lng: -73.13001653773449},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, kellerHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", kellerHall.ID)

	// Create building
	wagnerHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Wagner Hall",
		Location: models.Location{Lat: 40.912613352145165, Lng: -73.1307931038551},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, wagnerHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", wagnerHall.ID)

	// Create building
	lauterburHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Lauterbur Hall",
		Location: models.Location{Lat: 40.91268037991284, Lng: -73.12989426818393},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, lauterburHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", lauterburHall.ID)

	// Create building
	nobelHalls := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Nobel Halls",
		Location: models.Location{Lat: 40.91267392421254, Lng: -73.12934241763313},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, nobelHalls)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", nobelHalls.ID)

	// Create building
	rooseveltComm := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Roosevelt Community",
		Location: models.Location{Lat: 40.91191356870777, Lng: -73.13028799322616},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, rooseveltComm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", rooseveltComm.ID)

	// Create building
	deweyHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Dewey Hall",
		Location: models.Location{Lat: 40.912980081023974, Lng: -73.13107559407771},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, deweyHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", deweyHall.ID)

	// Create building
	hamiltonHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Hamilton Hall",
		Location: models.Location{Lat: 40.91389250645616, Lng: -73.13055520096474},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, hamiltonHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", hamiltonHall.ID)

	// Create building
	schickHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Schick Hall",
		Location: models.Location{Lat: 40.91421637231116, Lng: -73.13123212004388},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, schickHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", schickHall.ID)

	// Create building
	eisenhower := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Eisenhower Hall",
		Location: models.Location{Lat: 40.91455101893642, Lng: -73.13189249007522},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, eisenhower)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", eisenhower.ID)

	// Create building
	baruchHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Baruch Hall",
		Location: models.Location{Lat: 40.913934736934195, Lng: -73.13186069217814},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, baruchHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", baruchHall.ID)

	// Create building
	kellyComm := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Kelly Community",
		Location: models.Location{Lat: 40.91370739024136, Lng: -73.13158297586781},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, kellyComm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", kellyComm.ID)

	// Create building
	schomburgApts := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Schomburg Apartments Community",
		Location: models.Location{Lat: 40.91321770006627, Lng: -73.13258127733559},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, schomburgApts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", schomburgApts.ID)

	// Create building
	westAptA := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment A",
		Location: models.Location{Lat: 40.91157790420869, Lng: -73.1325683553833},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptA)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptA.ID)

	// Create building
	westAptB := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment B",
		Location: models.Location{Lat: 40.911750059152645, Lng: -73.13394307705065},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptB.ID)

	// Create building
	westAptC := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment C",
		Location: models.Location{Lat: 40.91194762204886, Lng: -73.13326179018661},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptC.ID)

	// Create building
	westAptD := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment D",
		Location: models.Location{Lat: 40.91225343476852, Lng: -73.13265854285808},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptD)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptD.ID)

	// Create building
	westAptE := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment E",
		Location: models.Location{Lat: 40.912343781717446, Lng: -73.13437905126445},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptE.ID)

	// Create building
	westAptF := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment F",
		Location: models.Location{Lat: 40.91215104362147, Lng: -73.135326005103},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptF)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptF.ID)

	// Create building
	westAptG := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment G",
		Location: models.Location{Lat: 40.91268117237905, Lng: -73.13557578022572},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptG)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptG.ID)

	// Create building
	westAptH := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment H",
		Location: models.Location{Lat: 40.91297033172032, Lng: -73.13487959847946},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptH)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptH.ID)

	// Create building
	westAptI := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment I",
		Location: models.Location{Lat: 40.91340808443563, Lng: -73.1345713653274},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptI)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptI.ID)

	// Create building
	westAptJ := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment J",
		Location: models.Location{Lat: 40.9113118469511, Lng: -73.13535431454514},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptJ)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptJ.ID)

	// Create building
	westAptK := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartment K",
		Location: models.Location{Lat: 40.912360716037014, Lng: -73.13647515059893},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westAptK)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westAptK.ID)

	// Create building
	westApts := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "West Apartments Community",
		Location: models.Location{Lat: 40.91270736994323, Lng: -73.13417103579732},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westApts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westApts.ID)

	// Create building
	nassauHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Nassau Hall",
		Location: models.Location{Lat: 40.90610636622708, Lng: -73.12094730268488},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, nassauHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", nassauHall.ID)

	// Create building
	danaHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Dana Hall",
		Location: models.Location{Lat: 40.90604204530964, Lng: -73.11888075056801},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, danaHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", danaHall.ID)

	// Create building
	putmanHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Putman Hall",
		Location: models.Location{Lat: 40.90599053154775, Lng: -73.1199061035842},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, putmanHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", putmanHall.ID)

	// Create building
	suffolkHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Suffolk Hall",
		Location: models.Location{Lat: 40.90527605392648, Lng: -73.12091960280975},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, suffolkHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", suffolkHall.ID)

	// Create building
	dutchessHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Dutchess Hall",
		Location: models.Location{Lat: 40.90533025143301, Lng: -73.119095650023},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, dutchessHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", dutchessHall.ID)

	// Create building
	challengerHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Challenger Hall",
		Location: models.Location{Lat: 40.90485635622481, Lng: -73.11827657168574},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, challengerHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", challengerHall.ID)

	// Create building
	endeavourHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Endeavour Hall",
		Location: models.Location{Lat: 40.904420801122384, Lng: -73.11780370425126},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, endeavourHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", endeavourHall.ID)

	// Create building
	discoveryHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Discovery Hall",
		Location: models.Location{Lat: 40.904118229029194, Lng: -73.11847496718458},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, discoveryHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", discoveryHall.ID)

	// Create building
	westchesterHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Westchester Hall",
		Location: models.Location{Lat: 40.90361353389892, Lng: -73.12069131602006},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, westchesterHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", westchesterHall.ID)

	// Create building
	sullivanHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Sullivan Hall",
		Location: models.Location{Lat: 40.90281042627816, Lng: -73.12006333971621},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, sullivanHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", sullivanHall.ID)

	// Create building
	rocklandHall := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Rockland Hall",
		Location: models.Location{Lat: 40.90352782461682, Lng: -73.11980219130649},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, rocklandHall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", rocklandHall.ID)

	// Create building
	marineAndAtomspheric := models.Building{
		ID:       primitive.NewObjectID(),
		Name:     "Marine and Atmospheric Sciences",
		Location: models.Location{Lat: 40.90473124741314, Lng: -73.12103309122304},
		Node:     nodeID,
	}

	// Insert building
	_, err = buildingCollection.InsertOne(ctx, marineAndAtomspheric)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted building with ID: %v\n", marineAndAtomspheric.ID)

	// Create parking lot
	parkingLotCollection := database.Collection("parking_lots") //models.Location{Lat: number, Lng: number}
	parkingLot1 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 1",
		Spaces:    676, // Based on provided data
		Faculty:   138,
		Premium:   0,
		Metered:   0,
		Resident:  524,
		Ada:       14,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91656639080932, Lng: -73.1169195914031},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot1)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot1.ID)

	// Create parking lot
	parkingLot2 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 2",
		Spaces:    436, // Based on provided data
		Faculty:   303,
		Premium:   121,
		Metered:   0,
		Resident:  0,
		Ada:       8,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.9168587205784, Lng: -73.11842147620798},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot2)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot2.ID)

	// Create parking lot
	parkingLot3 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 3 Stadium (EV Charging)",
		Spaces:    977, // Based on provided data
		Faculty:   445, // sum
		Premium:   456, // sum
		Metered:   0,
		Resident:  0,
		Ada:       22,   // sum
		Ev:        true, // EV charging available
		Active:    true,
		Location:  models.Location{Lat: 40.91797832262033, Lng: -73.12303664052753},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	_, err = parkingLotCollection.InsertOne(ctx, parkingLot3)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot3.ID)

	// Create parking lot
	parkingLot4 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 4",
		Spaces:    977,  // Based on provided data
		Faculty:   0,    // Not explicitly listed, assuming 0
		Premium:   0,    // "by section" in data, assuming 0 for now
		Metered:   0,    // "by section" in data, assuming 0 for now
		Resident:  0,    // "by section" in data, assuming 0 for now
		Ada:       0,    // "by section" in data, assuming 0 for now
		Ev:        true, // EV Charging available
		Active:    true,
		Location:  models.Location{Lat: 40.917786471918056, Lng: -73.12222810291995},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	_, err = parkingLotCollection.InsertOne(ctx, parkingLot4)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot4.ID)

	// Create parking lot
	parkingLot5 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 5 North P",
		Spaces:    510, // Based on provided data
		Faculty:   0,
		Premium:   363, // sum
		Metered:   30,  // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91942285486984, Lng: -73.12889450386172},
		Node:      nodeID,
		Fee:       2.50, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot5)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot5.ID)

	// Create parking lot
	parkingLot6 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 6 Gym Road",
		Spaces:    370, // Based on provided data
		Faculty:   94,
		Premium:   228, // sum
		Metered:   0,   // sum
		Resident:  0,
		Ada:       8,
		Ev:        true,
		Active:    true,
		Location:  models.Location{Lat: 40.91848563742475, Lng: -73.12798149154817},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot6)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot6.ID)

	// Create parking lot
	parkingLot7A := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 7 ISC Metered Lot (EV Charging)",
		Spaces:    162, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   136, // sum
		Resident:  0,
		Ada:       9,
		Ev:        true,
		Active:    true,
		Location:  models.Location{Lat: 40.91748584051804, Lng: -73.12679441540226},
		Node:      nodeID,
		Fee:       2.50, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot7A)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot7A.ID)

	// Create parking lot
	parkingLot7B := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 7 (Behind IF Arena)",
		Spaces:    23, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91783190861551, Lng: -73.12593599245233},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot7B)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot7B.ID)

	// Create parking lot
	parkingLot9 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 9 COM Business Office/CoGen",
		Spaces:    135, // Based on provided data
		Faculty:   73,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       2,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.916597932623525, Lng: -73.12841197122015},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot9)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot9.ID)

	// Create parking lot
	parkingLot10 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 10 Simons Gated Lot",
		Spaces:    39, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91550435639073, Lng: -73.12707435721556},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot10)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot10.ID)

	// Create parking lot
	parkingLot11 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 11 Math/Physics Lot",
		Spaces:    108, // Based on provided data
		Faculty:   96,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       10,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.915142496986675, Lng: -73.12694676385931},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot11)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot11.ID)

	// Create parking lot
	parkingLot12 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 12 Old H Metered",
		Spaces:    16, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   16, // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91477186582409, Lng: -73.12827782682756},
		Node:      nodeID,
		Fee:       2.50, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot12)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot12.ID)

	// Create parking lot
	parkingLot13 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 13 Old H ",
		Spaces:    547, // Based on provided data
		Faculty:   532,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       15,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91444953670846, Lng: -73.12796205150377},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot13)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot13.ID)

	// Create parking lot
	parkingLot13A := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 13 A  ESS (EV Charging)",
		Spaces:    39, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.914483393850645, Lng: -73.12579283665907},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot13A)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot13A.ID)

	// Create parking lot
	parkingLot14 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 14 Automotive/COM Compound (EV State Only)",
		Spaces:    152, // Based on provided data
		Faculty:   119,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       0,
		Ev:        true,
		Active:    true,
		Location:  models.Location{Lat: 40.91607591672189, Lng: -73.12982682255023},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot14)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot14.ID)

	// Create parking lot
	parkingLot15 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 15 Lower Kelly",
		Spaces:    160, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   0, // sum
		Resident:  143,
		Ada:       6,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.915272171049835, Lng: -73.1312272050479},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot15)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot15.ID)

	// Create parking lot
	parkingLot16A := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 16A",
		Spaces:    42, // Based on provided data
		Faculty:   40,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       2,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91335262764112, Lng: -73.12647231666436},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot16A)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot16A.ID)

	// Create parking lot
	parkingLot16B := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 16B",
		Spaces:    38, // Based on provided data
		Faculty:   38,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.912723903612125, Lng: -73.1266893525137},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot16B)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot16B.ID)

	// Create parking lot
	parkingLot17 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 17  West Apartments",
		Spaces:    1104, // Based on provided data
		Faculty:   18,
		Premium:   0,
		Metered:   0, // sum
		Resident:  957,
		Ada:       55,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.911920221572636, Lng: -73.12860650252043},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot17)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot17.ID)

	// Create parking lot
	parkingLot18 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 18  Heavy Engineering (EV Charging)",
		Spaces:    57, // Based on provided data
		Faculty:   49,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       2,
		Ev:        true,
		Active:    true,
		Location:  models.Location{Lat: 40.91219132196822, Lng: -73.12652367182991},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot18)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot18.ID)

	// Create parking lot
	parkingLot19 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 19",
		Spaces:    59, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   32, // sum
		Resident:  0,
		Ada:       14,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.913624369026046, Lng: -73.12339894704783},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot19)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot19.ID)

	// Create parking lot
	parkingLot20 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 20 Cardozo Metered",
		Spaces:    31, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   30, // sum
		Resident:  0,
		Ada:       1,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.9112515508758, Lng: -73.12558308538264},
		Node:      nodeID,
		Fee:       2.50, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot20)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot20.ID)

	// Create parking lot
	parkingLot21 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 21 Tabler Community",
		Spaces:    363, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   0, // sum
		Resident:  346,
		Ada:       10,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91019299217586, Lng: -73.12566114345968},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot21)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot21.ID)

	// Create parking lot
	parkingLot22 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 22 Tabler Metered (EV Charging)",
		Spaces:    12, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   10, // sum
		Resident:  0,
		Ada:       0,
		Ev:        true,
		Active:    true,
		Location:  models.Location{Lat: 40.90961290794953, Lng: -73.12501065970245},
		Node:      nodeID,
		Fee:       2.50, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot22)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot22.ID)

	// Create parking lot
	parkingLot23 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 23 Tabler",
		Spaces:    342, // Based on provided data
		Faculty:   132,
		Premium:   0,
		Metered:   0, // sum
		Resident:  210,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.90883617702413, Lng: -73.125249170409},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot23)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot23.ID)

	// Create parking lot
	parkingLot24 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 24 Roth/Lake Drive",
		Spaces:    270, // Based on provided data
		Faculty:   251,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       16,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91031249242267, Lng: -73.12257053530992},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot24)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot24.ID)

	// Create parking lot
	parkingLot25 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 25 Lake Drive Metered (EV Charging)",
		Spaces:    20, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   16, // sum
		Resident:  0,
		Ada:       0,
		Ev:        true,
		Active:    true,
		Location:  models.Location{Lat: 40.911568437001485, Lng: -73.12146569709925},
		Node:      nodeID,
		Fee:       2.50, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot25)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot25.ID)

	// Create parking lot
	parkingLot26 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 26 Life Sciences Metered (EV Charging)",
		Spaces:    101, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   89, // sum
		Resident:  0,
		Ada:       9,
		Ev:        true,
		Active:    true,
		Location:  models.Location{Lat: 40.910922625357635, Lng: -73.11954293907264},
		Node:      nodeID,
		Fee:       2.50, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot26)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot26.ID)

	// Create parking lot
	parkingLot27 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 27  Life Sciences Premium A & B",
		Spaces:    24, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91046732436041, Lng: -73.11936775446073},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot27)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot27.ID)

	// Create parking lot
	parkingLot28 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 28  Humanities Metered Admin Special Service",
		Spaces:    49, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   9, // sum
		Resident:  0,
		Ada:       18,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91365436555092, Lng: -73.12018385835721},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot28)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot28.ID)

	// Create parking lot
	parkingLot30 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 30 South Campus Metered",
		Spaces:    5, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   5, // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.904475936328566, Lng: -73.1211684813502},
		Node:      nodeID,
		Fee:       2.50, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot30)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot30.ID)

	// Create parking lot
	parkingLot31 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 31",
		Spaces:    87, // Based on provided data
		Faculty:   55,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       3,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.90408438490092, Lng: -73.12024532343736},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot31)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot31.ID)

	// Create parking lot
	parkingLot32 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 32 SoMAS",
		Spaces:    49, // Based on provided data
		Faculty:   47,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       2,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.903851808137404, Lng: -73.11827436170257},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot32)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot32.ID)

	// Create parking lot
	parkingLot33 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 33 Challenger",
		Spaces:    18, // Based on provided data
		Faculty:   17,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       1,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.90500585186523, Lng: -73.11772514114762},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot33)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot33.ID)

	// Create parking lot
	parkingLot34 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 34 Dana",
		Spaces:    65, // Based on provided data
		Faculty:   40,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       1,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.9064955436956, Lng: -73.11937754308114},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot34)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot34.ID)

	// Create parking lot
	parkingLot35 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 35 Putnam",
		Spaces:    82, // Based on provided data
		Faculty:   52,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       11,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.90636551794018, Lng: -73.11954480594586},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot35)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot35.ID)

	// Create parking lot
	parkingLot36 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 36 Nassau/Suffolk",
		Spaces:    71, // Based on provided data
		Faculty:   30,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       5,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.90585624793648, Lng: -73.12130823450028},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot36)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot36.ID)

	// Create parking lot
	parkingLot37 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 37 Across from Fire Area",
		Spaces:    40, // Based on provided data
		Faculty:   40,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       0,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.905209722211005, Lng: -73.12061050940738},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot37)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot37.ID)

	// Create parking lot
	parkingLot38A := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 38A Dental School",
		Spaces:    250, // Based on provided data
		Faculty:   88,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       11,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.90262100317112, Lng: -73.11966163467063},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot38A)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot38A.ID)

	// Create parking lot
	parkingLot38B := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 38B Side of Dental",
		Spaces:    72, // Based on provided data
		Faculty:   69,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       3,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.90312615942779, Lng: -73.12108498255458},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot38B)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot38B.ID)

	// Create parking lot
	parkingLot40 := models.ParkingLot{
		ID:        primitive.NewObjectID(),
		Name:      "Lot 40 South P",
		Spaces:    2438, // Based on provided data
		Faculty:   0,
		Premium:   0,
		Metered:   0, // sum
		Resident:  0,
		Ada:       34,
		Ev:        false,
		Active:    true,
		Location:  models.Location{Lat: 40.91550435639073, Lng: -73.12707435721556},
		Node:      nodeID,
		Fee:       0.0, // Assuming no fee provided
		Occupancy: 0,
	}

	// Insert parking lot
	_, err = parkingLotCollection.InsertOne(ctx, parkingLot40)
	if err != nil {
		log.Fatal(err)
	}
	// Extract and print the ID of the inserted document
	fmt.Printf("Inserted parking lot with ID: %v\n", parkingLot40.ID)

	// Create vehicle
	vehicleCollection := database.Collection("vehicles")
	vehicleID := primitive.NewObjectID()
	newVehicle := models.Vehicle{
		ID:          vehicleID,
		Name:        "Toyota",
		Model:       "Camry",
		Year:        2020,
		PlateNumber: "SBU-1234",
	}

	// Insert vehicle
	_, err = vehicleCollection.InsertOne(ctx, newVehicle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted vehicle with ID: %v\n", vehicleID)

	// Create user
	userCollection := database.Collection("users")
	userID := primitive.NewObjectID()

	newUser := models.User{
		ID: userID,
		Name: struct {
			First string `bson:"first"`
			Last  string `bson:"last"`
		}{
			First: "John",
			Last:  "Smith",
		},
		Role:    "faculty",
		SbuID:   "123456789",
		Vehicle: vehicleID,
		Address: struct {
			Street  string `bson:"street"`
			City    string `bson:"city"`
			State   string `bson:"state"`
			ZipCode string `bson:"zipCode"`
		}{
			Street:  "123 Campus Drive",
			City:    "Stony Brook",
			State:   "NY",
			ZipCode: "11794",
		},
		DriverLicense: struct {
			Number         string    `bson:"number"`
			State          string    `bson:"state"`
			ExpirationDate time.Time `bson:"expirationDate"`
		}{
			Number:         "DL12345678",
			State:          "NY",
			ExpirationDate: time.Now().AddDate(5, 0, 0), // 5 years from now
		},
		Username:     "jsmith",
		PasswordHash: "$2a$10$dkjfalkdjflakdjfla9u39.adlkfjaldf", // This would be a properly hashed password in production
	}

	// Insert user
	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted user with ID: %v\n", userID)

	// Create reservation
	reservationCollection := database.Collection("reservations")
	reservationID := primitive.NewObjectID()

	newReservation := models.Reservation{
		ID:         reservationID,
		ReservedBy: userID,
		ParkingLot: parkingLotID,
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(2 * time.Hour),
		Status:     "Reserved",
	}

	// Insert reservation
	_, err = reservationCollection.InsertOne(ctx, newReservation)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted reservation with ID: %v\n", reservationID)

	// Create violation
	violationCollection := database.Collection("violations")
	violationID := primitive.NewObjectID()

	newViolation := models.Violation{
		ID:         violationID,
		User:       userID,
		ParkingLot: parkingLotID,
		Reason:     "Parked in unauthorized zone",
		Fine:       50.00,
		PayBy:      time.Now().AddDate(0, 0, 30), // 30 days from now
	}

	// Insert violation
	_, err = violationCollection.InsertOne(ctx, newViolation)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted violation with ID: %v\n", violationID)

	// Create violation rebute
	violationRebuteCollection := database.Collection("violation_rebutes")

	newViolationRebute := models.ViolationRebute{
		ID:         primitive.NewObjectID(),
		User:       userID,
		ParkingLot: parkingLotID,
		Violation:  violationID,
		Reason:     "I had a temporary permit for this zone",
	}

	// Insert violation rebute
	result, err := violationRebuteCollection.InsertOne(ctx, newViolationRebute)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted violation rebute with ID: %v\n", result.InsertedID)

	// Update node with building and parking lot
	updateNode := models.Node{
		ID:         nodeID,
		Type:       "Intersection",
		ParkingLot: parkingLotID,
		Building:   buildingID,
		Neighbors: struct {
			Members []primitive.ObjectID `bson:"members"`
		}{
			Members: []primitive.ObjectID{},
		},
		Location: models.Location{Lat: 40.9123, Lng: -73.1234},
	}

	// Replace the node document
	_, err = nodeCollection.ReplaceOne(ctx, primitive.M{"_id": nodeID}, updateNode)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated node with ID: %v\n", nodeID)

	fmt.Println("\nAll instances created and inserted successfully!")
}
