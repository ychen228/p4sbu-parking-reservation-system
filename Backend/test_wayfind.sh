#!/bin/bash

# Configuration 
API_URL="http://localhost:8080/api/v1"
USER_CREDENTIALS_FILE=".user_credentials.json"
ADMIN_CREDENTIALS_FILE=".admin_credentials.json"

# Colors for better readability
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Check if credential files exist (for authenticated requests if needed)
has_credentials=true
if [ ! -f "$USER_CREDENTIALS_FILE" ] || [ ! -f "$ADMIN_CREDENTIALS_FILE" ]; then
    has_credentials=false
    echo -e "${YELLOW}Credentials files not found. Some tests may fail if authentication is required.${NC}"
else
    # Load credentials
    USER_USERNAME=$(cat "$USER_CREDENTIALS_FILE" | jq -r '.username')
    USER_PASSWORD=$(cat "$USER_CREDENTIALS_FILE" | jq -r '.password')
    
    ADMIN_USERNAME=$(cat "$ADMIN_CREDENTIALS_FILE" | jq -r '.username')
    ADMIN_PASSWORD=$(cat "$ADMIN_CREDENTIALS_FILE" | jq -r '.password')
    
    echo -e "${GREEN}Loaded credentials for user: $USER_USERNAME and admin: $ADMIN_USERNAME${NC}"
fi

echo -e "${BLUE}${BOLD}==== P4SBU Nearest Parking Lots ====${NC}"
echo -e "${YELLOW}Using API URL: ${BASE_URL}${NC}\n"

# Building ID for Engineering Building
BUILDING_ID="67f5beac8bd61c7932d00cc1"
BUILDING_NAME="Engineering"

echo -e "${GREEN}[Query]${NC} Finding nearest parking lots to ${BLUE}${BUILDING_NAME}${NC} building..."

# Make the API call - if credentials are required, use them
if [ "$has_credentials" = true ]; then
    RESPONSE=$(curl -s -u "$USER_USERNAME:$USER_PASSWORD" "${API_URL}/wayfind/building/${BUILDING_ID}/nearest-parking")
else
    RESPONSE=$(curl -s "${API_URL}/wayfind/building/${BUILDING_ID}/nearest-parking")
fi

