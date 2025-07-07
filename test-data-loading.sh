#!/bin/bash

# Test data loading after fixing the paths
cd "$(dirname "$0")"

echo "🧪 Testing data loading fix..."

# Build the CLI
echo "📦 Building CLI..."
go build -o codequest main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

# Test list command
echo "📋 Testing list command..."
./codequest list | head -10

if [ $? -eq 0 ]; then
    echo "✅ List command works!"
else
    echo "❌ List command failed"
    exit 1
fi

# Test fetching a challenge
echo "🎯 Testing fetch command..."
./codequest fetch go-multiple-return-values

if [ -d "challenge-go-multiple-return-values" ]; then
    echo "✅ Fetch command works!"
    rm -rf challenge-go-multiple-return-values
else
    echo "❌ Fetch command failed"
    exit 1
fi

echo ""
echo "🎉 All tests passed! Ready to push the fix."
echo ""
echo "Next steps:"
echo "1. git add ."
echo "2. git commit -m 'Fix: Include challenges data and update loader paths'"
echo "3. git push origin main"
echo "4. git tag v1.0.2"
echo "5. git push origin v1.0.2"