#!/bin/bash

# Comprehensive Campus Parking API Testing Script
# Tests all API endpoints including the newly added ones

# Configuration
API_URL="http://localhost:8080/api/v1"
USER_CREDENTIALS_FILE=".user_credentials.json"
ADMIN_CREDENTIALS_FILE=".admin_credentials.json"
LIMIT_OUTPUT=true     # Limit JSON array output to a few items
MAX_ITEMS=3           # Maximum number of items to show from arrays
SAVE_IDS=true         # Save IDs for use in subsequent tests

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

        # Check if response is an array
        if [[ $response == \[* ]] && [ "$LIMIT_OUTPUT" = true ]; then
            # Check if it's valid JSON before processing
            if echo "$response" | jq . &>/dev/null; then
                # Count items in array
                local count=$(echo "$response" | jq length)
                echo -e "${YELLOW}Found $count items. Showing first $MAX_ITEMS:${NC}"
                
                # Show just the first few items
                echo "$response" | jq ".[:$MAX_ITEMS]"
                
                if [ "$count" -gt "$MAX_ITEMS" ]; then
                    echo -e "${YELLOW}...and $(($count - $MAX_ITEMS)) more items (not shown for brevity)${NC}"
                fi
                
                # Extract ID if needed and endpoint matches
                if [ "$SAVE_IDS" = true ] && [ "$endpoint" = "/parking-lots" ]; then
                    PARKING_LOT_ID=$(echo "$response" | jq -r '.[0]._id // .[0].ID // .[0].id')
                    echo -e "${GREEN}Extracted Parking Lot ID: $PARKING_LOT_ID${NC}"
                elif [ "$SAVE_IDS" = true ] && [ "$endpoint" = "/buildings" ]; then
                    BUILDING_ID=$(echo "$response" | jq -r '.[0]._id // .[0].ID // .[0].id')
                    echo -e "${GREEN}Extracted Building ID: $BUILDING_ID${NC}"
                elif [ "$SAVE_IDS" = true ] && [ "$endpoint" = "/nodes" ]; then
                    NODE_ID=$(echo "$response" | jq -r '.[0]._id // .[0].ID // .[0].id')
                    echo -e "${GREEN}Extracted Node ID: $NODE_ID${NC}"
                elif [ "$SAVE_IDS" = true ] && [ "$endpoint" = "/vehicles" ]; then
                    VEHICLE_ID=$(echo "$response" | jq -r '.[0]._id // .[0].ID // .[0].id')
                    echo -e "${GREEN}Extracted Vehicle ID: $VEHICLE_ID${NC}"
                elif [ "$SAVE_IDS" = true ] && [ "$endpoint" = "/reservations" ]; then
                    RESERVATION_ID=$(echo "$response" | jq -r '.[0]._id // .[0].ID // .[0].id')
                    echo -e "${GREEN}Extracted Reservation ID: $RESERVATION_ID${NC}"
                elif [ "$SAVE_IDS" = true ] && [ "$endpoint" = "/violations" ]; then
                    VIOLATION_ID=$(echo "$response" | jq -r '.[0]._id // .[0].ID // .[0].id')
                    echo -e "${GREEN}Extracted Violation ID: $VIOLATION_ID${NC}"
                elif [ "$SAVE_IDS" = true ] && [ "$endpoint" = "/violation-rebutes" ]; then
                    REBUTE_ID=$(echo "$response" | jq -r '.[0]._id // .[0].ID // .[0].id')
                    echo -e "${GREEN}Extracted Rebute ID: $REBUTE_ID${NC}"
                fi
            else
                # Not valid JSON, just print as-is
                echo "$response"
            fi
        else
            # Not an array or limiting disabled, show full response
            if echo "$response" | jq . &>/dev/null; then
                echo "$response" | jq
                
                # Extract ID if specified
                if [ -n "$id_path" ] && [ "$SAVE_IDS" = true ]; then
                    local extracted_id=$(echo "$response" | jq -r "$id_path")
                    if [[ $endpoint == */users ]]; then
                        CREATED_USER_ID=$extracted_id
                        echo -e "${GREEN}Extracted Created User ID: $CREATED_USER_ID${NC}"
                    elif [[ $endpoint == */vehicles ]]; then
                        VEHICLE_ID=$extracted_id
                        echo -e "${GREEN}Extracted Vehicle ID: $VEHICLE_ID${NC}"
                    elif [[ $endpoint == */parking-lots ]]; then
                        PARKING_LOT_ID=$extracted_id
                        echo -e "${GREEN}Extracted Parking Lot ID: $PARKING_LOT_ID${NC}"
                    elif [[ $endpoint == */buildings ]]; then
                        BUILDING_ID=$extracted_id
                        echo -e "${GREEN}Extracted Building ID: $BUILDING_ID${NC}"
                    elif [[ $endpoint == */nodes ]]; then
                        NODE_ID=$extracted_id
                        echo -e "${GREEN}Extracted Node ID: $NODE_ID${NC}"
                    elif [[ $endpoint == */reservations ]]; then
                        RESERVATION_ID=$extracted_id
                        echo -e "${GREEN}Extracted Reservation ID: $RESERVATION_ID${NC}"
                    elif [[ $endpoint == */violations ]]; then
                        VIOLATION_ID=$extracted_id
                        echo -e "${GREEN}Extracted Violation ID: $VIOLATION_ID${NC}"
                    elif [[ $endpoint == */violation-rebutes ]]; then
                        REBUTE_ID=$extracted_id
                        echo -e "${GREEN}Extracted Rebute ID: $REBUTE_ID${NC}"
                    fi
                fi
            else
                # Not valid JSON, just print as-is
                echo "$response"
            fi
        fi
    else
        echo -e "${RED}Response:${NC} $response"
    fi
}

