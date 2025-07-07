#!/bin/bash

# Test build script for CodeQuest CLI
cd "$(dirname "$0")"

echo "Testing Go module and build..."

# Clean and update dependencies
go mod tidy

# Build the CLI
echo "Building CLI..."
go build -o codequest main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "Testing CLI..."
    ./codequest --help
    echo ""
    echo "✅ CLI is working correctly!"
else
    echo "❌ Build failed!"
    exit 1
fi