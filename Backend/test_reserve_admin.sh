#!/bin/bash

# P4SBU - Test Reservation System and Admin Routes
# This script specifically tests the reservation endpoints and admin privileged routes
# Author: Claude (Based on example test_api.sh)

# Configuration
API_URL="http://localhost:8080/api/v1"
USER_CREDENTIALS_FILE=".user_credentials.json"
ADMIN_CREDENTIALS_FILE=".admin_credentials.json"
SAVE_ORIGINAL_DATA=true     # Save original data for later restoration
TEST_DATA_FILE=".test_data.json"

# Colors for terminal output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Check if jq is installed for JSON formatting
if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}jq is not installed. Install it for better JSON formatting:${NC}"
    echo -e "${YELLOW}  sudo apt-get install jq (Debian/Ubuntu)${NC}"
    echo -e "${YELLOW}  brew install jq (macOS)${NC}"
    exit 1
fi

# Check if credentials files exist
if [ ! -f "$USER_CREDENTIALS_FILE" ] || [ ! -f "$ADMIN_CREDENTIALS_FILE" ]; then
    echo -e "${RED}One or both credential files not found:${NC}"
    echo -e "- User: $USER_CREDENTIALS_FILE"
    echo -e "- Admin: $ADMIN_CREDENTIALS_FILE"
    echo -e "${YELLOW}Please run the test_account.sh script first to create test credentials.${NC}"
    exit 1
fi

# Load credentials - using Basic Auth
USER_ID=$(cat "$USER_CREDENTIALS_FILE" | jq -r '.userId')
USER_USERNAME=$(cat "$USER_CREDENTIALS_FILE" | jq -r '.username')
USER_PASSWORD=$(cat "$USER_CREDENTIALS_FILE" | jq -r '.password')
USER_ROLE=$(cat "$USER_CREDENTIALS_FILE" | jq -r '.role')

ADMIN_ID=$(cat "$ADMIN_CREDENTIALS_FILE" | jq -r '.userId')
ADMIN_USERNAME=$(cat "$ADMIN_CREDENTIALS_FILE" | jq -r '.username')
ADMIN_PASSWORD=$(cat "$ADMIN_CREDENTIALS_FILE" | jq -r '.password')
ADMIN_ROLE=$(cat "$ADMIN_CREDENTIALS_FILE" | jq -r '.role')

if [ -z "$USER_ID" ] || [ "$USER_ID" == "null" ] || [ -z "$ADMIN_ID" ] || [ "$ADMIN_ID" == "null" ]; then
    echo -e "${RED}Invalid user ID in one or both credentials files.${NC}"
    echo -e "${YELLOW}Please run the test_account.sh script again to create valid credentials.${NC}"
    exit 1
fi

echo -e "${GREEN}Loaded user credentials for: $USER_USERNAME (ID: $USER_ID, Role: $USER_ROLE)${NC}"
echo -e "${GREEN}Loaded admin credentials for: $ADMIN_USERNAME (ID: $ADMIN_ID, Role: $ADMIN_ROLE)${NC}"

# Function to print section headers
print_header() {
    echo -e "\n${BLUE}======================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}======================================================${NC}"
}

# Function to print test name
print_test() {
    echo -e "\n${CYAN}>>> Testing: $1${NC}"
}

# Function to format and print JSON response
print_response() {
    local response="$1"
    local status="$2"
    local endpoint="$3"
    local id_path="$4" # Optional: path to extract ID from response
    
    echo -e "${PURPLE}Status: $status${NC}"
    
    if [ "$status" -eq 200 ] || [ "$status" -eq 201 ]; then
        echo -e "${GREEN}Response:${NC}"
        # Try to parse as JSON, if it fails, print as plain text
        if echo "$response" | jq '.' &>/dev/null; then
            echo "$response" | jq '.'
        else
            echo "$response"
        fi
        
        # Extract ID if specified
        if [ -n "$id_path" ]; then
            local extracted_id=$(echo "$response" | jq -r "$id_path" 2>/dev/null)
            if [ -n "$extracted_id" ] && [ "$extracted_id" != "null" ]; then
                if [[ $endpoint == */parking-lots ]]; then
                    PARKING_LOT_ID=$extracted_id
                    echo -e "${GREEN}Extracted Parking Lot ID: $PARKING_LOT_ID${NC}"
                elif [[ $endpoint == */buildings ]]; then
                    BUILDING_ID=$extracted_id
                    echo -e "${GREEN}Extracted Building ID: $BUILDING_ID${NC}"
                elif [[ $endpoint == */reservations ]]; then
                    RESERVATION_ID=$extracted_id
                    echo -e "${GREEN}Extracted Reservation ID: $RESERVATION_ID${NC}"
                elif [[ $endpoint == */vehicles ]]; then
                    VEHICLE_ID=$extracted_id
                    echo -e "${GREEN}Extracted Vehicle ID: $VEHICLE_ID${NC}"
                fi
            fi
        fi
    else
        echo -e "${RED}Response:${NC} $response"
    fi
}

