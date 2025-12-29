#!/bin/bash

# Script to check server status

echo "🔍 Checking server status..."
echo ""

# Check if port 8080 is in use
if lsof -ti:8080 > /dev/null 2>&1; then
    PID=$(lsof -ti:8080)
    echo "✅ Port 8080 is in use by PID: $PID"
    
    # Check if it's our Go server
    if ps -p $PID > /dev/null 2>&1; then
        CMD=$(ps -p $PID -o comm=)
        echo "   Process: $CMD"
    fi
else
    echo "❌ Port 8080 is not in use"
fi

echo ""

# Check for go run main.go processes
GO_PROCESSES=$(ps aux | grep "[g]o run main.go" | wc -l)
if [ "$GO_PROCESSES" -gt 0 ]; then
    echo "✅ Found $GO_PROCESSES Go server process(es):"
    ps aux | grep "[g]o run main.go" | grep -v grep
else
    echo "❌ No 'go run main.go' processes found"
fi

echo ""

# Test server response
echo "🌐 Testing server response..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/regions 2>/dev/null || echo "000")

case $HTTP_CODE in
    200)
        echo "✅ Server is responding (HTTP 200 OK)"
        ;;
    404)
        echo "⚠️  Server is responding but endpoint not found (HTTP 404)"
        ;;
    500)
        echo "⚠️  Server is responding but has errors (HTTP 500)"
        ;;
    000)
        echo "❌ Server is not responding (Connection refused)"
        ;;
    *)
        echo "⚠️  Server returned HTTP $HTTP_CODE"
        ;;
esac

echo ""

# Check for log file
if [ -f "server.log" ]; then
    echo "📋 Last 10 lines of server.log:"
    tail -10 server.log
else
    echo "📋 No server.log file found"
fi


