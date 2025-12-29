#!/bin/bash

# Script to stop the server

echo "🛑 Stopping server..."

# Kill processes on port 8080
if lsof -ti:8080 > /dev/null 2>&1; then
    PIDS=$(lsof -ti:8080)
    echo "   Killing processes on port 8080: $PIDS"
    echo $PIDS | xargs kill -9 2>/dev/null
fi

# Kill go run main.go processes
GO_PIDS=$(ps aux | grep "[g]o run main.go" | grep -v grep | awk '{print $2}')
if [ ! -z "$GO_PIDS" ]; then
    echo "   Killing Go server processes: $GO_PIDS"
    echo $GO_PIDS | xargs kill -9 2>/dev/null
fi

sleep 1

# Verify
if lsof -ti:8080 > /dev/null 2>&1; then
    echo "⚠️  Some processes still using port 8080"
else
    echo "✅ Port 8080 is now free"
fi

echo "✅ Server stopped"


