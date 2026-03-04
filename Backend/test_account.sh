#!/bin/bash

# Create Account Script for Campus Parking API
# Creates both a regular user and an admin account for testing

# Configuration
API_URL="http://localhost:8080/api/v1"
USER_CREDENTIALS_FILE=".user_credentials.json"
ADMIN_CREDENTIALS_FILE=".admin_credentials.json"

# Fixed credentials as requested
USER_USERNAME="testuser1"
USER_PASSWORD="password123"
ADMIN_USERNAME="testadmin1"
ADMIN_PASSWORD="password123"

# Colors for terminal output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}jq is not installed. Install it for better JSON formatting:${NC}"
    echo -e "${YELLOW}  sudo apt-get install jq (Debian/Ubuntu)${NC}"
    echo -e "${YELLOW}  brew install jq (macOS)${NC}"
    exit 1
fi

# Function to print section headers
print_header() {
    echo -e "\n${BLUE}======================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}======================================================${NC}"
}

# Function to create an account and save credentials
create_account() {
    local username="$1"
    local password="$2"
    local role="$3"
    local credentials_file="$4"
    
    print_header "CREATING $role ACCOUNT"
    echo -e "${CYAN}Username: $username${NC}"
    echo -e "${CYAN}Password: $password${NC}"
    
    # User registration data
    local user_data='{
        "name": {
            "first": "Test",
            "last": "'$role'"
        },
        "role": "'$role'",
        "sbuId": "123456789",
        "address": {
            "street": "123 Campus Drive",
            "city": "Stony Brook",
            "state": "NY",
            "zipCode": "11794"
        },
        "driverLicense": {
            "number": "DL12345678",
            "state": "NY",
            "expirationDate": "2028-03-25T00:00:00Z"
        },
        "username": "'$username'",
        "passwordHash": "'$password'"
    }'
    
    # Register the user
    echo -e "${CYAN}Registering $role...${NC}"
    local register_response=$(curl -s -H "Content-Type: application/json" -d "$user_data" -X POST "$API_URL/users")
    local register_status=$(curl -s -o /dev/null -H "Content-Type: application/json" -d "$user_data" -w "%{http_code}" -X POST "$API_URL/users")
    
    echo -e "${CYAN}Register status: $register_status${NC}"
    
    # Login data
    local login_data='{
        "username": "'$username'",
        "password": "'$password'"
    }'
    
    local login_needed=true
    local user_id=""
    
    # If registration succeeded, extract user ID
    if [ "$register_status" -eq 201 ] || [ "$register_status" -eq 200 ]; then
        echo -e "${GREEN}$role registered successfully!${NC}"
        user_id=$(echo "$register_response" | jq -r '.id')
    elif [ "$register_status" -eq 409 ]; then
        echo -e "${YELLOW}$role already exists. Trying to login...${NC}"
    else
        echo -e "${RED}$role registration failed with status code: $register_status${NC}"
        echo -e "${RED}Response: $register_response${NC}"
        echo -e "${YELLOW}Trying to login anyway in case the account exists...${NC}"
    fi
    
    # Login to get user ID and role
    echo -e "${CYAN}Logging in to get user information...${NC}"
    
    local login_response=$(curl -s -H "Content-Type: application/json" -d "$login_data" -X POST "$API_URL/login")
    local login_status=$(curl -s -o /dev/null -H "Content-Type: application/json" -d "$login_data" -w "%{http_code}" -X POST "$API_URL/login")
    
    echo -e "${CYAN}Login status: $login_status${NC}"
    
    if [ "$login_status" -eq 200 ]; then
        # If we didn't get a user ID from registration, extract it from login response
        if [ -z "$user_id" ]; then
            user_id=$(echo "$login_response" | jq -r '.userId')
        fi
        
        # Extract role from the response
        local user_role=$(echo "$login_response" | jq -r '.role')
        
        echo -e "${GREEN}Login successful!${NC}"
        
        # Save credentials to file
        echo '{
            "username": "'$username'",
            "password": "'$password'",
            "userId": "'$user_id'",
            "role": "'$user_role'"
        }' | jq > "$credentials_file"
        
        echo -e "${GREEN}Credentials saved to $credentials_file${NC}"
        return 0
    else
        echo -e "${RED}Login failed with status code: $login_status${NC}"
        echo -e "${RED}Response: $login_response${NC}"
        return 1
    fi
}

# Check if credentials files already exist
user_credentials_exist=false
admin_credentials_exist=false

if [ -f "$USER_CREDENTIALS_FILE" ]; then
    user_credentials_exist=true
fi

if [ -f "$ADMIN_CREDENTIALS_FILE" ]; then
    admin_credentials_exist=true
fi

if [ "$user_credentials_exist" = true ] && [ "$admin_credentials_exist" = true ]; then
    print_header "EXISTING CREDENTIALS FOUND"
    echo -e "${YELLOW}Existing credentials files were found:${NC}"
    echo -e "- User: $USER_CREDENTIALS_FILE"
    echo -e "- Admin: $ADMIN_CREDENTIALS_FILE"
    echo -e "${YELLOW}Would you like to use the existing credentials or create new accounts?${NC}"
    echo -e "1) Use existing credentials"
    echo -e "2) Create new accounts (will overwrite existing credentials)"
    
    read -p "Enter your choice (1-2): " choice
    
    if [ "$choice" == "1" ]; then
        echo -e "${GREEN}Using existing credentials.${NC}"
        echo -e "${CYAN}User credentials:${NC}"
        cat "$USER_CREDENTIALS_FILE" | jq
        echo -e "${CYAN}Admin credentials:${NC}"
        cat "$ADMIN_CREDENTIALS_FILE" | jq
        exit 0
    fi
    
    echo -e "${YELLOW}Creating new accounts (existing credentials will be overwritten)...${NC}"
fi

# Create regular user account
create_account "$USER_USERNAME" "$USER_PASSWORD" "user" "$USER_CREDENTIALS_FILE"
user_success=$?

# Create admin account
create_account "$ADMIN_USERNAME" "$ADMIN_PASSWORD" "admin" "$ADMIN_CREDENTIALS_FILE"
admin_success=$?

# Check if both accounts were created/logged in successfully
if [ $user_success -eq 0 ] && [ $admin_success -eq 0 ]; then
    echo -e "\n${GREEN}Both user and admin accounts are ready for testing!${NC}"
    echo -e "${CYAN}You can now run test_api.sh to test protected and admin endpoints.${NC}"
    exit 0
else
    echo -e "\n${RED}Failed to set up one or both accounts.${NC}"
    exit 1
fi