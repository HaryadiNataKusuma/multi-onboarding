#!/bin/bash

echo "🔄 Restarting server..."

# Stop existing server
./stop_server.sh

# Wait a bit
sleep 2

# Start server
./start_server.sh


