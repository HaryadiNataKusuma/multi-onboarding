#!/bin/bash

echo "🧪 Testing /api/insurances endpoint..."
echo ""

# Test endpoint
echo "📡 Making request to http://localhost:8080/api/insurances..."
echo ""

response=$(curl -s -w "\nHTTP_CODE:%{http_code}" http://localhost:8080/api/insurances)
http_code=$(echo "$response" | grep "HTTP_CODE" | cut -d: -f2)
body=$(echo "$response" | sed '/HTTP_CODE/d')

echo "📊 Response Status: $http_code"
echo ""
echo "📦 Response Body:"
echo "$body" | jq '.' 2>/dev/null || echo "$body"
echo ""

if [ "$http_code" = "200" ]; then
    count=$(echo "$body" | jq 'length' 2>/dev/null || echo "0")
    echo "✅ Success! Found $count insurance(s)"
    
    if [ "$count" = "0" ]; then
        echo "⚠️  Warning: No insurance data found in database"
        echo "💡 Check if data exists in travel_service_development.insurances table"
    else
        echo "📋 First insurance:"
        echo "$body" | jq '.[0]' 2>/dev/null || echo "Could not parse JSON"
    fi
else
    echo "❌ Error: HTTP $http_code"
    echo "💡 Check server logs for details"
fi