# Function to run a test with user authentication
run_user_test() {
    local method="$1"
    local endpoint="$2"
    local name="$3"
    local data="$4"
    local id_path="$5" # Optional: path to extract ID from response
    
    print_test "$name (as Regular User)"
    
    local response
    local status
    
    if [ "$method" == "GET" ]; then
        # GET request with auth - using Basic Auth
        response=$(curl -s -u "$USER_USERNAME:$USER_PASSWORD" "$API_URL$endpoint")
        status=$?
        
        # Check if curl failed
        if [ $status -ne 0 ]; then
            echo -e "${RED}Curl request failed with status $status${NC}"
            return 1
        fi
        
        # Get HTTP status code with a separate request
        status=$(curl -s -o /dev/null -u "$USER_USERNAME:$USER_PASSWORD" -w "%{http_code}" "$API_URL$endpoint")
    else
        # POST, PUT, DELETE with data and auth
        if [ -n "$data" ]; then
            response=$(curl -s -H "Content-Type: application/json" -u "$USER_USERNAME:$USER_PASSWORD" -d "$data" -X "$method" "$API_URL$endpoint")
            status=$?
            
            # Check if curl failed
            if [ $status -ne 0 ]; then
                echo -e "${RED}Curl request failed with status $status${NC}"
                return 1
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -H "Content-Type: application/json" -u "$USER_USERNAME:$USER_PASSWORD" -d "$data" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        else
            response=$(curl -s -u "$USER_USERNAME:$USER_PASSWORD" -X "$method" "$API_URL$endpoint")
            status=$?
            
            # Check if curl failed
            if [ $status -ne 0 ]; then
                echo -e "${RED}Curl request failed with status $status${NC}"
                return 1
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -u "$USER_USERNAME:$USER_PASSWORD" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        fi
    fi
    
    # Print formatted output
    print_response "$response" "$status" "$endpoint" "$id_path"
    return 0
}