# Function to run a test without authentication
run_public_test() {
    local method="$1"
    local endpoint="$2"
    local name="$3"
    local data="$4"
    local id_path="$5" # Optional: path to extract ID from response
    
    print_test "$name"
    
    local response
    local status
    
    if [ "$method" == "GET" ]; then
        # GET request
        response=$(curl -s "$API_URL$endpoint")
        status=$?
        
        # Check if curl failed
        if [ $status -ne 0 ]; then
            echo -e "${RED}Curl request failed with status $status${NC}"
            return
        fi
        
        # Get HTTP status code with a separate request
        status=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL$endpoint")
    else
        # POST, PUT, DELETE with data
        if [ -n "$data" ]; then
            response=$(curl -s -H "Content-Type: application/json" -d "$data" -X "$method" "$API_URL$endpoint")
            status=$?
            
            # Check if curl failed
            if [ $status -ne 0 ]; then
                echo -e "${RED}Curl request failed with status $status${NC}"
                return
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -H "Content-Type: application/json" -d "$data" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        else
            response=$(curl -s -X "$method" "$API_URL$endpoint")
            status=$?
            
            # Check if curl failed
            if [ $status -ne 0 ]; then
                echo -e "${RED}Curl request failed with status $status${NC}"
                return
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        fi
    fi
    
    # Print formatted output
    print_response "$response" "$status" "$endpoint" "$id_path"
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
            return
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
                return
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -H "Content-Type: application/json" -u "$USER_USERNAME:$USER_PASSWORD" -d "$data" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        else
            response=$(curl -s -u "$USER_USERNAME:$USER_PASSWORD" -X "$method" "$API_URL$endpoint")
            status=$?
            
            # Check if curl failed
            if [ $status -ne 0 ]; then
                echo -e "${RED}Curl request failed with status $status${NC}"
                return
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -u "$USER_USERNAME:$USER_PASSWORD" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        fi
    fi
    
    # Print formatted output
    print_response "$response" "$status" "$endpoint" "$id_path"
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
            return
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
                return
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -H "Content-Type: application/json" -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" -d "$data" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        else
            response=$(curl -s -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" -X "$method" "$API_URL$endpoint")
            status=$?
            
            # Check if curl failed
            if [ $status -ne 0 ]; then
                echo -e "${RED}Curl request failed with status $status${NC}"
                return
            fi
            
            # Get HTTP status code with a separate request
            status=$(curl -s -o /dev/null -u "$ADMIN_USERNAME:$ADMIN_PASSWORD" -w "%{http_code}" -X "$method" "$API_URL$endpoint")
        fi
    fi
    
    # Print formatted output
    print_response "$response" "$status" "$endpoint" "$id_path"
}

# Start testing
print_header "CAMPUS PARKING API COMPREHENSIVE TESTS"

# --------------------------
# PUBLIC ROUTES
# --------------------------
print_header "PUBLIC ROUTES (NO AUTHENTICATION REQUIRED)"

# Test Authentication
print_header "AUTHENTICATION ENDPOINTS"

# Test User Login
USER_LOGIN_DATA='{
    "username": "'$USER_USERNAME'",
    "password": "'$USER_PASSWORD'"
}'
run_public_test "POST" "/login" "Login with Regular User" "$USER_LOGIN_DATA"

