#!/bin/bash

# Setup script for CodeQuest CLI repository
cd "$(dirname "$0")"

echo "Setting up git repository..."

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: CodeQuest CLI with native executor

- Native Go and Node.js execution (no Docker required)
- Multi-language support: Go, JavaScript, TypeScript
- Cross-platform binary releases
- Standalone CLI for offline challenge practice

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# Add remote repository
git remote add origin https://github.com/crisecheverria/codequest.git

# Set main branch
git branch -M main

echo "Repository setup complete!"
echo ""
echo "To push to GitHub, run:"
echo "  git push -u origin main"
echo ""
echo "To test the CLI, run:"
echo "  go build -o codequest main.go"
echo "  ./codequest --help"