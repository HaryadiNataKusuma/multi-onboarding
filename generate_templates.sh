#!/bin/bash

# Script to generate template drafts for all existing products

echo "🔄 Generating template drafts for all existing products..."

# Call the API endpoint
response=$(curl -s -X POST http://localhost:8080/api/products/generate-templates \
  -H "Content-Type: application/json" \
  -d '{"created_by": 1}')

# Check if curl was successful
if [ $? -eq 0 ]; then
  echo "✅ Response:"
  echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
else
  echo "❌ Failed to connect to server. Make sure the server is running on http://localhost:8080"
  exit 1
fi


