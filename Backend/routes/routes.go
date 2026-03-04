package routes

import (
	database "github.com/416Coders/P4SBU/Backend/database"
	handler "github.com/416Coders/P4SBU/Backend/handler"
	"github.com/416Coders/P4SBU/Backend/middleware"
	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(router *mux.Router, repo *database.Repository) {
	// Create handler instances
	parkingLotHandler := handler.NewParkingLotHandler(repo)
	buildingHandler := handler.NewBuildingHandler(repo)
	nodeHandler := handler.NewNodeHandler(repo)
	userHandler := handler.NewUserHandler(repo)
	vehicleHandler := handler.NewVehicleHandler(repo)
	reservationHandler := handler.NewReservationHandler(repo)
	violationHandler := handler.NewViolationHandler(repo)
	rebuteHandler := handler.NewViolationRebuteHandler(repo)
	wayfindHandler := handler.NewWayfindHandler(repo)

	// Create auth handler
	authHandler := handler.NewAuthHandler(repo)

	// Public routes (no authentication required)
	publicRouter := router.PathPrefix("").Subrouter()
	registerPublicRoutes(publicRouter, parkingLotHandler, buildingHandler, nodeHandler, userHandler, authHandler, wayfindHandler)

	// Protected routes (authentication required)
	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(repo))
	registerProtectedRoutes(protectedRouter, parkingLotHandler, buildingHandler, nodeHandler, userHandler, vehicleHandler, reservationHandler, violationHandler, rebuteHandler)

	// Admin routes (authentication + admin role required)
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware(repo))
	adminRouter.Use(middleware.RoleMiddleware("admin"))
	registerAdminRoutes(adminRouter, parkingLotHandler, buildingHandler, nodeHandler, userHandler, vehicleHandler, reservationHandler, violationHandler, rebuteHandler)
}

// Register public routes (no authentication)
func registerPublicRoutes(
	router *mux.Router,
	parkingLotHandler *handler.ParkingLotHandler,
	buildingHandler *handler.BuildingHandler,
	nodeHandler *handler.NodeHandler,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	wayfindHandler *handler.WayfindHandler, // Added wayfind handler parameter
) {
	// Authentication
	router.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")

	// Parking Lot routes
	router.HandleFunc("/parking-lots", parkingLotHandler.GetAllParkingLots).Methods("GET")
	router.HandleFunc("/parking-lots/{id}", parkingLotHandler.GetParkingLot).Methods("GET")

	// Building routes
	router.HandleFunc("/buildings", buildingHandler.GetAllBuildings).Methods("GET")
	router.HandleFunc("/buildings/{id}", buildingHandler.GetBuilding).Methods("GET")

	// Node routes
	router.HandleFunc("/nodes", nodeHandler.GetAllNodes).Methods("GET")
	router.HandleFunc("/nodes/{id}", nodeHandler.GetNode).Methods("GET")
	router.HandleFunc("/navigate", nodeHandler.GetNavigationPath).Methods("GET")

	// User routes
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST") // Registration endpoint

	// Wayfind routes
	router.HandleFunc("/wayfind/building/{id}/nearest-parking", wayfindHandler.FindNearestParkingLot).Methods("GET")
	router.HandleFunc("/wayfind/nearest-parking", wayfindHandler.FindNearestParkingLotByCoordinates).Methods("GET")
}

