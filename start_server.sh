#!/bin/bash

# Script to start the Go server with proper error handling

echo "🔍 Checking port 8080..."

# Kill any process using port 8080
if lsof -ti:8080 > /dev/null 2>&1; then
    echo "⚠️  Port 8080 is in use. Killing existing process..."
    lsof -ti:8080 | xargs kill -9 2>/dev/null
    sleep 2
fi

# Kill any existing go run main.go processes
pkill -f "go run main.go" 2>/dev/null
sleep 1

echo "✅ Port 8080 is now free"
echo ""
echo "🚀 Starting server..."
echo ""

# Start the server and capture output
cd "$(dirname "$0")"
go run main.go 2>&1 | tee server.log &
SERVER_PID=$!

echo "📝 Server PID: $SERVER_PID"
echo "📋 Logs are being written to server.log"
echo ""

# Wait a bit for server to start
sleep 3

# Check if server is running
if ps -p $SERVER_PID > /dev/null 2>&1; then
    echo "✅ Server process is running (PID: $SERVER_PID)"
    
    # Test if server is responding
    sleep 2
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/regions 2>/dev/null || echo "000")
    
    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "404" ] || [ "$HTTP_CODE" = "500" ]; then
        echo "✅ Server is responding (HTTP $HTTP_CODE)"
        echo "🌐 Server is ready at http://localhost:8080"
    else
        echo "⚠️  Server started but not responding yet (HTTP $HTTP_CODE)"
        echo "📋 Check server.log for details"
    fi
else
    echo "❌ Server failed to start"
    echo "📋 Check server.log for error details:"
    tail -20 server.log 2>/dev/null || echo "No log file found"
fi

echo ""
echo "💡 To stop the server, run: kill $SERVER_PID"
echo "💡 Or use: pkill -f 'go run main.go'"
echo "💡 To view logs: tail -f server.log"


