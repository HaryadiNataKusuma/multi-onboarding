#!/bin/bash

# Quick restart script - stops and starts server immediately

echo "🔄 Quick Restart Server..."

# Kill processes on port 8080
lsof -ti:8080 | xargs kill -9 2>/dev/null
pkill -f "go run main.go" 2>/dev/null
sleep 1

# Start server in background
cd "$(dirname "$0")"
nohup go run main.go > server.log 2>&1 &
SERVER_PID=$!

echo "✅ Server restarted! PID: $SERVER_PID"
echo "📋 Check logs: tail -f server.log"
echo "🌐 Server: http://localhost:8080"

sleep 2
echo ""
echo "📊 Last 10 lines of server.log:"
tail -10 server.log 2>/dev/null || echo "Log file not created yet"