# Test Admin Login
ADMIN_LOGIN_DATA='{
    "username": "'$ADMIN_USERNAME'",
    "password": "'$ADMIN_PASSWORD'"
}'
run_public_test "POST" "/login" "Login with Admin User" "$ADMIN_LOGIN_DATA"

# Test Parking Lot Routes
print_header "PARKING LOT ENDPOINTS"
run_public_test "GET" "/parking-lots" "Get All Parking Lots"

if [ -n "$PARKING_LOT_ID" ]; then
    run_public_test "GET" "/parking-lots/$PARKING_LOT_ID" "Get Parking Lot by ID"
else
    echo -e "\n${YELLOW}>>> Skipping Get Parking Lot by ID test (no parking lots found)${NC}"
fi

# Test Building Routes
print_header "BUILDING ENDPOINTS"
run_public_test "GET" "/buildings" "Get All Buildings"

if [ -n "$BUILDING_ID" ]; then
    run_public_test "GET" "/buildings/$BUILDING_ID" "Get Building by ID"
else
    echo -e "\n${YELLOW}>>> Skipping Get Building by ID test (no buildings found)${NC}"
fi

# Test Node Routes
print_header "NODE ENDPOINTS"
run_public_test "GET" "/nodes" "Get All Nodes"

if [ -n "$NODE_ID" ]; then
    run_public_test "GET" "/nodes/$NODE_ID" "Get Node by ID"
else
    echo -e "\n${YELLOW}>>> Skipping Get Node by ID test (no nodes found)${NC}"
fi

run_public_test "GET" "/navigate?from=000000000000000000000000&to=000000000000000000000000" "Get Navigation Path (with dummy IDs)"

# Test User Routes - Create user (registration)
print_header "USER REGISTRATION ENDPOINT"
NEW_USER_DATA='{
    "name": {
        "first": "New",
        "last": "User"
    },
    "role": "user",
    "sbuId": "987654321",
    "address": {
        "street": "456 Campus Drive",
        "city": "Stony Brook",
        "state": "NY",
        "zipCode": "11794"
    },
    "driverLicense": {
        "number": "DL98765432",
        "state": "NY",
        "expirationDate": "2028-03-25T00:00:00Z"
    },
    "username": "newuser_'$(date +%s)'",
    "passwordHash": "password123"
}'

run_public_test "POST" "/users" "Register New User" "$NEW_USER_DATA" ".id"

# --------------------------
# PROTECTED ROUTES (with user auth)
# --------------------------
print_header "PROTECTED ROUTES (WITH USER AUTHENTICATION)"