# Try to parse the response and check if it's valid JSON with parking lots
if [[ "$RESPONSE" == "["* ]]; then
    # It looks like a JSON array
    echo -e "\n${GREEN}[Results]${NC} Nearest parking lots to ${BLUE}${BUILDING_NAME}${NC}:"
    echo -e "${YELLOW}─────────────────────────────────────────────────────────────────────────────${NC}"
    echo -e "${BOLD}RANK   PARKING LOT                                     DISTANCE    SPACES${NC}"
    echo -e "${YELLOW}─────────────────────────────────────────────────────────────────────────────${NC}"
    
    # Create a temporary file to store the response
    TEMP_FILE=$(mktemp)
    echo "$RESPONSE" > "$TEMP_FILE"
    
    # Extract the top 5 parking lots - improved extraction to maintain full names
    COUNT=0
    LIMIT=5
    
    # Process each entry in the JSON array
    jq_available=0
    if command -v jq &> /dev/null; then
        jq_available=1
        # Use jq to extract data if available (much more reliable)
        TOTAL_COUNT=$(jq '. | length' "$TEMP_FILE")
        if [ "$TOTAL_COUNT" -lt "$LIMIT" ]; then
            LIMIT=$TOTAL_COUNT
        fi
        
        for ((i=0; i<LIMIT; i++)); do
            NAME=$(jq -r ".[$i].parking_lot.Name // .[$i].parking_lot.name" "$TEMP_FILE")
            DISTANCE=$(jq -r ".[$i].distance_km" "$TEMP_FILE")
            SPACES=$(jq -r ".[$i].parking_lot.Spaces // .[$i].parking_lot.spaces" "$TEMP_FILE")
            
            # Format the distance to be readable
            FORMATTED_DISTANCE=$(printf "%.6f" $DISTANCE)
            
            # Display with nice formatting
            echo -e " ${CYAN}${BOLD}$(($i+1))${NC}     $(printf "%-45s" "$NAME") ${FORMATTED_DISTANCE}  ${SPACES}"
            COUNT=$((COUNT+1))
        done
        
        TOTAL_COUNT=$(jq '. | length' "$TEMP_FILE")
    else
        # Fallback to manual extraction if jq is not available
        # This approach is less reliable but should handle spaces better
        while IFS= read -r line; do
            if [[ $line == *'"Name":'* || $line == *'"name":'* ]]; then
                NAME=$(echo "$line" | sed -E 's/.*"[Nn]ame":"([^"]+)".*/\1/')
                # Read the next lines to get distance and spaces
                read -r distance_line
                DISTANCE=$(echo "$distance_line" | grep -o '"distance_km":[^,}]*' | cut -d':' -f2)
                
                # Look for Spaces within next few lines
                SPACES=""
                for j in {1..10}; do
                    read -r spaces_line
                    if [[ $spaces_line == *'"Spaces":'* || $spaces_line == *'"spaces":'* ]]; then
                        SPACES=$(echo "$spaces_line" | grep -o '"[Ss]paces":[^,}]*' | cut -d':' -f2)
                        break
                    fi
                done
                
                # Format the distance
                FORMATTED_DISTANCE=$(printf "%.6f" $DISTANCE)
                
                # Display with nice formatting
                echo -e " ${CYAN}${BOLD}$((COUNT+1))${NC}     $(printf "%-45s" "$NAME") ${FORMATTED_DISTANCE}  ${SPACES}"
                
                COUNT=$((COUNT+1))
                if [ $COUNT -ge $LIMIT ]; then
                    break
                fi
            fi
        done < "$TEMP_FILE"
        
        # Count total entries
        TOTAL_COUNT=$(grep -o '"parking_lot"' "$TEMP_FILE" | wc -l)
    fi
    
    echo -e "${YELLOW}─────────────────────────────────────────────────────────────────────────────${NC}"
    echo -e "${GREEN}Found ${TOTAL_COUNT} parking lots in total${NC}"
    
    # Clean up temp file
    rm -f "$TEMP_FILE"
else
    echo -e "${RED}[Error]${NC} Invalid or empty response from the API"
    echo -e "Make sure your server is running and the wayfinding handler is correctly implemented"
    echo -e "Response: $RESPONSE"
fi

# Try with direct coordinates for a central campus location
echo -e "\n${GREEN}[Query]${NC} Finding nearest parking lots to campus center coordinates..."
LAT="40.9140"
LNG="-73.1240"

echo -e "Using coordinates: Lat=${BLUE}${LAT}${NC}, Lng=${BLUE}${LNG}${NC}"

# Make the API call with coordinates - if credentials are required, use them
if [ "$has_credentials" = true ]; then
    COORD_RESPONSE=$(curl -s -u "$USER_USERNAME:$USER_PASSWORD" "${API_URL}/wayfind/nearest-parking?lat=${LAT}&lng=${LNG}")
else
    COORD_RESPONSE=$(curl -s "${API_URL}/wayfind/nearest-parking?lat=${LAT}&lng=${LNG}")
fi