# Function to run a test with admin authentication
run_admin_test() {
    local method="$1"
    local endpoint="$2"
    local name="$3"
    local data="$4"
    local id_path="$5" # Optional: path to extract ID from response
    
    print_test "$name (as Admin)"
    
    local response
    local status
    
    if [ "$method" == "GET" ]; then
        # GET request with auth - using Basic Auth
        response=$(curl -s -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" "$API_URL$endpoint")
        status=$?
        
        # Check if curl failed
        if [ $status -ne 0 ]; then
            echo -e "${RED}Curl request failed with status $status${NC}"
            return 1
        fi
        
        # Get HTTP status code with a separate request
        status=$(curl -s -o /dev/null -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" -w "%{http_code}" "$API_URL$endpoint")
    else
        # POST, PUT, DELETE with data and auth
        if [ -n "$data" ]; then
            response=$(curl -s -H "Content-Type: application/json" -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" -d "$data" -X "$method" "$API_URL$endpoint")
            status=$?
            
            # Check if curl failed
            if [ $status -ne 0 ]; then
                echo -e "${RED}Curl request failed with status $status${NC}"
                return 1
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -H "Content-Type: application/json" -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" -d "$data" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        else
            response=$(curl -s -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" -X "$method" "$API_URL$endpoint")
            status=$?
            
            # Check if curl failed
            if [ $status -ne 0 ]; then
                echo -e "${RED}Curl request failed with status $status${NC}"
                return 1
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        fi
    fi
    
    # Print formatted output
    print_response "$response" "$status" "$endpoint" "$id_path"
    return 0
}

# Function to save original data for later restoration
save_original_data() {
    local entity="$1"
    local id="$2"
    local auth="$3"  # "user" or "admin"
    
    if [ "$SAVE_ORIGINAL_DATA" != true ]; then
        return
    fi

    if [ -z "$id" ]; then
        echo -e "${YELLOW}>>> Skipping save original data for $entity (no ID provided)${NC}"
        return
    fi
    
    echo -e "\n${CYAN}>>> Saving original $entity data (ID: $id)${NC}"
    
    local username="$USER_USERNAME"
    local password="$USER_PASSWORD"
    if [ "$auth" == "admin" ]; then
        username="$ADMIN_USERNAME"
        password="$ADMIN_PASSWORD"
    fi
    
    local data=$(curl -s -u "$username:$password" "$API_URL/$entity/$id")
    
    # Create test data file if it doesn't exist
    if [ ! -f "$TEST_DATA_FILE" ]; then
        echo "{}" > "$TEST_DATA_FILE"
    fi
    
    # Save data to test data file
    local existing_data=$(cat "$TEST_DATA_FILE")
    echo "$existing_data" | jq --arg entity "${entity}_${id}" --arg data "$data" '. + {($entity): $data}' > "$TEST_DATA_FILE"
    
    echo -e "${GREEN}Saved original $entity data to $TEST_DATA_FILE${NC}"
}

# Function to restore original data
restore_original_data() {
    local entity="$1"
    local id="$2"
    local auth="$3"  # "user" or "admin"
    
    if [ "$SAVE_ORIGINAL_DATA" != true ]; then
        return
    fi

    if [ -z "$id" ] || [ ! -f "$TEST_DATA_FILE" ]; then
        echo -e "${YELLOW}>>> Skipping restore original data for $entity (no ID provided or no test data file)${NC}"
        return
    fi
    
    echo -e "\n${CYAN}>>> Restoring original $entity data (ID: $id)${NC}"
    
    local username="$USER_USERNAME"
    local password="$USER_PASSWORD"
    local admin_prefix=""
    if [ "$auth" == "admin" ]; then
        username="$ADMIN_USERNAME"
        password="$ADMIN_PASSWORD"
        admin_prefix="/admin"
    fi
    
    local data=$(cat "$TEST_DATA_FILE" | jq -r --arg entity "${entity}_${id}" '.[$entity]')
    
    if [ "$data" == "null" ]; then
        echo -e "${YELLOW}No original data found for $entity with ID $id${NC}"
        return
    fi
    
    # Update entity with original data
    local status=$(curl -s -o /dev/null -H "Content-Type: application/json" -u "$username:$password" -d "$data" -w "%{http_code}" -X "PUT" "$API_URL$admin_prefix/$entity/$id")
    
    if [ "$status" -eq 200 ]; then
        echo -e "${GREEN}Successfully restored original data for $entity with ID $id${NC}"
    else
        echo -e "${RED}Failed to restore original data for $entity with ID $id (Status: $status)${NC}"
    fi
}

# Function to setup test data
setup_test_data() {
    print_header "SETTING UP TEST DATA"
    
    # Create a vehicle for reservations if needed
    VEHICLE_DATA='{
        "name": "Test Reservation Vehicle",
        "model": "Honda Accord",
        "year": 2022,
        "plateNumber": "RSV-123"
    }'
    
    run_user_test "POST" "/vehicles" "Create Test Vehicle for Reservations" "$VEHICLE_DATA" ".id"
    USER_VEHICLE_ID=$VEHICLE_ID
    
    # Get the first available parking lot for testing
    print_test "Get a parking lot for testing"
    
    # Get parking lots - using Basic Auth
    PARKING_LOTS_RESPONSE=$(curl -s -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" "$API_URL/parking-lots")
    
    # Try to parse the response
    if echo "$PARKING_LOTS_RESPONSE" | jq . &>/dev/null; then
        PARKING_LOT_COUNT=$(echo "$PARKING_LOTS_RESPONSE" | jq '. | length' 2>/dev/null || echo 0)
        
        if [ "$PARKING_LOT_COUNT" -gt 0 ]; then
            PARKING_LOT_ID=$(echo "$PARKING_LOTS_RESPONSE" | jq -r '.[0]._id // .[0].id // .[0].ID' 2>/dev/null)
            PARKING_LOT_NAME=$(echo "$PARKING_LOTS_RESPONSE" | jq -r '.[0].name // .[0].Name' 2>/dev/null)
            echo -e "${GREEN}Found parking lot: $PARKING_LOT_NAME (ID: $PARKING_LOT_ID)${NC}"
        else
            echo -e "${YELLOW}No parking lots found, creating one${NC}"
            NEW_PARKING_LOT='{
                "name": "Test Reservation Lot",
                "spaces": 50,
                "faculty": 20,
                "premium": 5,
                "metered": 5,
                "resident": 15,
                "ada": 5,
                "ev": true,
                "active": true,
                "location": {
                    "lat": 40.91,
                    "lng": -73.12
                }
            }'
            
            run_admin_test "POST" "/admin/parking-lots" "Create Test Parking Lot" "$NEW_PARKING_LOT" ".id"
            PARKING_LOT_NAME="Test Reservation Lot"
        fi
    else
        echo -e "${RED}Failed to get parking lots. Response: $PARKING_LOTS_RESPONSE${NC}"
        # Try creating a parking lot anyway
        echo -e "${YELLOW}Creating a test parking lot${NC}"
        NEW_PARKING_LOT='{
            "name": "Test Reservation Lot",
            "spaces": 50,
            "faculty": 20,
            "premium": 5,
            "metered": 5,
            "resident": 15,
            "ada": 5,
            "ev": true,
            "active": true,
            "location": {
                "lat": 40.91,
                "lng": -73.12
            }
        }'
        
        run_admin_test "POST" "/admin/parking-lots" "Create Test Parking Lot" "$NEW_PARKING_LOT" ".id"
        PARKING_LOT_NAME="Test Reservation Lot"
    fi
    
    if [ -z "$PARKING_LOT_ID" ]; then
        echo -e "${RED}Failed to get or create a parking lot for testing.${NC}"
        exit 1
    fi
    
    # Save original parking lot data
    save_original_data "parking-lots" "$PARKING_LOT_ID" "admin"
}

# Start testing
print_header "CAMPUS PARKING API - RESERVATION & ADMIN ROUTES TEST"

# Setup test data
setup_test_data

# --------------------------
# RESERVATION SYSTEM TESTS
# --------------------------
print_header "RESERVATION SYSTEM TESTS"

# Current time plus 1 hour for start time
START_TIME=$(date -u -d "+1 hour" "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+1H "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u "+%Y-%m-%dT%H:%M:%SZ")
# Current time plus 3 hours for end time
END_TIME=$(date -u -d "+3 hour" "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+3H "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u "+%Y-%m-%dT%H:%M:%SZ")

# 1. Create Reservation
print_header "1. CREATE RESERVATION"
RESERVATION_DATA='{
    "reserved_by": "'$USER_ID'",
    "parking_lot": "'$PARKING_LOT_ID'",
    "start_time": "'$START_TIME'",
    "end_time": "'$END_TIME'",
    "status": "confirmed"
}'

run_user_test "POST" "/reservations" "Create New Reservation" "$RESERVATION_DATA" ".id"
USER_RESERVATION_ID=$RESERVATION_ID

if [ -z "$USER_RESERVATION_ID" ]; then
    echo -e "${RED}Failed to create reservation.${NC}"
    exit 1
fi

# 2. Get Specific Reservation
print_header "2. GET SPECIFIC RESERVATION"
run_user_test "GET" "/reservations/$USER_RESERVATION_ID" "Get Reservation by ID"

# 3. Get All User Reservations
print_header "3. GET ALL USER RESERVATIONS"
run_user_test "GET" "/users/$USER_ID/reservations" "Get All User Reservations"

# 4. Update Reservation
print_header "4. UPDATE RESERVATION"
# New end time (4 hours from now instead of 3)
UPDATED_END_TIME=$(date -u -d "+4 hour" "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+4H "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u "+%Y-%m-%dT%H:%M:%SZ")

UPDATED_RESERVATION_DATA='{
    "reserved_by": "'$USER_ID'",
    "parking_lot": "'$PARKING_LOT_ID'",
    "start_time": "'$START_TIME'",
    "end_time": "'$UPDATED_END_TIME'",
    "status": "confirmed"
}'