# Test User Protected Routes
print_header "USER PROTECTED ENDPOINTS"
run_user_test "GET" "/users/$USER_ID" "Get Own User Profile"

UPDATED_USER_DATA='{
    "name": {
        "first": "Updated",
        "last": "User"
    },
    "role": "user",
    "sbuId": "123456789",
    "address": {
        "street": "123 Updated Drive",
        "city": "Stony Brook",
        "state": "NY",
        "zipCode": "11794"
    },
    "driverLicense": {
        "number": "DL12345678",
        "state": "NY",
        "expirationDate": "2028-03-25T00:00:00Z"
    },
    "username": "'$USER_USERNAME'",
    "passwordHash": "'$USER_PASSWORD'"
}'

run_user_test "PUT" "/users/$USER_ID" "Update Own Profile" "$UPDATED_USER_DATA"

# Test if regular user can access admin's profile (should fail)
run_user_test "GET" "/users/$ADMIN_ID" "Access Admin Profile (should fail)"

# Test Vehicle Routes
print_header "VEHICLE ENDPOINTS (USER)"
VEHICLE_DATA='{
    "name": "User Test Vehicle",
    "model": "Toyota Camry",
    "year": 2023,
    "plateNumber": "USER123"
}'

run_user_test "POST" "/vehicles" "Create Vehicle" "$VEHICLE_DATA" ".id"

# Save separate vehicle ID for user vehicle
USER_VEHICLE_ID=$VEHICLE_ID

if [ -n "$USER_VEHICLE_ID" ]; then
    run_user_test "GET" "/vehicles/$USER_VEHICLE_ID" "Get Vehicle by ID"
    
    UPDATED_VEHICLE_DATA='{
        "name": "Updated User Vehicle",
        "model": "Toyota Camry",
        "year": 2023,
        "plateNumber": "USER456"
    }'
    
    run_user_test "PUT" "/vehicles/$USER_VEHICLE_ID" "Update Vehicle" "$UPDATED_VEHICLE_DATA"
else
    echo -e "\n${YELLOW}>>> Skipping individual vehicle tests (no vehicle created)${NC}"
fi

run_user_test "GET" "/users/$USER_ID/vehicles" "Get User Vehicles"

# Test Reservation Routes (if we have a parking lot ID)
if [ -n "$PARKING_LOT_ID" ]; then
    print_header "RESERVATION ENDPOINTS (USER)"
    
    # Current time plus 1 hour
    START_TIME=$(date -u -d "+1 hour" "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+1H "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u "+%Y-%m-%dT%H:%M:%SZ") 
    # Current time plus 2 hours
    END_TIME=$(date -u -d "+2 hour" "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+2H "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u "+%Y-%m-%dT%H:%M:%SZ")
    
    USER_RESERVATION_DATA='{
        "reserved_by": "'$USER_ID'",
        "parking_lot": "'$PARKING_LOT_ID'",
        "start_time": "'$START_TIME'",
        "end_time": "'$END_TIME'",
        "status": "confirmed"
    }'
    
    run_user_test "POST" "/reservations" "Create Reservation" "$USER_RESERVATION_DATA" ".id"
    
    # Save user reservation ID
    USER_RESERVATION_ID=$RESERVATION_ID
    
    run_user_test "GET" "/users/$USER_ID/reservations" "Get User Reservations"
    
    if [ -n "$USER_RESERVATION_ID" ]; then
        run_user_test "GET" "/reservations/$USER_RESERVATION_ID" "Get Reservation by ID"
        
        # Update reservation: extend end time by 1 hour
        UPDATED_END_TIME=$(date -u -d "+3 hour" "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+3H "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u "+%Y-%m-%dT%H:%M:%SZ")
        
        UPDATED_RESERVATION_DATA='{
            "reserved_by": "'$USER_ID'",
            "parking_lot": "'$PARKING_LOT_ID'",
            "start_time": "'$START_TIME'",
            "end_time": "'$UPDATED_END_TIME'",
            "status": "confirmed"
        }'
        
        run_user_test "PUT" "/reservations/$USER_RESERVATION_ID" "Update Reservation" "$UPDATED_RESERVATION_DATA"
    else
        echo -e "\n${YELLOW}>>> Skipping individual reservation tests (no reservation created)${NC}"
    fi
