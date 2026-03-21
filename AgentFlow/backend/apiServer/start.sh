#!/bin/bash

# Configuration
PORT=3000
BINARY_NAME="api_server"
LOG_FILE="api_server.log"

echo "🔍 Checking for processes on port $PORT..."

# Find PID of process using the port
PID=$(lsof -t -i:$PORT)

if [ -n "$PID" ]; then
    echo "⚠️ Found process $PID using port $PORT. Terminating..."
    kill -9 $PID
    sleep 1
    echo "✅ Process terminated."
else
    echo "✅ Port $PORT is free."
fi

# Build the application
echo "🔨 Building $BINARY_NAME..."
go build -o $BINARY_NAME main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed. Exiting."
    exit 1
fi

# Start the application in the background
echo "🚀 Starting $BINARY_NAME on port $PORT..."
nohup ./$BINARY_NAME > $LOG_FILE 2>&1 &

# Wait a moment to check if it started successfully
sleep 2
if ps -p $! > /dev/null; then
    echo "✅ Server is running! (PID: $!)"
    echo "📝 Logs are being written to $LOG_FILE"
    echo "🔗 URL: http://localhost:$PORT"
else
    echo "❌ Server failed to start. Check $LOG_FILE for details."
    cat $LOG_FILE
    exit 1
fi