// Register protected routes (authentication required)
func registerProtectedRoutes(
	router *mux.Router,
	parkingLotHandler *handler.ParkingLotHandler,
	buildingHandler *handler.BuildingHandler,
	nodeHandler *handler.NodeHandler,
	userHandler *handler.UserHandler,
	vehicleHandler *handler.VehicleHandler,
	reservationHandler *handler.ReservationHandler,
	violationHandler *handler.ViolationHandler,
	rebuteHandler *handler.ViolationRebuteHandler,
) {
	// User routes - user can only access their own data
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")

	// Vehicle routes
	router.HandleFunc("/vehicles/{id}", vehicleHandler.GetVehicle).Methods("GET")
	router.HandleFunc("/users/{id}/vehicles", vehicleHandler.GetUserVehicles).Methods("GET")
	router.HandleFunc("/vehicles", vehicleHandler.CreateVehicle).Methods("POST")
	router.HandleFunc("/vehicles/{id}", vehicleHandler.UpdateVehicle).Methods("PUT")
	router.HandleFunc("/vehicles/{id}", vehicleHandler.DeleteVehicle).Methods("DELETE")

	// Reservation routes
	router.HandleFunc("/reservations/{id}", reservationHandler.GetReservation).Methods("GET")
	router.HandleFunc("/users/{id}/reservations", reservationHandler.GetUserReservations).Methods("GET")
	router.HandleFunc("/reservations", reservationHandler.CreateReservation).Methods("POST")
	router.HandleFunc("/reservations/{id}", reservationHandler.UpdateReservation).Methods("PUT")
	router.HandleFunc("/reservations/{id}", reservationHandler.CancelReservation).Methods("DELETE")

	// Violation routes
	router.HandleFunc("/violations/{id}", violationHandler.GetViolation).Methods("GET")
	router.HandleFunc("/users/{id}/violations", violationHandler.GetUserViolations).Methods("GET")

	// Violation Rebute routes
	router.HandleFunc("/violation-rebutes/{id}", rebuteHandler.GetRebute).Methods("GET")
	router.HandleFunc("/users/{id}/violation-rebutes", rebuteHandler.GetUserRebutes).Methods("GET")
	router.HandleFunc("/violations/{id}/rebutes", rebuteHandler.GetViolationRebutes).Methods("GET")
	router.HandleFunc("/violation-rebutes", rebuteHandler.CreateRebute).Methods("POST")
}

// Register admin routes (admin role required)
func registerAdminRoutes(
	router *mux.Router,
	parkingLotHandler *handler.ParkingLotHandler,
	buildingHandler *handler.BuildingHandler,
	nodeHandler *handler.NodeHandler,
	userHandler *handler.UserHandler,
	vehicleHandler *handler.VehicleHandler,
	reservationHandler *handler.ReservationHandler,
	violationHandler *handler.ViolationHandler,
	rebuteHandler *handler.ViolationRebuteHandler,
) {
	// User routes (admins can manage all users)
	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// ParkingLot admin routes
	router.HandleFunc("/parking-lots", parkingLotHandler.CreateParkingLot).Methods("POST")
	router.HandleFunc("/parking-lots/{id}", parkingLotHandler.UpdateParkingLot).Methods("PUT")
	router.HandleFunc("/parking-lots/{id}", parkingLotHandler.DeleteParkingLot).Methods("DELETE")

	// Building admin routes
	router.HandleFunc("/buildings", buildingHandler.CreateBuilding).Methods("POST")
	router.HandleFunc("/buildings/{id}", buildingHandler.UpdateBuilding).Methods("PUT")
	router.HandleFunc("/buildings/{id}", buildingHandler.DeleteBuilding).Methods("DELETE")

	// Node admin routes
	router.HandleFunc("/nodes", nodeHandler.CreateNode).Methods("POST")
	router.HandleFunc("/nodes/{id}", nodeHandler.UpdateNode).Methods("PUT")
	router.HandleFunc("/nodes/{id}", nodeHandler.DeleteNode).Methods("DELETE")

	// Reservation admin routes
	router.HandleFunc("/reservations", reservationHandler.GetAllReservations).Methods("GET")
	router.HandleFunc("/parking-lots/{id}/reservations", reservationHandler.GetParkingLotReservations).Methods("GET")

	// Violation admin routes
	router.HandleFunc("/violations", violationHandler.GetAllViolations).Methods("GET")
	router.HandleFunc("/violations", violationHandler.CreateViolation).Methods("POST")
	router.HandleFunc("/violations/{id}", violationHandler.UpdateViolation).Methods("PUT")
	router.HandleFunc("/violations/{id}", violationHandler.DeleteViolation).Methods("DELETE")

	// Violation Rebute admin routes
	router.HandleFunc("/violation-rebutes", rebuteHandler.GetAllRebutes).Methods("GET")
	router.HandleFunc("/violation-rebutes/{id}", rebuteHandler.UpdateRebute).Methods("PUT")
	router.HandleFunc("/violation-rebutes/{id}", rebuteHandler.DeleteRebute).Methods("DELETE")
}
