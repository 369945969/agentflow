#!/bin/bash
# Validate if the output directory follows the four-folder rule

PROJECT_DIR=$1

if [ -z "$PROJECT_DIR" ]; then
    echo "Usage: ./validate.sh [path_to_output/project_name]"
    exit 1
fi

REQUIRED_FOLDERS=("Global&Context" "Feature&Plan" "Style&Guide" "Screen&Prototype")
MISSING=0

for folder in "${REQUIRED_FOLDERS[@]}"; do
    if [ ! -d "$PROJECT_DIR/$folder" ]; then
        echo "[ERROR] Missing folder: $folder"
        MISSING=$((MISSING + 1))
    else
        echo "[OK] Found folder: $folder"
    fi
done

if [ $MISSING -eq 0 ]; then
    echo "Validation Successful: All 4 standard folders exist."
    echo "Starting simple HTTP server at http://localhost:8000"
    cd "$PROJECT_DIR"
    python3 -m http.server 8000 &
    SERVER_PID=$!
    
    # Wait a moment for the server to start
    sleep 2
    
    echo "Opening browser..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open "http://localhost:8000"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v xdg-open > /dev/null; then
            xdg-open "http://localhost:8000"
        else
            echo "Please open http://localhost:8000 manually"
        fi
    elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        start "http://localhost:8000"
    else
        echo "Please open http://localhost:8000 manually"
    fi

    echo "Server is running (PID: $SERVER_PID). Press Ctrl+C to stop."
    wait $SERVER_PID

else
    echo "Validation Failed: $MISSING folders missing."
    exit 1
fi