else
    echo -e "\n${YELLOW}>>> Skipping reservation tests (no parking lot found)${NC}"
fi

# --------------------------
# ADMIN ROUTES (with admin auth)
# --------------------------
print_header "ADMIN ROUTES (WITH ADMIN AUTHENTICATION)"

# Test Admin User Routes
print_header "ADMIN USER ENDPOINTS"
run_admin_test "GET" "/admin/users" "Get All Users"

# Test Admin Parking Lot Routes
print_header "ADMIN PARKING LOT ENDPOINTS"
NEW_PARKING_LOT='{
    "name": "Test Admin Lot",
    "spaces": 100,
    "faculty": 50,
    "premium": 10,
    "metered": 10,
    "resident": 25,
    "ada": 5,
    "ev": true,
    "active": true,
    "location": {
        "lat": 40.91,
        "lng": -73.12
    }
}'

run_admin_test "POST" "/admin/parking-lots" "Create Parking Lot" "$NEW_PARKING_LOT" ".id"

if [ -n "$PARKING_LOT_ID" ]; then
    UPDATED_PARKING_LOT='{
        "name": "Updated Admin Lot",
        "spaces": 120,
        "faculty": 60,
        "premium": 15,
        "metered": 10,
        "resident": 30,
        "ada": 5,
        "ev": true,
        "active": true,
        "location": {
            "lat": 40.91,
            "lng": -73.12
        }
    }'
    
    run_admin_test "PUT" "/admin/parking-lots/$PARKING_LOT_ID" "Update Parking Lot" "$UPDATED_PARKING_LOT"
    
    if [ -n "$PARKING_LOT_ID" ]; then
        run_admin_test "GET" "/admin/parking-lots/$PARKING_LOT_ID/reservations" "Get Parking Lot Reservations"
    fi
fi

# Test Admin Building Routes
print_header "ADMIN BUILDING ENDPOINTS"
NEW_BUILDING='{
    "name": "Test Admin Building",
    "location": {
        "lat": 40.91,
        "lng": -73.12
    },
    "node": "000000000000000000000000"
}'

run_admin_test "POST" "/admin/buildings" "Create Building" "$NEW_BUILDING" ".id"

if [ -n "$BUILDING_ID" ]; then
    UPDATED_BUILDING='{
        "name": "Updated Admin Building",
        "location": {
            "lat": 40.91,
            "lng": -73.12
        },
        "node": "000000000000000000000000"
    }'
    
    run_admin_test "PUT" "/admin/buildings/$BUILDING_ID" "Update Building" "$UPDATED_BUILDING"
fi

# Test Admin Node Routes
print_header "ADMIN NODE ENDPOINTS"
NEW_NODE='{
    "type": "Intersection",
    "location": {
        "lat": 40.91,
        "lng": -73.12
    },
    "neighbors": {
        "members": []
    }
}'

run_admin_test "POST" "/admin/nodes" "Create Node" "$NEW_NODE" ".id"

if [ -n "$NODE_ID" ]; then
    UPDATED_NODE='{
        "type": "Building",
        "location": {
            "lat": 40.92, 
            "lng": -73.13
        },
        "neighbors": {
            "members": []
        }
    }'
    
    run_admin_test "PUT" "/admin/nodes/$NODE_ID" "Update Node" "$UPDATED_NODE"
fi

# Test Admin Reservation Routes
print_header "ADMIN RESERVATION ENDPOINTS"
run_admin_test "GET" "/admin/reservations" "Get All Reservations"

# Test Admin Violation Routes
print_header "ADMIN VIOLATION ENDPOINTS"
run_admin_test "GET" "/admin/violations" "Get All Violations"