# Try to parse the response and check if it's valid JSON with parking lots
if [[ "$COORD_RESPONSE" == "["* ]]; then
    # It looks like a JSON array
    echo -e "\n${GREEN}[Results]${NC} Nearest parking lots to coordinates (${LAT}, ${LNG}):"
    echo -e "${YELLOW}─────────────────────────────────────────────────────────────────────────────${NC}"
    echo -e "${BOLD}RANK   PARKING LOT                                     DISTANCE    SPACES${NC}"
    echo -e "${YELLOW}─────────────────────────────────────────────────────────────────────────────${NC}"
    
    # Create a temporary file to store the response
    TEMP_FILE=$(mktemp)
    echo "$COORD_RESPONSE" > "$TEMP_FILE"
    
    # Extract the top 5 parking lots
    COUNT=0
    LIMIT=5
    
    # Process each entry in the JSON array
    jq_available=0
    if command -v jq &> /dev/null; then
        jq_available=1
        # Use jq to extract data if available (much more reliable)
        TOTAL_COUNT=$(jq '. | length' "$TEMP_FILE")
        if [ "$TOTAL_COUNT" -lt "$LIMIT" ]; then
            LIMIT=$TOTAL_COUNT
        fi
        
        for ((i=0; i<LIMIT; i++)); do
            NAME=$(jq -r ".[$i].parking_lot.Name // .[$i].parking_lot.name" "$TEMP_FILE")
            DISTANCE=$(jq -r ".[$i].distance_km" "$TEMP_FILE")
            SPACES=$(jq -r ".[$i].parking_lot.Spaces // .[$i].parking_lot.spaces" "$TEMP_FILE")
            
            # Format the distance to be readable
            FORMATTED_DISTANCE=$(printf "%.6f" $DISTANCE)
            
            # Display with nice formatting
            echo -e " ${CYAN}${BOLD}$(($i+1))${NC}     $(printf "%-45s" "$NAME") ${FORMATTED_DISTANCE}  ${SPACES}"
            COUNT=$((COUNT+1))
        done
        
        TOTAL_COUNT=$(jq '. | length' "$TEMP_FILE")
    else
        # Fallback to manual extraction if jq is not available
        python_available=0
        if command -v python3 &> /dev/null; then
            python_available=1
            # Use Python to parse JSON if available
            RESULT=$(python3 -c "
import json, sys
data = json.load(sys.stdin)
count = min(5, len(data))
total = len(data)
for i in range(count):
    name = data[i]['parking_lot'].get('Name', data[i]['parking_lot'].get('name', 'Unknown'))
    dist = data[i]['distance_km']
    spaces = data[i]['parking_lot'].get('Spaces', data[i]['parking_lot'].get('spaces', 0))
    print(f'{i+1}|{name}|{dist:.6f}|{spaces}')
print(f'total|{total}')
            " < "$TEMP_FILE")
            
            # Process the Python output
            while IFS= read -r line; do
                if [[ $line == total* ]]; then
                    TOTAL_COUNT=${line#*|}
                else
                    IFS='|' read -r rank name dist spaces <<< "$line"
                    echo -e " ${CYAN}${BOLD}${rank}${NC}     $(printf "%-45s" "$name") ${dist}  ${spaces}"
                    COUNT=$((COUNT+1))
                fi
            done <<< "$RESULT"
        else
            # Last resort: manual parsing with grep/sed
            while IFS= read -r line; do
                if [[ $line == *'"Name":'* || $line == *'"name":'* ]]; then
                    NAME=$(echo "$line" | sed -E 's/.*"[Nn]ame":"([^"]+)".*/\1/')
                    # Read the next lines to get distance and spaces
                    read -r distance_line
                    DISTANCE=$(echo "$distance_line" | grep -o '"distance_km":[^,}]*' | cut -d':' -f2)
                    
                    # Look for Spaces within next few lines
                    SPACES=""
                    for j in {1..10}; do
                        read -r spaces_line
                        if [[ $spaces_line == *'"Spaces":'* || $spaces_line == *'"spaces":'* ]]; then
                            SPACES=$(echo "$spaces_line" | grep -o '"[Ss]paces":[^,}]*' | cut -d':' -f2)
                            break
                        fi
                    done
                    
                    # Format the distance
                    FORMATTED_DISTANCE=$(printf "%.6f" $DISTANCE)
                    
                    # Display with nice formatting
                    echo -e " ${CYAN}${BOLD}$((COUNT+1))${NC}     $(printf "%-45s" "$NAME") ${FORMATTED_DISTANCE}  ${SPACES}"
                    
                    COUNT=$((COUNT+1))
                    if [ $COUNT -ge $LIMIT ]; then
                        break
                    fi
                fi
            done < "$TEMP_FILE"
            
            # Count total entries
            TOTAL_COUNT=$(grep -o '"parking_lot"' "$TEMP_FILE" | wc -l)
        fi
    fi
    
    echo -e "${YELLOW}─────────────────────────────────────────────────────────────────────────────${NC}"
    echo -e "${GREEN}Found ${TOTAL_COUNT} parking lots in total${NC}"
    
    # Clean up temp file
    rm -f "$TEMP_FILE"
else
    echo -e "${RED}[Error]${NC} Invalid or empty response from the API"
    echo -e "Response: $COORD_RESPONSE"
fi

echo -e "\n${BLUE}${BOLD}==== Display Complete ====${NC}"