run_user_test "PUT" "/reservations/$USER_RESERVATION_ID" "Update Reservation (Change End Time)" "$UPDATED_RESERVATION_DATA"

# --------------------------
# ADMIN ROUTES TESTS
# --------------------------
print_header "ADMIN ROUTES TESTS"

# 1. Admin - Get All Reservations
print_header "1. ADMIN - GET ALL RESERVATIONS"
run_admin_test "GET" "/admin/reservations" "Get All Reservations"

# 2. Admin - Get Parking Lot Reservations
print_header "2. ADMIN - GET PARKING LOT RESERVATIONS"
run_admin_test "GET" "/admin/parking-lots/$PARKING_LOT_ID/reservations" "Get Parking Lot Reservations"

# 3. Admin - Edit a parking lot
print_header "3. ADMIN - EDIT PARKING LOT"

# Get the current parking lot data
PARKING_LOT_RESPONSE=$(curl -s -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" "$API_URL/parking-lots/$PARKING_LOT_ID")

# Check if we got valid JSON
if echo "$PARKING_LOT_RESPONSE" | jq . &>/dev/null; then
    # Create updated parking lot data
    UPDATED_PARKING_LOT="{
        \"name\": \"$PARKING_LOT_NAME (EDITED)\",
        \"spaces\": 60,
        \"faculty\": 25,
        \"premium\": 10,
        \"metered\": 5,
        \"resident\": 15,
        \"ada\": 5,
        \"ev\": true,
        \"active\": true,
        \"location\": $(echo "$PARKING_LOT_RESPONSE" | jq '.location' 2>/dev/null || echo '{\"lat\": 40.91, \"lng\": -73.12}')
    }"
    
    run_admin_test "PUT" "/admin/parking-lots/$PARKING_LOT_ID" "Update Parking Lot" "$UPDATED_PARKING_LOT"
