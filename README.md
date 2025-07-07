# CodeQuest CLI

A command-line tool for practicing coding challenges locally. Write solutions in Go, JavaScript, or TypeScript and test them against the same test cases used in the online platform.

## Features

- **Multi-language support**: Go, JavaScript, and TypeScript
- **Local testing**: Run tests without Docker - just needs Go and/or Node.js
- **Offline practice**: Download challenges and work on them locally
- **Fast execution**: Native runtime execution for quick feedback
- **Same test cases**: Uses identical test cases as the web platform

## Installation

### Install from releases (recommended)

Download the latest binary from the [releases page](https://github.com/crisecheverria/codequest/releases):

```bash
# For Linux/macOS
curl -L https://github.com/crisecheverria/codequest/releases/latest/download/codequest-$(uname -s)-$(uname -m) -o codequest
chmod +x codequest
sudo mv codequest /usr/local/bin/

# For Windows
# Download codequest-windows-amd64.exe and add it to your PATH
```

### Install from source

```bash
go install github.com/crisecheverria/codequest@latest
```

## Prerequisites

- **Go 1.23+** (for Go challenges)
- **Node.js 18+** (for JavaScript/TypeScript challenges)

## Usage

### List available challenges

```bash
codequest list
codequest list --language go --difficulty easy
```

### Download a challenge

```bash
codequest fetch challenge-slug
cd challenge-slug
```

### Test your solution

```bash
codequest test
```

### Example workflow

```bash
# List challenges
codequest list --language go

# Download a challenge
codequest fetch add-two-numbers-go

# Move to the challenge directory
cd add-two-numbers-go

# Edit the solution file
vim solution.go

# Test your solution
codequest test
```

## Supported Languages

- **Go**: Full support with native `go run` execution
- **JavaScript**: Full support with Node.js execution
- **TypeScript**: Basic support with type annotation removal

## Challenge Structure

Each challenge creates a workspace with:
- `solution.go` / `solution.js` - Your solution file
- `README.md` - Challenge description and examples
- `.challenge.json` - Challenge metadata (don't modify)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details