if [ -n "$CREATED_USER_ID" ] && [ -n "$PARKING_LOT_ID" ]; then
    # Create a violation as admin
    PAY_BY_TIME=$(date -u -d "+7 days" "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+7d "+%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u "+%Y-%m-%dT%H:%M:%SZ")
    
    VIOLATION_DATA='{
        "user": "'$CREATED_USER_ID'",
        "parking_lot": "'$PARKING_LOT_ID'",
        "reason": "Parked in reserved space",
        "fine": 50.00,
        "pay_by": "'$PAY_BY_TIME'"
    }'
    
    run_admin_test "POST" "/admin/violations" "Create Violation" "$VIOLATION_DATA" ".id"
    
    if [ -n "$VIOLATION_ID" ]; then
        UPDATED_VIOLATION_DATA='{
            "user": "'$CREATED_USER_ID'",
            "parking_lot": "'$PARKING_LOT_ID'",
            "reason": "Parked in faculty space",
            "fine": 75.00,
            "pay_by": "'$PAY_BY_TIME'"
        }'
        
        run_admin_test "PUT" "/admin/violations/$VIOLATION_ID" "Update Violation" "$UPDATED_VIOLATION_DATA"
    fi
fi

# Test Admin Violation Rebute Routes
print_header "ADMIN VIOLATION REBUTE ENDPOINTS"
run_admin_test "GET" "/admin/violation-rebutes" "Get All Violation Rebutes"

if [ -n "$REBUTE_ID" ]; then
    UPDATED_REBUTE_DATA='{
        "reason": "Updated test rebute reason"
    }'
    
    run_admin_test "PUT" "/admin/violation-rebutes/$REBUTE_ID" "Update Violation Rebute" "$UPDATED_REBUTE_DATA"
fi

# --------------------------
# Test permission boundaries
# --------------------------
print_header "TESTING PERMISSION BOUNDARIES"

# Test if regular user can access admin-only routes
print_header "TESTING REGULAR USER ACCESSING ADMIN ROUTES (SHOULD FAIL)"
run_user_test "GET" "/admin/users" "Access Admin-only User List"
run_user_test "POST" "/admin/parking-lots" "Create Parking Lot via Admin Route" "$NEW_PARKING_LOT"

# --------------------------
# CLEANUP
# --------------------------
print_header "CLEANING UP"

# Clean up - delete resources in reverse order of creation

# Delete objects created during testing
if [ -n "$REBUTE_ID" ]; then
    run_admin_test "DELETE" "/admin/violation-rebutes/$REBUTE_ID" "Delete Violation Rebute"
fi

if [ -n "$VIOLATION_ID" ]; then
    run_admin_test "DELETE" "/admin/violations/$VIOLATION_ID" "Delete Violation"
fi

if [ -n "$USER_RESERVATION_ID" ]; then
    run_user_test "DELETE" "/reservations/$USER_RESERVATION_ID" "Cancel User Reservation"
fi

if [ -n "$USER_VEHICLE_ID" ]; then
    run_user_test "DELETE" "/vehicles/$USER_VEHICLE_ID" "Delete User Vehicle"
fi

if [ -n "$NODE_ID" ]; then
    run_admin_test "DELETE" "/admin/nodes/$NODE_ID" "Delete Node"
fi

if [ -n "$BUILDING_ID" ]; then
    run_admin_test "DELETE" "/admin/buildings/$BUILDING_ID" "Delete Building"
fi

if [ -n "$PARKING_LOT_ID" ]; then
    run_admin_test "DELETE" "/admin/parking-lots/$PARKING_LOT_ID" "Delete Parking Lot"
fi

if [ -n "$CREATED_USER_ID" ]; then
    run_admin_test "DELETE" "/admin/users/$CREATED_USER_ID" "Delete Created Test User"
fi

echo -e "\n${GREEN}Testing Complete!${NC}"