else
    echo -e "${RED}Failed to get parking lot data. Skipping update test.${NC}"
    echo -e "${RED}Response: $PARKING_LOT_RESPONSE${NC}"
fi

# 4. Admin - Create a new building
print_header "4. ADMIN - CREATE NEW BUILDING"
NEW_BUILDING='{
    "name": "Test Admin Building",
    "location": {
        "lat": 40.92,
        "lng": -73.13
    },
    "node": "000000000000000000000000"
}'

run_admin_test "POST" "/admin/buildings" "Create Building" "$NEW_BUILDING" ".id"
ADMIN_BUILDING_ID=$BUILDING_ID

# 5. Admin - Update the building
if [ -n "$ADMIN_BUILDING_ID" ]; then
    print_header "5. ADMIN - UPDATE BUILDING"
    UPDATED_BUILDING='{
        "name": "Updated Test Admin Building",
        "location": {
            "lat": 40.92,
            "lng": -73.13
        },
        "node": "000000000000000000000000"
    }'
    
    run_admin_test "PUT" "/admin/buildings/$ADMIN_BUILDING_ID" "Update Building" "$UPDATED_BUILDING"
else
    echo -e "${YELLOW}>>> Skipping building update test (no building created)${NC}"
fi

# --------------------------
# Permission Boundary Tests
# --------------------------
print_header "PERMISSION BOUNDARY TESTS"

# 1. Regular user tries to access admin routes
print_header "1. USER ATTEMPTS TO ACCESS ADMIN ROUTES (SHOULD FAIL)"
run_user_test "GET" "/admin/reservations" "Regular User Tries to Get All Reservations"
run_user_test "GET" "/admin/parking-lots/$PARKING_LOT_ID/reservations" "Regular User Tries to Get Parking Lot Reservations"

# 2. Regular user tries to access another user's reservations
if [ -n "$ADMIN_ID" ] && [ "$ADMIN_ID" != "$USER_ID" ]; then
    print_header "2. USER ATTEMPTS TO ACCESS ANOTHER USER'S RESERVATIONS (SHOULD FAIL)"
    run_user_test "GET" "/users/$ADMIN_ID/reservations" "Regular User Tries to Access Admin's Reservations"
fi

# --------------------------
# CLEANUP
# --------------------------
print_header "CLEANING UP"

# 1. Cancel the reservation
if [ -n "$USER_RESERVATION_ID" ]; then
    print_test "Cancel User Reservation"
    run_user_test "DELETE" "/reservations/$USER_RESERVATION_ID" "Cancel User Reservation"
fi

# 2. Delete the vehicle
if [ -n "$USER_VEHICLE_ID" ]; then
    print_test "Delete User Vehicle"
    run_user_test "DELETE" "/vehicles/$USER_VEHICLE_ID" "Delete User Vehicle"
fi

# 3. Delete the building created by admin
if [ -n "$ADMIN_BUILDING_ID" ]; then
    print_test "Delete Admin Building"
    run_admin_test "DELETE" "/admin/buildings/$ADMIN_BUILDING_ID" "Delete Admin Building"
fi

# 4. Restore original parking lot data
if [ -n "$PARKING_LOT_ID" ]; then
    print_test "Restore Original Parking Lot Data"
    restore_original_data "parking-lots" "$PARKING_LOT_ID" "admin"
else
    echo -e "${YELLOW}>>> Skipping parking lot restoration (no ID available)${NC}"
fi

# Remove test data file if it exists
if [ -f "$TEST_DATA_FILE" ]; then
    rm "$TEST_DATA_FILE"
    echo -e "${GREEN}Removed test data file.${NC}"
fi

echo -e "\n${GREEN}Testing Complete!${NC}"